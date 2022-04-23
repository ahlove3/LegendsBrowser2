package df

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

var eventsTemplate = template.Must(template.New("").Parse(`// Code generated by legendsbrowser; DO NOT EDIT.
package model

{{- range $name, $obj := $.Objects }}
{{- if $obj.IsSubTypeOf "HistoricalEvent" }}
func (x *{{ $obj.Name }}) Html(c *context) string { return "UNKNWON {{ $obj.Name }}" }
{{- end }}
{{- end }}
`))

func GenerateEventsCode(objects *Metadata) error {
	file, _ := json.MarshalIndent(objects, "", "  ")
	_ = ioutil.WriteFile("model.json", file, 0644)

	f, err := os.Create("../backend/model/events.go")
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	err = eventsTemplate.Execute(&buf, struct {
		Objects *Metadata
		Modes   []bool
	}{
		Objects: objects,
		Modes:   []bool{false, true},
	})
	if err != nil {
		return err
	}
	p, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("WARN: could not format source", err)
		p = buf.Bytes()
	}
	_, err = f.Write(p)
	return err
}

func (o Object) IsSubTypeOf(t string) bool {
	return o.SubTypeOf != nil && *o.SubTypeOf == t
}
