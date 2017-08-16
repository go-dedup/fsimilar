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
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/go-dedup/simhash"
	"github.com/go-dedup/simhash/sho"

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

type Files []FileT

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

	oracle := sho.NewOracle()
	sh := simhash.NewSimhash()
	r := Opts.Distance
	fAll := make(FAll)
	fc := NewFCollection()
	fClose := make(FCs, 16384)

	// read input line by line
	scanner := bufio.NewScanner(cin)
	for scanner.Scan() {
		file := FileT{}
		fn := scanner.Text()

		// == Gather file info
		if Opts.SizeGiven {
			_, err := fmt.Sscan(fn, &file.Size)
			abortOn("Parsing file size", err)
			il := regexp.MustCompile(`^ *\d+\s+(.*)$`).FindStringSubmatchIndex(fn)
			//fmt.Println(il)
			fn = fn[il[2]:]
		} else {
			s, err := os.Stat(fn)
			warnOn("Get file size", err)
			file.Size = int(s.Size())
		}
		p, n := filepath.Split(fn)
		file.Dir, file.Name, file.Ext = p, Basename(n), filepath.Ext(n)
		verbose(1, " n='%s', e='%s', s='%d', d='%s'",
			file.Name, file.Ext, file.Size, file.Dir)

		hash := sh.GetSimhash(sh.NewWordFeatureSet([]byte(file.Name)))
		fi := FCItem{hash, fc.LenOf(hash)}
		fAll[fn] = fi
		fc.Set(hash, file)

		// == Check similarity
		if h, d, seen := oracle.Find(hash, r); seen == true {
			fClose.SetClose(fi, h, d)
			verbose(1, "=: Simhash of %x ignored for %x (%d).", hash, h, d)
		} else {
			oracle.See(hash)
			verbose(1, "+: Simhash of %x added.", hash)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// process all, the sorted fAll map
	visited := make(HVisited)
	var keys []string
	for k := range fAll {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fi := fAll[k]
		if visited[fi.Hash] {
			continue
		}

		visited[fi.Hash] = true
	}

	return nil
}
