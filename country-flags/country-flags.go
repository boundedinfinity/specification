package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/boundedinfinity/go-commoner/idiomatic/extentioner"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"gopkg.in/yaml.v3"
)

func main() {
	outputDir := os.Args[1]
	inputDirs := os.Args[2:]
	m := map[string]*flagRecord{}

	fmt.Printf("Processing %v\n", inputDirs)

	for _, inputDir := range inputDirs {
		entries, err := os.ReadDir(inputDir)

		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			svgName := entry.Name()
			alpha2 := extentioner.Strip(svgName)

			if stringer.Contains(alpha2, "-") {
				continue
			}

			record, ok := m[alpha2]

			if !ok {
				record = &flagRecord{
					Alpha2: alpha2,
					Svg:    map[string]string{},
				}

				m[alpha2] = record
			}

			svgBs, err := os.ReadFile(filepath.Join(inputDir, svgName))

			if err != nil {
				log.Fatal(err)
			}

			svgDim := filepath.Base(inputDir)
			record.Svg[svgDim] = string(svgBs)
		}
	}

	outputYamlFn := filepath.Join(outputDir, "country-flags.yaml")
	outputJsonFn := filepath.Join(outputDir, "country-flags.json")

	var data flagData

	for _, record := range m {
		data.Records = append(data.Records, *record)
	}

	bs, err := yaml.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(outputYamlFn, bs, os.FileMode(0755)); err != nil {
		log.Fatal(err)
	}

	bs, err = json.MarshalIndent(data, "", "    ")

	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(outputJsonFn, bs, os.FileMode(0755)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processed %v records\n", len(data.Records))
}

type flagData struct {
	Records []flagRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type flagRecord struct {
	Alpha2 string            `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Svg    map[string]string `json:"svg,omitempty" yaml:"svg,omitempty"`
}
