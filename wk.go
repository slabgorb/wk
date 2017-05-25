package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/slabgorb/wk/workflow"
)

var filename *string

func init() {
	filename = flag.String("file", "", "path to workflow file")
	flag.Parse()
}

func main() {
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	file, err := workflow.LoadFile(*filename)
	if err != nil {
		panic(err)
	}
	steps, err := workflow.ParseSteps(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("done parsing")
	for _, step := range steps.List {
		fmt.Printf("%s, %s\n", step.Command, strings.Trim(fmt.Sprintf("%v", step.Arguments), "map[]"))
	}
	b, err := steps.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%d bytes", b)

}
