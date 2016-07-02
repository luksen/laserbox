package main

import (
	"flag"
	"log"
	"os"
)

import (
	"github.com/luksen/laserbox"
)

var (
	outputFlag   string
	widthFlag    float64
	heightFlag   float64
	depthFlag    float64
	materialFlag float64
	teethFlag    float64
)

func init() {
	const (
		outputDefault = "laserbox.svg"
		outputUsage   = "the output file"
	)
	flag.StringVar(&outputFlag, "output", outputDefault, outputUsage)
	flag.StringVar(&outputFlag, "o", outputDefault, outputUsage+" (shorthand)")
	const (
		widthDefault = 0
		widthUsage   = "inner width of the base area in mm"
	)
	flag.Float64Var(&widthFlag, "width", widthDefault, widthUsage)
	flag.Float64Var(&widthFlag, "w", widthDefault, widthUsage+" (shorthand)")

	const (
		heightDefault = 0
		heightUsage   = "inner height of the base area in mm"
	)
	flag.Float64Var(&heightFlag, "height", heightDefault, heightUsage)
	flag.Float64Var(&heightFlag, "h", heightDefault, heightUsage+" (shorthand)")

	const (
		depthDefault = 0
		depthUsage   = "innner depth of the box/height of the walls in mm"
	)
	flag.Float64Var(&depthFlag, "depth", depthDefault, depthUsage)
	flag.Float64Var(&depthFlag, "d", depthDefault, depthUsage+" (shorthand)")

	const (
		materialDefault = 3
		materialUsage   = "thickness of material in mm"
	)
	flag.Float64Var(&materialFlag, "material", materialDefault, materialUsage)
	flag.Float64Var(&materialFlag, "m", materialDefault, materialUsage+" (shorthand)")

	const (
		teethDefault = 10
		teethUsage   = "length of teeth in mm"
	)
	flag.Float64Var(&teethFlag, "teeth", teethDefault, teethUsage)
	flag.Float64Var(&teethFlag, "t", teethDefault, teethUsage+" (shorthand)")
}

func main() {
	flag.Parse()
	xml := laserbox.Do(widthFlag, heightFlag, depthFlag, materialFlag, teethFlag)

	f, err := os.Create(outputFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(xml)
}
