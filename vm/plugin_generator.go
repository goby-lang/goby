package vm

import (
	"bytes"
	"text/template"
)

type importSection struct {
	HasLib bool
	Pkgs   []pkg
}

type pluginSections struct {
	ImportSection, VarSection string
}

type pkg struct {
	Prefix string
	Name   string
}

type function struct {
	Prefix string
	Name   string
}

func compileTemplate(obj interface{}, sn, tn string) string {
	buffer := &bytes.Buffer{}

	t := template.Must(template.New(sn).Parse(tn))
	err := t.Execute(buffer, obj)

	if err != nil {
		panic(err.Error())
	}

	return buffer.String()
}

func compilePluginTemplate(pkgs []pkg, funcs []function) string {
	is := compileImportSection(pkgs)
	vs := compileVarsSection(funcs)
	p := pluginSections{ImportSection: is, VarSection: vs}

	return compileTemplate(p, "pluginSections", pluginTemplate)
}

func compileImportSection(pkgs []pkg) string {
	return compileTemplate(pkgs, "importSection", importSectionTemplate)
}

func compileVarsSection(funcs []function) string {
	return compileTemplate(funcs, "functionSection", functionSectionTemplate)
}

const importSectionTemplate = `
import(
  {{- range $pkg := .}}
	{{printf "%s %s\n" $pkg.Prefix $pkg.Name -}}
  {{- end}}
)
`

const functionSectionTemplate = `

{{- range $f := . }}
{{printf "var %s = %s.%s\n" $f.Name $f.Prefix $f.Name -}}
{{- end}}
`

const pluginTemplate = `
package main

{{ .ImportSection -}}

{{- .VarSection -}}

func main {}
`
