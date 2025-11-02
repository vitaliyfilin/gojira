package load

import (
	"bytes"
	"os"
	"text/template"
)

func LoadBodyTemplate(path string) (*template.Template, error) {
	if path == "" {
		return nil, nil
	}
	bb, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return template.New("body").Parse(string(bb))
}

func RenderBody(tmpl *template.Template) ([]byte, error) {
	if tmpl == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
