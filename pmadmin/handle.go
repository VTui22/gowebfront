package pmadmin

import (
    "net/http"
    // "fmt"
    "log"
    // "encoding/json"
    "time"
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

    log.Printf("token cookie:%+v", cookie )
    return true
}


// return is wheher need login
func Login(w http.ResponseWriter, r *http.Request) bool {
    valid_token := CheckRequestToken( r )
    if ! valid_token {
        action := r.URL.Query().Get("action")
        if action == "login" {
            bValiduser := true // TODO
            if !bValiduser {
                http.Redirect( w, r, r.URL.Path , http.StatusSeeOther )
                return true
            }

            username := r.URL.Query().Get("username")
            password := r.URL.Query().Get("password")
            log.Println( username, password )

            token := "token"
            expiration := time.Now().Add( 10 * time.Second)
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

