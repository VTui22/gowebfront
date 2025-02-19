package main

import (
    "flag"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/mebusy/goweb/db"
    "github.com/mebusy/goweb/webrouter"
    "github.com/mebusy/goweb/webserver"
    "net/http"
    "runtime"
    "strings"
    "server/dbconn"
    "os"
    "github.com/mebusy/gowebfront/pmadmin"
)

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
    // time.Sleep( 10 * time.Second ) // test gracefully shutdown
    fmt.Fprintf(w, "ok")
}
func docHandler(w http.ResponseWriter, r *http.Request) {
    if runtime.GOOS == "darwin"  {
        fmt.Fprintf(w, strings.Replace(webrouter.GetAPIDoc(), "{DOMAIN}", r.Host, -1))
    } else if strings.Contains(r.Host, "bisoft.org" ) {
        uri_prefix := os.Getenv( "INTERNAL_URI_PREFIX" ) 
        fmt.Fprintf(w, strings.Replace(webrouter.GetAPIDoc(), "{DOMAIN}", "https://" + r.Host + uri_prefix , -1))
    } else {
        fmt.Fprintf(w, "forbidden")
    }
}

func ipTestHandler(w http.ResponseWriter, r *http.Request) {
    // time.Sleep( 10 * time.Second ) // test gracefully shutdown
    fmt.Fprintf(w, "x-for:%s, x-real-ip:%s", r.Header.Get("X-Forwarded-For") , r.Header.Get("X-Real-Ip") )
}


var listenPort = flag.Int("p", 5757, "port")
var verbose = flag.Int("v", 0, "verbose")

var GitCommit string

func main() {
    // runtime.GOMAXPROCS(1)  // even if you set max procs to 1, you also have to handle concurrency data racing problem.
    defer db.MysqlClose()

    flag.Parse()

    r := mux.NewRouter()
    // r.HandleFunc( "/bot", webhookHandleGET).Methods("GET")
    // r.HandleFunc( "/bot", webhookHandlePOST).Methods("POST")

    r.HandleFunc("/", catchAllHandler)
    r.HandleFunc("/doc", docHandler)
    r.HandleFunc("/iptest", ipTestHandler)


    // r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    dbconn.PrepareMysqlTable()

    // 
    // key for calc hmac
    // game title
    // db instance
    // whitelist
    // uri_prefix, used when a prefix is add by nginx proxy
    pmadmin.InitKeyAndPage( "<secret_key>",  "游戏名", db.GetMysqlDB() , []string{ "127.0.0.1/32", "::1/64" }, "" ) // , "10.192.0.0/16"
    r.HandleFunc("/pmadmin", func(w http.ResponseWriter, r *http.Request) {
        bNeedLogin := pmadmin.Login( w,r )
        if bNeedLogin {
            return
        }
        // your main page
        fmt.Fprintf(w, "admin" )
    } )


    webserver.StartServer(r, *listenPort, *verbose, GitCommit)
}


