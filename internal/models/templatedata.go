package models

import "github.com/satya-kr/bookings/internal/forms"

// TemplateData holds data to send between handler to template
type TemplateData struct {
	StringMap  map[string]string
	IntMap     map[string]int
	FloatMap   map[string]float32
	Data       map[string]interface{}
	CSRFToken  string
	FlashMsg   string
	WarningMsg string
	ErrorMsg   string
	Form       *forms.Form
}
