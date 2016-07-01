package main

import (
	"log"
	"os"
)

import (
	"github.com/luksen/laserbox"
)

func main() {
	xml := laserbox.Do(60, 95, 33.5, 3.5, 10)

	f, err := os.Create("svgtest.svg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(xml)
}
