// Laserbox generates a vector image for laser-cutting a box with fitting lid.
// You can specify the internal dimensions of the box.
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

func (ps *pathString) line(dir string, l, material, teeth, sign float64, cut bool) float64 {
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
		l -= material
		teeth -= material
	}

	for l > teeth {
		if l-teeth < material {
			ps.draw(dir, inv*(l-material))
			l -= l - material
		} else {
			ps.draw(dir, inv*(teeth))
			l -= teeth
		}
		ps.draw(otherdir, sign*material)
		sign *= -1
		if cut {
			teeth += material
			cut = false
		}
	}
	if l > 0 {
		ps.draw(dir, inv*l)
	}
	return sign
}

func draw(x, y, width, height, depth, material, teeth float64) {
	base := start(x, y)
	sign := base.line("h", width, material, teeth, 1, false)
	sign = base.line("v", height, material, teeth, -1, (sign == -1))
	sign = base.line("h", -width, material, teeth, -1, (sign == 1))
	sign = base.line("v", -height, material, teeth, 1, (sign == 1))
	if sign == -1 {
		base = start(x+material, y)
		sign = base.line("h", width, material, teeth, 1, true)
		sign = base.line("v", height, material, teeth, -1, (sign == -1))
		sign = base.line("h", -width, material, teeth, -1, (sign == 1))
		sign = base.line("v", -height, material, teeth, 1, (sign == 1))
	}
	base.close()

	top := start(x, y-material-depth)
	sign = top.line("h", width, material, teeth, 1, true)
	top.line("v", -depth, material, teeth, -1, (sign == 1))
	top.sub()
	sign = top.line("v", -depth, material, teeth, -1, true)
	if sign == -1 {
		top.move(x+material, y-material-depth)
	}
	top.close()

	right := start(x+material+depth+width, y)
	sign = right.line("v", height, material, teeth, -1, true)
	right.line("h", depth, material, teeth, -1, (sign == -1))
	right.sub()
	sign = right.line("h", depth, material, teeth, -1, true)
	if sign == -1 {
		right.move(x+material+depth+width, y+material)
	}
	right.close()

	bottom := start(x+width, y+height+material+depth)
	sign = bottom.line("h", -width, material, teeth, -1, true)
	bottom.line("v", depth, material, teeth, 1, (sign == -1))
	bottom.sub()
	sign = bottom.line("v", depth, material, teeth, 1, true)
	if sign == 1 {
		bottom.move(x+width-material, y+height+material+depth)
	}
	bottom.close()

	left := start(x-material-depth, y+height)
	sign = left.line("v", -height, material, teeth, 1, true)
	left.line("h", -depth, material, teeth, 1, (sign == 1))
	left.sub()
	sign = left.line("h", -depth, material, teeth, 1, true)
	if sign == 1 {
		left.move(x-material-depth, y+height-material)
	}
	left.close()
}

func do(width, height, depth, material, teeth float64, lid bool) {
	width = width + 2*material
	height = height + 2*material
	depth = depth + material
	draw(depth+material, depth+material, width, height, depth, material, teeth)
	totalHeight := int(height+2*depth+2*material) + 1
	if lid {
		width = width + 2*material
		height = height + 2*material
		depth = depth + material
		draw(depth+material, depth+2*depth+height, width, height, depth, material, teeth)
		totalHeight += int(height + 2*depth + 4*material)
	}
	totalWidth := int(width+2*depth+2*material) + 1
	svg.Width = fmt.Sprintf("%dmm", totalWidth)
	svg.Height = fmt.Sprintf("%dmm", totalHeight)
	svg.ViewBox = fmt.Sprintf("0 0 %d %d", totalWidth, totalHeight)
}

// Do will generate the box for the given values.
//     width: inner width of the base area in mm
//     height: inner height of the base area in mm
//     depth: inner depth of the box/height of the walls in mm
//     material: thickness of material in mm
//     teeth: length of teeth in mm
//     lid: whether to draw a lid as well
func Do(width, height, depth, material, teeth float64, lid bool) string {
	if width < 0 {
		width *= -1
	}
	if height < 0 {
		height *= -1
	}
	if depth < 0 {
		depth *= -1
	}
	if material < 0 {
		material *= -1
	}
	if teeth < 0 {
		teeth *= -1
	}
	do(width, height, depth, material, teeth, lid)

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "	")
	err := enc.Encode(svg)
	if err != nil {
		log.Fatal(err)
	}
	svg.Paths = []*Path{}
	return buf.String()
}
