package utils

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var templatesBase string

func ensureTemplatesBase() string {
	if templatesBase != "" {
		return templatesBase
	}
	candidates := []string{
		"templates",
		filepath.Join("..", "templates"),
		filepath.Join("..", "..", "templates"),
	}
	for _, dir := range candidates {
		if _, err := os.Stat(filepath.Join(dir, "layout.html")); err == nil {
			templatesBase = dir
			return templatesBase
		}
	}
	// fallback to current dir
	templatesBase = "templates"
	return templatesBase
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	base := ensureTemplatesBase()
	layout := filepath.Join(base, "layout.html")
	page := filepath.Join(base, name)
	tmpl, err := template.ParseFiles(layout, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, data)
}
