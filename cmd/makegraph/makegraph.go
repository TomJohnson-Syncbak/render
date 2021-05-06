package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tmc/dot"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [file]\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	args := flag.Args()

	if len(args) != 1 {
		usage()
		os.Exit(1)
	}

	g := dot.NewGraph("G")
	g.Set("label", "Example graph")

	jsonFile := args[0]

	f, err := os.Open(jsonFile)
	if err != nil {
		panic(err)
	}

	var nodeCount int
	var subgraphCount int

	var node1 *dot.Node
	var node2 *dot.Node

	// One line contains exactly one segment

	input := bufio.NewScanner(f)
	for input.Scan() {

		jsonInput := input.Bytes()

		var seg Segment

		json.Unmarshal(jsonInput, &seg)

		var nodeName = strconv.Itoa(seg.SequenceNumber)
		node1 = dot.NewNode(nodeName)
		nodeCount++
		node1.Set("rank", strconv.Itoa(nodeCount))
		g.AddNode(node1)
		g.Set("ranksep", "0.5 equally")

		if node2 != nil {
			edge := dot.NewEdge(node2, node1)
			g.AddEdge(edge)
		}
		node2 = node1

		var last *dot.Node
		var sg = dot.NewSubgraph("History_" + nodeName)
		g.AddSubgraph(sg)
		subgraphCount++
		sg.Set("rank", strconv.Itoa(nodeCount))

		for _, h := range seg.History {
			var hNode = dot.NewNode(nodeName + "_" + h.Message)
			nodeCount++
			if last != nil {
				var hEdge = dot.NewEdge(last, hNode)
				sg.AddEdge(hEdge)
			}
			last = hNode
		}

	}

	fmt.Println(g)

	fmt.Fprintf(os.Stderr, "Node count: %d\n", nodeCount)
	fmt.Fprintf(os.Stderr, "Subgraph count: %d\n", subgraphCount)

	outname := strings.Replace(jsonFile, ".json", ".gv", 1)
	err = ioutil.WriteFile(outname, []byte(g.String()), os.FileMode(os.O_CREATE|os.O_RDWR|0644))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file %s due to %s\n", outname, err)
	} else {
		fmt.Fprintf(os.Stderr, "Wrote file %s\n", outname)
	}
}

//Segment whatever
type Segment struct {
	RawStreamID    string        `json:"rawStreamID"`
	SourceStreamID string        `json:"sourceStreamID"`
	SequenceNumber int           `json:"sequenceNumber"`
	StartTime      time.Time     `json:"startTime"`
	Duration       time.Duration `json:"duration"`
	History        []Entry       `json:"history"`
}

//Entry whatever
type Entry struct {
	Fields struct {
		Age            int64     `json:"age"`
		Duration       int64     `json:"duration"`
		Host           string    `json:"host"`
		Name           string    `json:"name"`
		Queue          string    `json:"queue"`
		SegmentLength  int       `json:"segmentLength"`
		SequenceNumber int       `json:"sequenceNumber"`
		RawStreamID    string    `json:"rawStreamID"`
		SourceStreamID string    `json:"sourceStreamID"`
		SourceStreams  []string  `json:"sourceStreams"`
		Error          string    `json:"error"`
		StartTime      time.Time `json:"startTime"`
		Version        string    `json:"version"`
	} `json:"fields"`
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}
