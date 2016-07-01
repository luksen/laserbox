package laserbox

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
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
var Style = "opacity:1;fill:none;fill-opacity:1;stroke:#ff0000;stroke-width:0.282;stroke-linecap:butt;stroke-miterlimit:10;stroke-dasharray:none;stroke-opacity:1"

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
		if part[0] == '-' {
			parts[i] = part[1:len(part)]
		} else {
			parts[i] = "-" + part
		}
	}
	for i, _ := range parts {
		if i%2 == 0 {
			newParts[l-2-i] = parts[i]
			newParts[l-1-i] = parts[i+1]
		}
	}
	return strings.Join(newParts, " ")
}

func start(x, y float64) *pathString {
	ps := &pathString{
		Strings: []string{fmt.Sprintf("m %f,%f", x, y), ""},
		Cur:     1,
	}
	return ps
}

func (ps *pathString) move(x, y float64) {
	ps.Strings[0] = fmt.Sprintf("m %f,%f", x, y)
}

func (ps *pathString) h(l float64) {
	s := fmt.Sprintf(" h %f", l)
	ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + s
}

func (ps *pathString) v(l float64) {
	s := fmt.Sprintf(" v %f", l)
	ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + s
}

func (ps *pathString) sub() {
	ps.Cur += 1
	ps.Strings = append(ps.Strings, "")
}

func (ps *pathString) draw(dir string, l float64) {
	s := fmt.Sprintf(" %s %f", dir, l)
	ps.Strings[ps.Cur] = ps.Strings[ps.Cur] + s
}

func (ps *pathString) close() {
	s := ps.Strings[0]
	for i := 2; i < len(ps.Strings); i++ {
		s += revert(ps.Strings[i])
	}
	s += ps.Strings[1]
	appendPath(s + " z")
}

func (ps *pathString) line(dir string, l, wall, teeth, sign float64, cut bool) float64 {
	var otherdir string

	if dir == "v" {
		otherdir = "h"
	} else if dir == "h" {
		otherdir = "v"
	}

	inv := 1.0
	if l < 0 {
		l *= -1
		inv = -1
	}

	if cut {
		l -= wall
		teeth -= wall
	}

	for l > teeth {
		if l-teeth < wall {
			ps.draw(dir, inv*(l-wall))
			l -= l - wall
		} else {
			ps.draw(dir, inv*(teeth))
			l -= teeth
		}
		ps.draw(otherdir, sign*wall)
		sign *= -1
		if cut {
			teeth += wall
			cut = false
		}
	}
	if l > 0 {
		ps.draw(dir, inv*l)
	}
	return sign
}

func draw(x, y, width, height, depth, wall, teeth float64) {
	base := start(x, y)
	sign := base.line("h", width, wall, teeth, 1, false)
	sign = base.line("v", height, wall, teeth, -1, (sign == -1))
	sign = base.line("h", -width, wall, teeth, -1, (sign == 1))
	sign = base.line("v", -height, wall, teeth, 1, (sign == 1))
	if sign == -1 {
		base = start(x+wall, y)
		sign = base.line("h", width, wall, teeth, 1, true)
		sign = base.line("v", height, wall, teeth, -1, (sign == -1))
		sign = base.line("h", -width, wall, teeth, -1, (sign == 1))
		sign = base.line("v", -height, wall, teeth, 1, (sign == 1))
	}
	base.close()

	top := start(x, y-wall-depth)
	sign = top.line("h", width, wall, teeth, 1, true)
	top.line("v", -depth, wall, teeth, -1, (sign == 1))
	top.sub()
	sign = top.line("v", -depth, wall, teeth, -1, true)
	if sign == -1 {
		top.move(x+wall, y-wall-depth)
	}
	top.close()

	right := start(x+wall+depth+width, y)
	sign = right.line("v", height, wall, teeth, -1, true)
	right.line("h", depth, wall, teeth, -1, (sign == -1))
	right.sub()
	sign = right.line("h", depth, wall, teeth, -1, true)
	if sign == -1 {
		right.move(x+wall+depth+width, y+wall)
	}
	right.close()

	bottom := start(x+width, y+height+wall+depth)
	sign = bottom.line("h", -width, wall, teeth, -1, true)
	bottom.line("v", depth, wall, teeth, 1, (sign == -1))
	bottom.sub()
	sign = bottom.line("v", depth, wall, teeth, 1, true)
	if sign == 1 {
		bottom.move(x+width-wall, y+height+wall+depth)
	}
	bottom.close()

	left := start(x-wall-depth, y+height)
	sign = left.line("v", -height, wall, teeth, 1, true)
	left.line("h", -depth, wall, teeth, 1, (sign == 1))
	left.sub()
	sign = left.line("h", -depth, wall, teeth, 1, true)
	if sign == 1 {
		left.move(x-wall-depth, y+height-wall)
	}
	left.close()
}

func do(width, height, depth, wall, teeth float64) {
	width = width + 2*wall
	height = height + 2*wall
	depth = depth + wall
	draw(depth+wall, depth+wall, width, height, depth, wall, teeth)
	width = width + 2*wall
	height = height + 2*wall
	depth = depth + wall
	draw(depth+wall, depth+2*depth+height, width, height, depth, wall, teeth)
}

func Do(width, height, depth, wall, teeth float64) string {
	do(width, height, depth, wall, teeth)

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "	")
	err := enc.Encode(svg)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
