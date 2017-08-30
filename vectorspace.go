package main

import (
	"math"
	"strings"
)

type float = float64
type Concordance map[string]float

func (con Concordance) Magnitude() float {
	total := 0.0

	for _, v := range con {
		total = total + math.Pow(v, 2)
	}

	return math.Sqrt(total)
}

func BuildConcordance(document string) Concordance {
	var con map[string]float
	con = make(map[string]float)

	words := strings.Split(strings.ToLower(document), " ")

	for _, key := range words {

		_, ok := con[key]

		key = strings.Trim(key, " ")

		if ok && key != "" {
			con[key] = con[key] + 1
		} else {
			con[key] = 1
		}

	}

	return con
}

func Relation(con1 Concordance, con2 Concordance) float {
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
