////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar template handling
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var funcMap = template.FuncMap{
	"ToRel": ToRel,
}

////////////////////////////////////////////////////////////////////////////
// Function definitions

func outputSimilars(tmplFile string, files Files, stdout bool) {
	templateF := locateTemplate(tmplFile)
	verbose(3, "  Similar items (%s) --\n %#v.", templateF, files)
	tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles(templateF)
	abortOn("Parse template", err)
	err = tmpl.Execute(os.Stdout, files)
	abortOn("Executing template", err)
}

func locateTemplate(tmplFile string) string {
	var templateF string
	if IsExist(tmplFile) {
		templateF = tmplFile
	} else if tf := Opts.ExecPath + "/" + tmplFile; IsExist(tf) {
		templateF = tf
	} else if tf := Opts.CfgPath + "/" + tmplFile; IsExist(tf) {
		templateF = tf
	} else if tf := "/etc/fsimilar/" + tmplFile; IsExist(tf) {
		templateF = tf
	} else {
		abortOn("Locate template",
			fmt.Errorf("Template file '%s' not found", tmplFile))
	}
	return templateF
}

func ToRel(base, p string) (string, error) {
	return filepath.Rel(base, p)
}
