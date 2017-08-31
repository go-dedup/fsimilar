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

var fv = Files{}

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

	for _, f := range fv {
		verbose(1, "# C: %v.", f.concordance)
	}

	// each file from input
	for ii := range fv {
		similar := []int{}
		similar = append(similar, ii)
		fv[ii].Vstd = true

		// each remaining unvisited candidates
		for jj := ii + 1; jj < len(fv); jj++ {
			if !fv[jj].Vstd {
				// compare it with *each* similar file found so far
				for kk := range similar {
					if !fv[jj].Vstd &&
						Relation(fv[similar[kk]].concordance, fv[jj].concordance) > 0.5 {
						similar = append(similar, jj)
						fv[jj].Vstd = true
					}
				}
			}
		}

		// output unvisited potential similars by each row
		if len(similar) > 1 {
			verbose(1, "# S: %v.", similar)
		}
	}

	return nil
}

func buildVecs(fn string, file FileT) {
	file.concordance = BuildConcordance(file.Name)
	fv = append(fv, file)
}
