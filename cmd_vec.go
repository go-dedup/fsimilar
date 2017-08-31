////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"io"

	"github.com/mkideal/cli"
)

var concordance []Concordance

func vecCLI(ctx *cli.Context) error {
	rootArgv = ctx.RootArgv().(*rootT)
	// argv := ctx.Argv().(*vecT)
	// fmt.Printf("[vec]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
	Opts.SizeGiven, Opts.QuerySize, Opts.Phonetic, Opts.Verbose =
		rootArgv.SizeGiven, rootArgv.QuerySize,
		rootArgv.Phonetic, rootArgv.Verbose.Value()

	return cmdVec(rootArgv.Filei)
}

func cmdVec(cin io.Reader) error {
	processFileInfo(cin, buildVecs)

	for _, c := range concordance {
		verbose(1, "# C: %v.", c)
	}

	return nil
}

func buildVecs(fn string, file FileT) {
	concordance = append(concordance, BuildConcordance(file.Name))
}
