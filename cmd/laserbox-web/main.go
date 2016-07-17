package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

import (
	"github.com/luksen/laserbox"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("website.html")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	_, err = io.Copy(w, f)
	if err != nil {
		fmt.Fprint(w, err)
	}
}

func svgHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/svg/")
	params := strings.Split(path, "/")
	if len(params) != 6 {
		fmt.Fprintln(w, "width/height/depth/material/teeth/lid")
		return
	}
	width, err := strconv.ParseFloat(params[0], 64)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	height, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	depth, err := strconv.ParseFloat(params[2], 64)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	material, err := strconv.ParseFloat(params[3], 64)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	teeth, err := strconv.ParseFloat(params[4], 64)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	lid, err := strconv.ParseBool(params[5])
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	svg := laserbox.Do(width, height, depth, material, teeth, lid)
	w.Header().Set("Content-Type", "image/svg+xml")
	_, err = io.WriteString(w, svg)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/svg/", svgHandler)

	log.Println("serving...")
	http.ListenAndServe(":1234", nil)
}
