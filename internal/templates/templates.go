package templates

import "html/template"

func GetTemplates() (*template.Template, error) {
	templates, err := template.New("orders").ParseGlob("static/*.html")
	if err != nil {
		return nil, err
	}

	return templates, nil
}
