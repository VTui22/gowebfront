package pmadmin

import (
	"net/http"
	"fmt"
	"log"
	// "encoding/json"
    "strings"
	"time"

	"github.com/mebusy/gowebfront/dbconn"
    "crypto/hmac"
    "crypto/sha256"
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
    valid_token := CheckRequestToken( r )
    if ! valid_token {
        action := r.URL.Query().Get("action")
        if action == "login" {
            username := r.URL.Query().Get("username")
            password := r.URL.Query().Get("password")

            bValiduser := dbconn.IsValidUser(username, password)
            if !bValiduser {
                http.Redirect( w, r, r.URL.Path , http.StatusSeeOther )
                return true
            }

            expiration := time.Now().Add( 10 * time.Second)
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

