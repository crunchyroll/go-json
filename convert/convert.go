// convert takes a JSON object from a file as input and produces as output a set of Go structs
// into which it can be unmarshalled
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"./converter"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage:", os.Args[0], "<input> <package_name>")
		return
	}

	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error reading input file %q: %v\n", os.Args[1], err)
	}

	pkgName := strings.ToLower(os.Args[2])

	c, err := converter.Convert(strings.Title(pkgName), input)
	if err != nil {
		fmt.Printf("error converting JSON: %v\n", err)
		return
	}

	_, fn := filepath.Split(os.Args[1])
	os.Mkdir(pkgName, os.ModePerm)
	f, err := os.Create(fmt.Sprintf("%s/%s.go", pkgName, fn))
	if err != nil {
		fmt.Printf("unable to create package: %v\n", err)
		return
	}

	fmt.Fprintf(f, `// Package %s is an auto-generated schema for input JSON. Do not edit it by hand or check it in to
// source control.
package %[1]s
`, pkgName)

	c.WriteStructs(f)
}
