//  this file data we can access from anywhere to the application

package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// Appconfig hold application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template //pointer to template
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
