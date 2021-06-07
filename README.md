# gowebfront

```
├── pmadin   
│   ├── handle.go
│   ├── pages
│   │   └── login.tmpl
│   ├── static
│   │   └── style.css
│   └── templ.go
└── tmplfunc 
    └── funcmap.go
```

- pmadmin/tmplfunc
    - useful function for handle html, css, js content in http/template
- pmadmin/templ.go
    - load/init templates
- pmadmin/handle.go
    - main handle
- pmadmin/pages
    - template files
- pmadmin/static
    - css, js, etc...


## useage 

see example/server.go


```go
    pmadmin.InitKeyAndPage( "<secret_key>",  "游戏名", db.GetMysqlDB() , []string{ "127.0.0.1/32", "::1/64" } ) // , "10.192.0.0/16"
    r.HandleFunc("/pmadmin", func(w http.ResponseWriter, r *http.Request) {
        bNeedLogin := pmadmin.Login( w,r )
        if bNeedLogin {
            return
        }
        // your main page
        fmt.Fprintf(w, "admin" )
    } )


    // r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    dbconn.PrepareMysqlTable()
```

detailed example

```go
func getPmToken( r *http.Request  ) string {
    cookie, err := r.Cookie( pmadmin.PM_TOKEN_NAME )
    if err != nil {
        return ""
    }
    token := cookie.Value
    return token
}

const COOKIE_NAME_ERRINFO =  "errinfo" 
func getErrorInfoCookie( r *http.Request  ) string {
    cookie, err := r.Cookie( COOKIE_NAME_ERRINFO )
    if err != nil {
        return ""
    }
    return cookie.Value
}




func WebPMHandle(w http.ResponseWriter, r *http.Request) {
    bNeedLogin := pmadmin.Login( w,r )
    if bNeedLogin {
        return
    }

    // your main page
    var _page_data HTML_PAGE_t
    _page_data.Css = css
    _page_data.Title = "Yet Another Title"
    _page_data.GameList = conf.CROSS_PROMOTION_GAMES
    // _page_data.GameList = []string{}

    // use cookie to pass errinfo, because 301/302 can not use headers
    _page_data.Error = getErrorInfoCookie(r)

    action := r.URL.Query().Get("action")
    switch action {
    case "create":
        game := r.URL.Query().Get("game")
        var errinfo = ""
        if game == "" {
            errinfo = fmt.Sprintf( "game '%s' is not a valid game", game )
        }
        starttime, err := strconv.ParseInt( r.URL.Query().Get("starttime") , 10, 0)
        if err != nil {
            errinfo = err.Error()
        }
        endtime , err := strconv.ParseInt( r.URL.Query().Get("endtime") , 10, 0)
        if err != nil {
            errinfo = err.Error()
        }
        pm_token := getPmToken(r)
        operator := strings.Split( pm_token , "." )[0]

        if errinfo == "" {
            // no error, try create
            err := dbconn.PublishCrossPormotion( game, starttime, endtime, operator )
            if err != nil {
                log.Println(err)
                errinfo = err.Error()
            }
        }

        if errinfo != "" {  // use cookie to pass errinfo, because 301/302 can not use headers
            expiration := time.Now().Add( 2 * time.Second)
            cookie := http.Cookie{Name: COOKIE_NAME_ERRINFO, Value:errinfo, Expires: expiration }
            http.SetCookie(w, &cookie)
        }
        http.Redirect( w, r, r.URL.Path  , http.StatusSeeOther )

    default:
        err := t_crosspromotion.Execute(w, _page_data )
        if err != nil {
            fmt.Fprintf(w, err.Error())
        }
    }

}


const PAGE_FOLDER="pages"

//go:embed pages/*.html
var embedFiles embed.FS

//go:embed static/style.css
var css string

var t_crosspromotion *template.Template

func init() {

    funcMap := template.FuncMap{
        "SafeHtml": tmplfunc.SafeHtml ,
        "SafeCss": tmplfunc.SafeCss ,
    }
    t_crosspromotion = template.Must( template.New("index.html").Funcs(funcMap).ParseFS( embedFiles, PAGE_FOLDER + "/index.html"  ) )
}


type HTML_PAGE_t struct {
    Title string
    Css string
    GameList []string
    Error string
}
```

