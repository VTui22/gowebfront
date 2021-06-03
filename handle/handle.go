package handle

import (
    "net/http"
    "fmt"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
    // time.Sleep( 10 * time.Second ) // test gracefully shutdown
    fmt.Fprintf(w, "admin" )
}

