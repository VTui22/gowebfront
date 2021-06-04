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
    pmadmin.InitKeyAndPage( "<secret_key>",  "游戏名", db.GetMysqlDB()  )
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



