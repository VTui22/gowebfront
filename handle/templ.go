package handle

import (
	"html/template"
	// "path"
    "embed"
    "log"
)

const PAGE_FOLDER="pages"

//go:embed pages/*.tmpl
var embedFiles embed.FS

//go:embed static/style.css
var css string

var t_login *template.Template

func init() {

    funcMap := template.FuncMap{
        "toAttr": func(s string) template.HTMLAttr{
            return template.HTMLAttr(s)
        },
        "toHtml": func(s string) template.HTML{
            return template.HTML(s)
        },
    }

    /*
    pages, err := embedFiles.ReadDir( PAGE_FOLDER  )
    if err != nil {
        log.Fatal(err)
    }
    for i, fe := range pages {
        log.Println(i,fe.Name() )
    }
    //*/

    // embedFiles are FS, use patters to filter files
    // Make sure the argument you pass to template.New is the base name of one of the files in the list you pass to ParseFiles.
    t_login = template.Must( template.New("login.tmpl").Funcs(funcMap).ParseFS( embedFiles, PAGE_FOLDER + "/login.tmpl"  ) )
}



type HTML_PAGE_t struct {
    Title string
    Css string
}

var _page_data HTML_PAGE_t

func InitFrontPage( title string ) {
    _page_data.Title = title
    _page_data.Css = css
    log.Printf("css: %s",css)
}
