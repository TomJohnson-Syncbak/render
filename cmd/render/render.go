// Given the name of a file in Graphviz 'dot' format, render it as SVG
// and open the SVG file in the browser.
// Assumes that Graphviz' 'dot' command is installed on the local machine and is
// in the user's path.
// If the data comes from standard input, 'dot' is not run and the bytes from
// standard input are sent directly to the browser
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/browser"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [file]\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var outType string
	flag.StringVar(&outType, "t", "svg", "Type of output file, e.g. svg, png")

	var progType string
	flag.StringVar(&progType, "p", "dot", "Program to run, e.g. dot, neato")

	flag.Parse()
	args := flag.Args()
	switch len(args) {
	case 1:
		log.Printf("Reading from file: %s", args[0])
		check(render(args[0], progType, outType))
	default:
		usage()
	}
}

func render(dotFile string, progType string, outType string) error {
	var e error

	var extension = path.Ext(dotFile)
	var noExt = strings.TrimSuffix(dotFile, extension)
	var outFile = noExt + "." + outType

	log.Printf("Running "+progType+" on %s\n", dotFile)
	outBytes, e := exec.Command(progType, "-T"+outType, dotFile).Output()
	check(e)

	log.Printf("Writing %s\n", outFile)
	check(ioutil.WriteFile(outFile, outBytes, 0755))

	log.Printf("Opening browser on %s\n", outFile)
	e = browser.OpenFile(outFile)
	return e
}
