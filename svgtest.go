package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type pathString struct {
	String string
	Sign   float32
	Side   int
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

func start(x, y float32) *pathString {
	ps := &pathString{
		String: fmt.Sprintf("m %f,%f", x, y),
		Sign:   1,
		Side:   0,
	}
	return ps
}

func (ps *pathString) h(l float32) {
	s := fmt.Sprintf(" h %f", l)
	ps.String = ps.String + s
}

func (ps *pathString) v(l float32) {
	s := fmt.Sprintf(" v %f", l)
	ps.String = ps.String + s
}

func (ps *pathString) draw(dir string, l float32) {
	s := fmt.Sprintf(" %s %f", dir, l)
	ps.String = ps.String + s
}

func (ps *pathString) close() {
	ps.String = ps.String + " z"
	ps.done()
}

func (ps *pathString) done() {
	appendPath(string(ps.String))
}

func (ps *pathString) line(l, wall, teeth, sign float32, cut, hinv, end bool) (lowend bool) {
	lowend = false
	var otherdir string
	var dir string
	//var sign float32
	var inv float32

	if cut {
		l -= wall
		teeth -= wall
		if end {
			teeth += wall
		}
	}
	if ps.Sign == -1 {
		l -= wall
	}
	ps.Sign = sign
	switch ps.Side {
	case 0:
		sign = 1
		inv = 1
		dir = "h"
		otherdir = "v"
		if hinv {
			inv *= -1
		}
	case 1:
		sign = -1
		inv = 1
		dir = "v"
		otherdir = "h"
		if hinv {
			sign *= -1
		}
	case 2:
		sign = -1
		inv = -1
		dir = "h"
		otherdir = "v"
		if hinv {
			inv *= -1
		}
	case 3:
		sign = 1
		inv = -1
		dir = "v"
		otherdir = "h"
		if hinv {
			sign *= -1
		}
	}
	for l > teeth {
		if l-teeth < wall {
			ps.draw(dir, inv*(l-wall))
			l -= l - wall
		} else {
			ps.draw(dir, inv*teeth)
			l -= teeth
		}
		ps.draw(otherdir, sign*ps.Sign*wall)
		ps.Sign *= -1
		if cut {
			cut = false
			if !end {
				teeth += wall
			}
		}
	}
	if l > 0 {
		ps.draw(dir, inv*l)
	}
	ps.Side += 1
	ps.Side %= 4
	if ps.Sign == 1 {
		lowend = true
	}
	return
}

func base(x, y, width, height, depth, wall, teeth float32) {
	var low bool
	var nextlow bool
	base := start(x, y)
	above := start(x+wall, y-(wall+2))
	right := start(x+width+(wall+2), y+wall)
	below := start(x+width-wall, y+height+(wall+2))
	left := start(x-(wall+2), y+height-wall)
	base.line(width, wall, teeth, 1, false, false, false)
	base.line(height, wall, teeth, 1, false, false, false)
	base.line(width, wall, teeth, 1, false, false, false)
	base.line(height, wall, teeth, 1, false, false, false)
	above.Side = 2
	nextlow = above.line(width, wall, teeth, -1, true, true, false)
	above.line(depth, wall, teeth, 1, false, true, false)
	above.line(width, wall, width, 1, false, true, false)
	above.line(depth, wall, teeth, 1, true, true, false)
	right.Side = 1
	low = nextlow
	nextlow = right.line(height, wall, teeth, -1, true, true, low)
	right.line(depth, wall, teeth, 1, false, true, low)
	right.line(height, wall, height, 1, false, true, low)
	right.line(depth, wall, teeth, 1, true, true, low)
	below.Side = 0
	low = nextlow
	nextlow = below.line(width, wall, teeth, -1, true, true, low)
	below.line(depth, wall, teeth, 1, false, true, low)
	below.line(width, wall, width, 1, false, true, low)
	below.line(depth, wall, teeth, 1, true, true, low)
	left.Side = 3
	low = nextlow
	nextlow = left.line(height, wall, teeth, -1, true, true, low)
	left.line(depth, wall, teeth, 1, false, true, low)
	left.line(height, wall, height, 1, false, true, low)
	left.line(depth, wall, teeth, 1, true, true, low)
	base.close()
	above.close()
	right.close()
	below.close()
	left.close()
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
