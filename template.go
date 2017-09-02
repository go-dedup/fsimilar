////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar template handling
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"os"
	"text/template"
)

func outputSimilars(tmplFile string, files Files, stdout bool) {
	templateF := locateTemplate(tmplFile)
	verbose(3, "  Similar items (%s) --\n %#v.", templateF, files)
	tmpl, err := template.New(tmplFile).ParseFiles(templateF)
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
