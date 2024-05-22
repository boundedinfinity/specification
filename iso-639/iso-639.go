package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

type isoData struct {
	Records []isoRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type isoRecord struct {
	Name  isoName `json:"name,omitempty" yaml:"name,omitempty"`
	Set1  string  `json:"set-1,omitempty" yaml:"set-1,omitempty"`
	Set2t string  `json:"set-2t,omitempty" yaml:"set-2t,omitempty"`
	Set2b string  `json:"set-2b,omitempty" yaml:"set-2b,omitempty"`
	Set3  string  `json:"set-3,omitempty" yaml:"set-3,omitempty"`
}

type isoName map[string][]string

func main() {
	inputFn := os.Args[1]
	outputDir := os.Args[2]

	fmt.Printf("Processing %v\n", inputFn)

	outputPath := func(ext string) string {
		path := filepath.Base(inputFn)
		path = strings.Replace(path, filepath.Ext(path), ext, -1)
		path = filepath.Join(outputDir, path)
		return path
	}

	outputJsonFn := outputPath(".json")

	content, err := os.ReadFile(inputFn)
	if err != nil {
		log.Fatal(err)
	}

	var data isoData

	if err := yaml.Unmarshal(content, &data); err != nil {
		log.Fatal(err)
	}

	bs, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(outputJsonFn, bs, os.FileMode(0755)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processed %v records\n", len(data.Records))
}
