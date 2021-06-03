package handle

import (
	"html/template"
	// "path"
    "embed"
    "log"
)

const PAGE_FOLDER="pages"

//go:embed *
var fsLogin embed.FS

var t_login *template.Template

func init() {
    // var t *template.Template
    // t = template.New( tmpl_name  )
    // must will panic if error occur
    // t_login = template.Must( template.ParseFiles( []string { path.Join( PAGE_FOLDER, tmpl_name )  }... ) )
    log.Printf("%+v", t_login )
    // t_login = template.Must( template.ParseFS( fsLogin  ) )
}

