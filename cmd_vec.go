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

var (
	fv    = Files{}
	simTh = 0.5
)

func vecCLI(ctx *cli.Context) error {
	rootArgv = ctx.RootArgv().(*rootT)
	argv := ctx.Argv().(*vecT)
	// fmt.Printf("[vec]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
	Opts.SizeGiven, Opts.QuerySize, Opts.Phonetic, Opts.Final, Opts.Verbose =
		rootArgv.SizeGiven, rootArgv.QuerySize,
		rootArgv.Phonetic, rootArgv.Final, rootArgv.Verbose.Value()
	simTh = argv.Threshold
	cmdInit()

	return cmdVec(rootArgv.Filei)
}

func cmdVec(cin io.Reader) error {
	processFileInfo(cin, buildVecs)

	Opts.Ndx = 0
	for _, f := range fv {
		verbose(1, "# C: %v.", f.Conc.String())
	}

	// each file from input
	for ii := range fv {
		// skip those that have already been visited
		if fv[ii].Vstd {
			verbose(2, " = Visited file ignored")
			continue
		}
		// get next file item
		similar := []int{}
		similar = append(similar, ii)
		fSizeRef := fv[ii].Size
		fv[ii].Vstd, fv[ii].Relation, fv[ii].SizeRef = true, 1, fSizeRef

		// each remaining unvisited candidates
		for jj := ii + 1; jj < len(fv); jj++ {
			if !fv[jj].Vstd {
				// compare it with *each* similar file found so far
				for kk := range similar {
					rel := Relation(fv[similar[kk]].Conc, fv[jj].Conc)
					if !fv[jj].Vstd && rel >= simTh {
						// prepare info for Similarity() computation
						fv[jj].Relation, fv[jj].SizeRef = rel, fSizeRef
						if fv[jj].Similarity() >= int(simTh*100) {
							similar = append(similar, jj)
							fv[jj].Vstd = true
							// if similar to one, no need to compare with the rest
							break
						}
					}
				}
			}
		}

		// output unvisited potential similars by each row
		if len(similar) > 1 {
			verbose(1, "# S: %v.", similar)
			files := make(Files, len(similar))
			for kk := range similar {
				files[kk] = fv[similar[kk]]
			}
			Opts.File1st = files[0].Org
			Opts.Ndx++
			outputSimilars(tmplFileName[Opts.Final], files, true)
			outputFinal(files)
		}
	}

	return nil
}

func buildVecs(fn string, file FileT) {
	file.Conc = BuildConcordance(file.Name, doc2words)
	fv = append(fv, file)
}
