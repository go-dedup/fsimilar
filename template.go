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
	"strings"
	"text/template"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const (
	fTmplShellLn   = "shell_ln.tmpl"
	fTmplShellGlob = "/shell_*.tmpl"
	sTmplShellCap  = `#!/bin/bash
# -*- bash -*- 

# If env var FSIM_MIN is not set, exit
[ -z "$FSIM_MIN" ] && {
  echo Env var FSIM_MIN is not set, exiting
  exit 0
}

`
)

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var funcMap = template.FuncMap{
	"Quote": Quote,
	"ToRel": ToRel,
	"Ext":   Ext,
}

////////////////////////////////////////////////////////////////////////////
// Function definitions

func outputFinal(files Files) {
	if !Opts.Final {
		return
	}
	templateP := filepath.Dir(locateTemplate(fTmplShellLn))
	filesG, err := filepath.Glob(templateP + fTmplShellGlob)
	abortOn("Glob template", err)
	for _, f := range filesG {
		verbose(3, "F: '%s' = '%s' / '%s'", f, filepath.Dir(f), filepath.Base(f))
		outputSimilars(f, files, false)
	}
}

func outputSimilars(tmplFile string, files Files, stdout bool) {
	templateV := struct {
		Similars Files
		Opts     OptsT
	}{files, Opts}
	templateF := locateTemplate(tmplFile)
	templateB := filepath.Base(tmplFile)
	verbose(3, "  Similar items (%s) --\n %#v.", templateF, files)
	tmpl, err := template.New(templateB).Funcs(funcMap).
		ParseFiles(templateF)
	abortOn("Parse template", err)
	if stdout {
		err = tmpl.Execute(os.Stdout, templateV)
	} else {
		fo, e := os.OpenFile(Opts.TmpFileP+templateB+".sh",
			os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
		abortOn("Create shell output", e)
		defer fo.Close()
		err = tmpl.Execute(fo, templateV)
	}
	warnOn("Executing template", err)
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

// Quote will quote the file name so as to be savely used in shell.
// It uses single quote to quote the file, thus only need to escape single quote
// by Paul Borman https://play.golang.org/p/VYJ907b8w2
func Quote(s string) string {
	return "'" + strings.Join(strings.Split(s, "'"), `'\''`) + "'"
}

// ToRel will convert `p` to relative path to `base`
func ToRel(base, p string) (string, error) {
	return filepath.Rel(base, p)
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.
func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}
