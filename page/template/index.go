package template

import (
	// "github.com/MerinEREN/iiPackages/page/content"
	// usr "github.com/MerinEREN/iiPackages/datastore/user"
	"html/template"
	"log"
	"net/http"
)

var (
	html = template.Must(template.ParseFiles("../iiClient/public/index.html"))
	// RenderIndex      = renderTemplate("index")
)

// RenderIndex "Exported functions should have a comment"
func RenderIndex(w http.ResponseWriter) {
	err := html.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

/* func renderTemplate(page string) func(w http.ResponseWriter, pc *content.Page) {
	return func(w http.ResponseWriter, pc *content.Page) {
		if page == "index" {
			err := html.Execute(w, pc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatalln(err)
			}
		}
		TemplateRendered = true
	}
} */
