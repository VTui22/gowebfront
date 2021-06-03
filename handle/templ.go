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

var t_login *template.Template

func init() {
    pages, err := embedFiles.ReadDir( PAGE_FOLDER  )
    if err != nil {
        log.Fatal(err)
    }
    for i, fe := range pages {
        log.Println(i,fe.Name() )
    }

    // embedFiles are FS, use patters to filter files
    t_login = template.Must( template.ParseFS( embedFiles, PAGE_FOLDER + "/login.tmpl"  ) )
    log.Printf("%+v", t_login )
}

