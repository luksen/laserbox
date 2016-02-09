package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type pathString struct {
	Strings []string
	Cur     int
}
type Path struct {
	D     string `xml:"d,attr"`
	Style string `xml:"style,attr"`
}
type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Version string   `xml:"version,attr"`
	Paths   []*Path  `xml:"path"`
}

var svg *Svg = &Svg{
	Xmlns:   "http://www.w3.org/2000/svg",
	Width:   "210mm",
	Height:  "297mm",
	ViewBox: "0 0 210 297",
	Version: "1.1",
}
var Style = "opacity:1;fill:none;fill-opacity:1;stroke:#ff0000;stroke-width:0.99999994;stroke-linecap:butt;stroke-miterlimit:10;stroke-dasharray:none;stroke-opacity:1"

func appendPath(data string) {
	p := &Path{
		D:     data,
		Style: Style,
	}
	svg.Paths = append(svg.Paths, p)
}

func revert(s string) string {
	parts := strings.Fields(s)
	l := len(parts)
	newParts := make([]string, l)
	for i, part := range parts {
		if i%2 == 0 {
			continue
		}
		if i%4 == 3 {
			continue
		}
		if part[0] == '-' {
			parts[i] = part[:len(part)]
		} else {
			parts[i] = "-" + part
		}
	}
	fmt.Println(strings.Join(parts, " "))
	for i, _ := range parts {
		if i%2 == 0 {
			newParts[l-2-i] = parts[i]
			newParts[l-1-i] = parts[i+1]
		}
	}
	fmt.Println(strings.Join(newParts, " "))
	return strings.Join(newParts, " ")
}

func start(x, y float32) *pathString {
	ps := &pathString{
		Strings: []string{fmt.Sprintf("m %f,%f", x, y)},
		Cur:     0,
	}
	return ps
}

func (ps *pathString) h(l float32) {
	s := fmt.Sprintf(" h %f", l)
	ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + s
}

func (ps *pathString) v(l float32) {
	s := fmt.Sprintf(" v %f", l)
	ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + s
}

func (ps *pathString) m(x, y float32) {
	ps.Cur += 1
	ps.Strings = append(ps.Strings, "")
	//s := fmt.Sprintf(" M %f %f", x, y)
	//ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + s
}

func (ps *pathString) draw(dir string, l float32) {
	s := fmt.Sprintf(" %s %f", dir, l)
	ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + s
}

func (ps *pathString) close() {
	ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + " z"
	ps.done()
}

func (ps *pathString) done() {
	s := ps.Strings[0]
	s += revert(ps.Strings[2])
	s += revert(ps.Strings[1])
	//s := strings.Join(ps.Strings, " ")
	appendPath(s)
}

func (ps *pathString) line(dir string, l, wall, teeth, sign float32) float32 {
	var otherdir string

	if dir == "v" {
		otherdir = "h"
	} else if dir == "h" {
		otherdir = "v"
	}

	for l > teeth {
		if l-teeth < wall {
			ps.draw(dir, l-wall)
			l -= l - wall
		} else {
			ps.draw(dir, teeth)
			l -= teeth
		}
		ps.draw(otherdir, sign*wall)
		sign *= -1
	}
	if l > 0 {
		ps.draw(dir, l)
	}
	return sign
}

func base(x, y, width, height, depth, wall, teeth float32) {
	base := start(x, y)
	sign := base.line("h", width, wall, teeth, 1)
	if sign == -1 {
		base.line("v", height-wall, wall, teeth, -1)
	} else {
		base.line("v", height, wall, teeth, -1)
	}
	
	base.m(x, y)
	sign = base.line("v", height, wall, teeth, -1)
	base.m(x, y)
	if sign == 1 {
		base.line("h", width-wall, wall, teeth, 1)
	} else {
		base.line("h", width, wall, teeth, 1)
	}
	base.done()
}

func do(width, height, depth, wall, teeth float32) {
	base(depth+4, depth+4, width, height, depth, wall, teeth)
}

func main() {
	do(25, 21, 15, 3, 8)

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "	")
	if err := enc.Encode(svg); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	f, _ := os.Create("svgtest.svg")
	enc2 := xml.NewEncoder(f)
	enc2.Indent("", "	")
	if err := enc2.Encode(svg); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println()

}
