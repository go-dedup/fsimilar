package main

import (
	"bytes"
	"fmt"
	"math"
	"sort"

	"github.com/go-dedup/text"
)

type Concordance map[string]float64

func (con Concordance) Magnitude() float64 {
	total := 0.0

	for _, v := range con {
		total = total + math.Pow(v, 2)
	}

	return math.Sqrt(total)
}

func BuildConcordance(document string, doc2words text.Doc2Words) Concordance {
	var con map[string]float64
	con = make(map[string]float64)

	words := doc2words(document)

	for _, key := range words {
		_, ok := con[key]
		if ok && key != "" {
			con[key] = con[key] + 1
		} else {
			con[key] = 1
		}
	}

	return con
}

func Relation(con1 Concordance, con2 Concordance) float64 {
	topvalue := 0.0

	for name, count := range con1 {
		_, ok := con2[name]

		if ok {
			topvalue = topvalue + (count * con2[name])
		}
	}

	mag := con1.Magnitude() * con2.Magnitude()

	if mag != 0 {
		return topvalue / mag
	} else {
		return 0
	}
}

func (con Concordance) String() string {
	buf := bytes.NewBufferString("")

	var keys []string
	for k := range con {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf.WriteString("[ ")
	for _, k := range keys {
		fmt.Fprintf(buf, "%s:%d ", k, int(con[k]))
	}
	buf.WriteByte(']')
	return buf.String()
}
