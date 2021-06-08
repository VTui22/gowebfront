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
    pmadmin.InitKeyAndPage( "<secret_key>",  "游戏名", db.GetMysqlDB() , []string{ "127.0.0.1/32", "::1/64" }, "/cp" ) // , "10.192.0.0/16"
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

<details>
<summary>
detailed example
</summary>

```go
package webpm

import (
    "embed"
    "fmt"
    "html/template"
    "net/http"
    "os"
    "server/conf"
    "server/dbconn"
    "time"

    "strconv"
    "strings"

    // "log"

    // "github.com/mebusy/goweb/tools"
    "github.com/mebusy/goweb/tools"
    "github.com/mebusy/gowebfront/pmadmin"
    "github.com/mebusy/gowebfront/tmplfunc"
)

var INTERNAL_URI_PREFIX string
func init() {
    INTERNAL_URI_PREFIX  = os.Getenv( "INTERNAL_URI_PREFIX" )
}

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
    _page_data.Title = "管理后台"
    _page_data.IsProd = strings.Contains( r.Host, "prod" )
    // _page_data.GameList = []string{}
    confs, _ := dbconn.GetCrossPromotionConfs()
    _page_data.CrossPromotionConfs = confs
    
    games_not_in_confs := []string{}
    for _, game := range conf.CROSS_PROMOTION_GAMES {
        _, ok := confs[game]
        if !ok {
            games_not_in_confs = append(  games_not_in_confs  , game )
        }
    }
    _page_data.GameList = games_not_in_confs

    // use cookie to pass errinfo, because 301/302 can not use headers
    _page_data.Error = getErrorInfoCookie(r)

    action := r.URL.Query().Get("action")
    var errinfo = ""
    switch action {
    case "create":
        game := r.URL.Query().Get("game")
        if game == "" {
            errinfo = fmt.Sprintf( "game '%s' is not a valid game", game )
        }
        starttime, err := utctimestampFromBeijingDate( r.URL.Query().Get("starttime") )
        if err != nil {
            errinfo = err.Error()
        }
        endtime , err := utctimestampFromBeijingDate( r.URL.Query().Get("endtime") )
        if err != nil {
            errinfo = err.Error()
        }
        pm_token := getPmToken(r)
        operator := strings.Split( pm_token , "." )[0]

        if errinfo == "" {
            // no error, try create
            err := dbconn.PublishCrossPormotion( game, starttime, endtime, operator )
            if err != nil {
                errinfo = err.Error()
            }
        }
        if errinfo != "" {  // use cookie to pass errinfo, because 301/302 can not use headers
            expiration := time.Now().Add( 2 * time.Second)
            cookie := http.Cookie{Name: COOKIE_NAME_ERRINFO, Value:errinfo, Expires: expiration }
            http.SetCookie(w, &cookie)
        }
        http.Redirect( w, r, INTERNAL_URI_PREFIX + r.URL.Path  , http.StatusSeeOther )

    case "delete":
        id, err := strconv.ParseInt( r.URL.Query().Get("id"), 10,0 )
        if err != nil {
            errinfo = err.Error()
        }
        if errinfo == "" {
            err := dbconn.DeleteCrossPormotion(id)
            if err != nil {
                errinfo = err.Error()
            }
        }
        if errinfo != "" {  // use cookie to pass errinfo, because 301/302 can not use headers
            expiration := time.Now().Add( 2 * time.Second)
            cookie := http.Cookie{Name: COOKIE_NAME_ERRINFO, Value:errinfo, Expires: expiration }
            http.SetCookie(w, &cookie)
        }
        http.Redirect( w, r, INTERNAL_URI_PREFIX + r.URL.Path  , http.StatusSeeOther )

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
        "utctimestampToBeijingDate": utctimestampToBeijingDate,
        "ShowExp": ShowExp,
        "canDelete": canDelete,
    }
    t_crosspromotion = template.Must( template.New("index.html").Funcs(funcMap).ParseFS( embedFiles, PAGE_FOLDER + "/index.html"  ) )
}


type HTML_PAGE_t struct {
    Title string
    Css string
    GameList []string
    Error string
    CrossPromotionConfs map[string] dbconn.CrossPromotionConf_t
    IsProd bool
}

var secondsEastOfUTC = int((8 * time.Hour).Seconds())
var beijing = time.FixedZone("Beijing Time", secondsEastOfUTC)

func utctimestampFromBeijingDate( date string ) (int64,error) {

    RFC3339     := "2006010215"
    time_target, err  := time.ParseInLocation( RFC3339,   date , beijing )
    if err != nil {
        return  -1, err 
    }

    return time_target.UTC().Unix(), nil
}

func utctimestampToBeijingDate( seconds int64 ) string {
    if seconds == 0 {
        return "N/A"
    }
    FORMAT := "2006-01-02 15:04"

    t := time.Unix( seconds,0 ).UTC()
    return t.In(beijing).Format( FORMAT )
}

func ShowExp( startTime, endTime int64 ) string {
    sec_now := tools.GetSeconds()
    if startTime >= endTime {
        return "过期"
    }
    if sec_now < startTime {
        return "未生效"
    }
    if sec_now >= startTime && sec_now <= endTime {
        return "生效"
    }
    if sec_now > endTime {
        return "过期"
    }

    return "错误"
}

func canDelete( startTime, endTime int64 ) bool {
    sec_now := tools.GetSeconds()
    if sec_now >= startTime && sec_now <= endTime {
        return false
    }
    return true
}
```

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <style>{{ .Css | SafeCss }}</style>
</head>
<body>
<header>
    {{.Title}}
