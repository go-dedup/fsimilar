////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"regexp"
	"time"

	"github.com/mkideal/cli"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The FileT type defines all the necessary info for a file.
type FileT struct {
	Dir  string
	Name string
	Ext  string
	Size int
	Hash uint64
	Dist uint8
}

//==========================================================================
// Main dispatcher

func fsimilar(ctx *cli.Context) error {
	ctx.JSON(ctx.RootArgv())
	ctx.JSON(ctx.Argv())
	fmt.Println()
	rootArgv = ctx.RootArgv().(*rootT)

	Opts.Distance, Opts.SizeGiven, Opts.Template, Opts.Verbose =
		rootArgv.Distance, rootArgv.SizeGiven, rootArgv.Template,
		rootArgv.Verbose.Value()

	return fSimilar(rootArgv.Filei)
}

func fSimilar(cin io.Reader) error {
	rand.Seed(time.Now().UTC().UnixNano())
	//tmpfile := fmt.Sprintf("%s.%d", file, 99999999-rand.Int31n(90000000))

	scanner := bufio.NewScanner(cin)
	for scanner.Scan() {
		file := FileT{}
		fn := scanner.Text()

		if Opts.SizeGiven {
			_, err := fmt.Sscan(fn, &file.Size)
			abortOn("Parsing file size", err)
			il := regexp.MustCompile(`^ *\d+ *(.*)$`).FindStringSubmatchIndex(fn)
			//fmt.Println(il)
			fn = fn[il[2]:]
		}
		p, n := filepath.Split(fn)
		file.Dir, file.Name, file.Ext = p, Basename(n), filepath.Ext(n)
		fmt.Printf("  d='%s', n='%s', e='%s', s='%d'\n",
			file.Dir, file.Name, file.Ext, file.Size)
	}

	return scanner.Err()
}
