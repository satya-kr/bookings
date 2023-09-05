package render

import (
	"bytes"
	"fmt"
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/models"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"HumanDate": HumanDate,
}

var app *config.AppConfig

var pathToTemplate, _ = filepath.Abs("../../templates")

// NewRenderer NewTemplates sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.FlashMsg = app.Session.PopString(r.Context(), "flash")
	td.ErrorMsg = app.Session.PopString(r.Context(), "error")
	td.WarningMsg = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuth = 1
	}

	return td
}

// Template to render templates
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var err error

	var templates map[string]*template.Template

	if app.UseCache {
		templates = app.TemplateCache
	} else {
		templates, err = GetTemplatesCache()
		if err != nil {
			fmt.Println("Error getting template", err)
			return err
		}
	}

	t, ok := templates[tmpl+".page.tmpl"]
	if !ok {
		fmt.Println("Opps, Template doesn't exist")
		return err
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error waiting template to browse:", err)
		return err
	}

	return nil
}

// GetTemplatesCache collect all templates then merge them with layout
func GetTemplatesCache() (map[string]*template.Template, error) {

	//get the Template Cache from app congig

	myCache := map[string]*template.Template{}

	//pathToTemplate, pathErr := filepath.Abs("../../templates")
	//if pathErr != nil {
	//	return myCache, pathErr
	//}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	// fmt.Println("\n\n Pasges", pages)
	if err != nil {
		return myCache, err
	}

	matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
	// fmt.Println("\n\n Matches", matches)
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
