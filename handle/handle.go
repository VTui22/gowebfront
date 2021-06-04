package handle

import (
    "net/http"
    "fmt"
    "log"
    // "encoding/json"
    "time"
)


const PM_TOKEN_NAME = "pmtoken"
func checkRequestToken( r *http.Request  ) bool {
    cookie, err := r.Cookie( PM_TOKEN_NAME )
    if err != nil {
        log.Println(err)
        return false
    }

    log.Printf("token cookie:%+v", cookie )
    return true
}


func AdminHandler(w http.ResponseWriter, r *http.Request) {
    valid_token := checkRequestToken( r )
    if ! valid_token {
        action := r.URL.Query().Get("action")
        if action == "login" {
            username := r.URL.Query().Get("username")
            password := r.URL.Query().Get("password")
            log.Println( username, password )

            token := "token"
            expiration := time.Now().Add( 10 * time.Second)
            cookie := http.Cookie{Name: PM_TOKEN_NAME, Value:token, Expires:expiration}
            http.SetCookie(w, &cookie)

            http.Redirect( w, r, r.URL.Path , 200 )
        } else {
            // to login
            t_login.Execute( w, _page_data )
        }
        return
    }

    fmt.Fprintf(w, "admin" )
}

