package pmadmin

import (
    "net"
	"net/http"
	"fmt"
	"log"
	// "encoding/json"
    "strings"
	"time"

    "crypto/hmac"
    "crypto/sha256"
    "strconv"
)


const PM_TOKEN_NAME = "pmtoken"
var TOKEN_KEY = ""
func CheckRequestToken( r *http.Request  ) bool {
    if TOKEN_KEY == "" {
        log.Fatal( "TOKEN_KEY can not empty!" )
    }

    cookie, err := r.Cookie( PM_TOKEN_NAME )
    if err != nil {
        log.Println(err)
        return false
    }

    token := cookie.Value
    log.Printf("pmtoken in cookie:%s", token )

    items := strings.Split( token, "|" )
    if len(items) != 2 {
        return false
    }
    s1 := items[0]
    s2 := items[1]

    // check expiration
    subitems:= strings.Split( s1 , "." )
    if len(subitems) != 2 {
        return false
    }
    exp_seconds, err := strconv.ParseInt( subitems[1], 10, 0 )
    if err != nil {
        log.Println(err)
        return false
    }
    if time.Now().UTC().Unix() > exp_seconds {
        log.Println("token expired")
        return false
    }

    // check token
    key := []byte( TOKEN_KEY )
    mac := hmac.New(sha256.New, key)
    s2_recalc := fmt.Sprintf("%x", mac.Sum([]byte(s1)) )
    if s2 != s2_recalc {
        return false
    }
    return true
}

func generateToken( username string, exp_seconds int64 ) string {
    s1 := fmt.Sprintf( "%s.%d", username, exp_seconds ) 

    key := []byte( TOKEN_KEY )
    mac := hmac.New(sha256.New, key)
    s2 := fmt.Sprintf("%x", mac.Sum([]byte(s1)) )
    // log.Println( "token raw string:", s1, "hmac len:%d",len(s2) )
    return fmt.Sprintf( "%s|%s", s1, s2 )
}

// return is wheher need login
func Login(w http.ResponseWriter, r *http.Request) bool {
    xforip := r.Header.Get("X-Forwarded-For")
    if !isInWhiteList( xforip ) {
        log.Printf( "ip '%s' not in whitelist", xforip )
        fmt.Fprintf( w, "forbiden" )
        return true
    }

    valid_token := CheckRequestToken( r )
    if ! valid_token {
        action := r.URL.Query().Get("action")
        if action == "login" {
            username := r.URL.Query().Get("username")
            password := r.URL.Query().Get("password")

            bValiduser := isValidUser(username, password)
            if !bValiduser {
                http.Redirect( w, r, r.URL.Path , http.StatusSeeOther )
                return true
            }

            expiration := time.Now().Add( 15 * time.Minute)
            exp_seconds := expiration.UTC().Unix()
            token := generateToken( username, exp_seconds )
            cookie := http.Cookie{Name: PM_TOKEN_NAME, Value:token, Expires:expiration}
            http.SetCookie(w, &cookie)

            // redirect, clean the query in browser URL
            http.Redirect( w, r, r.URL.Path , http.StatusSeeOther )
        } else { // not login
            // to login
            t_login.Execute( w, _page_data )
        }

        // return nontheless if token is valid
        return true
    }
    return false
}


var whitelist = []string {}

func setWhiteList( _whitelist []string ) {
    if _whitelist != nil {
        whitelist = _whitelist
    }
}

func isInWhiteList( xforip string ) bool {
    // whiteList
    // realip := r.Header.Get("X-Real-Ip")
    for _, whiteip := range whitelist {
        // if ip == xforip {
        //     return true
        //     break
        // }
        xip := net.ParseIP(xforip)

        _,white_ipnet,err := net.ParseCIDR( whiteip )
        if err != nil {
            log.Println( err )
            continue
        }
        if white_ipnet.Contains( xip ) {
            log.Printf( "%s is in whitelist %s",xforip, whiteip  )
            return true
        }
    }
    return false
}

