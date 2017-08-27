////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/fatih/structs"
	"github.com/go-easygen/easygen"
	"github.com/go-easygen/easygen/egVar"

	"github.com/go-dedup/simhash"
	"github.com/go-dedup/simhash/sho"

	"github.com/mkideal/cli"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	oracle = sho.NewOracle()
	sh     = simhash.NewSimhash()
	r      = Opts.Distance

	empty = regexp.MustCompile(`^\s*$`)
	fAll  = make(FAll)       // the all file to FCItem map
	fc    = NewFCollection() // the FCollection that holds everything
)

var tmplFileName = map[bool]string{
	true:  "fsimilar_std.tmpl",
	false: "fsimilar_plain.tmpl",
}

func simCLI(ctx *cli.Context) error {
	// ctx.JSON(ctx.RootArgv())
	// ctx.JSON(ctx.Argv())
	// fmt.Println()
	rootArgv = ctx.RootArgv().(*rootT)

	Opts.Distance, Opts.SizeGiven, Opts.QuerySize, Opts.Verbose =
		rootArgv.Distance, rootArgv.SizeGiven, rootArgv.QuerySize,
		rootArgv.Verbose.Value()
	r = Opts.Distance

	return fSimilar(rootArgv.Filei)
}

func fSimilar(cin io.Reader) error {
	rand.Seed(time.Now().UTC().UnixNano())
	//tmpfile := fmt.Sprintf("%s.%d", file, 99999999-rand.Int31n(90000000))
	buildOracle(cin)
	return dealDups()
}

func buildOracle(cin io.Reader) error {
	// read input line by line
	scanner := bufio.NewScanner(cin)
	for scanner.Scan() {
		fn := scanner.Text()
		if empty.MatchString(fn) {
			continue
		}
		file := FileT{}

		// == Gather file info
		if Opts.SizeGiven {
			_, err := fmt.Sscan(fn, &file.Size)
			abortOn("Parsing file size", err)
			il := regexp.MustCompile(`^ *\d+\s+(.*)$`).FindStringSubmatchIndex(fn)
			//fmt.Println(il)
			fn = fn[il[2]:]
		} else if Opts.QuerySize {
			s, err := os.Stat(fn)
			warnOn("Get file size", err)
			file.Size = int(s.Size())
		} else {
			file.Size = 1
		}
		p, n := filepath.Split(fn)
		file.Dir, file.Name, file.Ext = p, Basename(n), filepath.Ext(n)
		verbose(2, " n='%s', e='%s', s='%d', d='%s'",
			file.Name, file.Ext, file.Size, file.Dir)

		hash := sh.GetSimhash(sh.NewWordFeatureSet([]byte(file.Name)))
		file.Hash = hash
		fi := FCItem{Hash: hash, Index: fc.LenOf(hash)}
		fAll[fn] = fi
		fc.Add(hash, file)

		// == Build similarity knowledge
		if h, d, seen := oracle.Find(hash, r); seen == true {
			verbose(2, "=: Simhash of %x ignored for %x (%d).", hash, h, d)
		} else {
			oracle.See(hash)
			verbose(2, "+: Simhash of %x added.", hash)
		}
	}

	return scanner.Err()
}

func dealDups() error {
	verbose(2, "Deal Dups\n")
	tmpl0 := easygen.NewTemplate().Customize()
	tmpl := tmpl0.Funcs(easygen.FuncDefs()).Funcs(egVar.FuncDefs())

	// process all, the sorted fAll map
	visited := make(HVisited)
	var keys []string
	for k := range fAll {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// for each sorted item according to file path & name
	for _, k := range keys {
		// get next file item from FAll file to FCItem map
		fi := fAll[k]
		verbose(2, "# Item: %v.", fi)
		// and skip if the hash has already been visited
		if visited[fi.Hash] {
			verbose(2, "  Visited hash ignored")
			continue
		}
		visited[fi.Hash] = true
		files, ok := fc.Get(fi.Hash)
		if !ok {
			abortOn("Internal error", errors.New("fc integrity checking"))
		}

		fSizeRef := files[0].Size
		// similarity exist, start digging
		for ii, _ := range files {
			files[ii].Vstd, files[ii].Dist, files[ii].SizeRef = true, 0, fSizeRef
		}
		neighbors := oracle.Search(fi.Hash, r+1)
		// skip if no similar items at this hash
		if len(neighbors) == 0 {
			verbose(2, "  No similar items at this hash")
			continue
		}
		verbose(2, "  Neighbors %v: %v.", oracle.Seen(fi.Hash, r), neighbors)
		for _, nigh := range neighbors {
			visited[nigh.H] = true
			// files to add
			fta := Files{}
			// files more
			fm, ok := fc.Get(nigh.H)
			if ok {
				for ii, _ := range fm {
					if !fm[ii].Vstd {
						fm[ii].Vstd, fm[ii].Dist, fm[ii].SizeRef = true, nigh.D, fSizeRef
						fta = append(fta, fm[ii])
					}
				}
				files = append(files, fta...)
			}
		}

		// One group of similar items found, output
		sort.Sort(files)
		m := structs.Map(struct{ Similars Files }{files})
		verbose(2, "  Similar items -- \n %v.", m)
		easygen.Execute(tmpl, os.Stdout, tmplFileName[false], easygen.EgData(m))
	}

	return nil
}
