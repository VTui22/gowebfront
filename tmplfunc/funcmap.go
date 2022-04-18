package tmplfunc

import (
	"html/template"
)

// use string directly may cause "ZgotmplZ"
// convert to template.HTML alias
func SafeHtml (s string) template.HTML {
    return template.HTML(s)
}

// use string directly may cause "ZgotmplZ"
// convert to template.CSS alias
func SafeCss (s string) template.CSS {
    return template.CSS(s)
}


func SafeJs (s string) template.JS {
    return template.JS(s)
}

