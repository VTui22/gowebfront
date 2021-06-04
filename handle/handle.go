package handle

import (
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
    "log"
    "bytes"
    // "encoding/json"
    "time"
)


func checkRequestToken( r *http.Request  ) bool {
    cookie, err := r.Cookie( "token" )
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
        // to login
        var b bytes.Buffer
        err := t_login.Execute( &b, _page_data )
        if err != nil {
            fmt.Fprintf( w , err.Error() )
            return
        }
        b.WriteTo(w)
        return
    }
    fmt.Fprintf(w, "admin" )
}


func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"] // verified by mux route
    password := vars["password"] // verified by mux route
    // is valid user ?
    log.Println( username, password )

    token := "token"

    // m := map[string]interface{} {}  // empty do nothing
    // m["errcode"] = -1 
    // b, _ := json.Marshal( &m )

    // set cookie
    expiration := time.Now().Add( 10 * time.Second)
    cookie := http.Cookie{Name: "pmtoken",Value:token, Expires:expiration}
    http.SetCookie(w, &cookie)

    fmt.Fprintf(w, "{}" )
}