</header>

<p><center>
<div id="id_separate_bar">==================  列表 ==================== </div>
<table border="1">
    <thead>
        <th>游戏</th>
        <th>开始时间</th>
        <th>结束时间</th>
        <th>生效</th>
        <th>Ops</th>
        <th>操作</th>
    </thead>
    <tbody>
        {{ $isdev := not .IsProd }}
        {{ range $key, $conf := .CrossPromotionConfs }}
            <tr>
                <td> {{$key}} </td>
                <td> {{$conf.StartTime | utctimestampToBeijingDate}} </td>
                <td> {{$conf.EndTime | utctimestampToBeijingDate}} </td>
                <td> {{ShowExp $conf.StartTime $conf.EndTime}} </td>
                <td> {{$conf.Operator}} </td>
                {{$candel := canDelete $conf.StartTime $conf.EndTime }}
                <td> {{if or $candel $isdev  }} 
                    <form id="id_form_delete{{$conf.Id}}"> <!-- default action, same url -->
                        <input type="hidden" name="action" value="delete" />
                        <input type="hidden" name="id" value="{{$conf.Id}}" />
                        <button class="cBtn_del"  type="submit" form="id_form_delete{{$conf.Id}}">删除</button>  
                    </form>
                    {{end}}
                </td>
            </tr>
        {{ end }}
    </tbody>
</table>
</center></p>

<!-- if game list is not empty, show create section-->
{{ if .GameList }}
<p><center>

<div id="id_separate_bar">================== 发布 ==================== </div>
<form id="id_form_create"> <!-- default action, same url -->
    <!-- hideen filed -->
    <input type="hidden" name="action" value="create" />
    <div>Choose game</div>
    <div>
    <select name="game" id="id_game" form="id_form_create">
        {{ range $index, $gamename := .GameList }}
        <option value="{{$gamename}}">{{$gamename}}</option>
        {{ end }}
    </select>
    </div>
    <div>时间格式(北京时间): 2016010215</div>
    <div>开始时间</div>
    <input type="number" id="id_starttime" name="starttime" minlength="10" maxlength="10" placeholder="1970010208" required oninput="javascript: if (this.value.length > this.maxLength) this.value = this.value.slice(0, this.maxLength);">
    <div>结束时间</div>
    <input type="number" id="id_endtime" name="endtime" minlength="10" maxlength="10" placeholder="1970010208" required oninput="javascript: if (this.value.length > this.maxLength) this.value = this.value.slice(0, this.maxLength);">

    <br><br>
    <button type="submit" form="id_form_create">发布</button>
</form>

</center></p>
{{ end }}


<!-- error info -->
{{ if .Error}}
<br><br>
<div id="id_errorinfo"> {{.Error}} </div>
{{ end }}

</body>
</html>

```

</details>



