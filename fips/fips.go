package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v1"
)

func main() {
	inputFn := os.Args[1]
	outputDir := os.Args[2]
	outputYamlFn := filepath.Join(outputDir, "fips.yaml")
	outputJsonFn := filepath.Join(outputDir, "fips.json")

	fmt.Printf("Processing %v\n", inputFn)

	content, err := os.ReadFile(inputFn)
	if err != nil {
		log.Fatal(err)
	}

	var records []*fipsRecord
	var parent *fipsRecord

	for _, line := range strings.Split(string(content), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		linedata, err := parseLine(line)
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasSuffix(linedata.Code, "00") {
			parent = &fipsRecord{
				Name:         linedata.Name,
				Code:         linedata.Code,
				FirstVersion: linedata.FirstVersion,
				LastVersion:  linedata.LastVersion,
				Designation:  linedata.Designation,
				Divisions:    []fipsDivision{},
			}

			records = append(records, parent)
		} else {
			parent.Divisions = append(parent.Divisions, fipsDivision(linedata))
		}
	}

	var data fipsData

	for _, record := range records {
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

func parseLine(line string) (fipsLine, error) {
	var data fipsLine

	comps := strings.Split(line, "____")

	if len(comps) == 1 {
		comps = strings.Split(line, "___")
	}

	if len(comps) == 1 {
		comps = strings.Split(line, "__")
	}

	nameStr := comps[1]
	nameStr = strings.Trim(nameStr, "\r\n")
	var names []string

	if strings.Contains(nameStr, "__") {
		names = strings.Split(nameStr, "__")
	} else {
		names = strings.Split(nameStr, "_")
	}

	data.Name = fipsName{
		"en": []string{names[0]},
	}

	other := strings.Split(comps[0], "_")
	data.Code = other[0]

	firstStr := other[1]
	lastStr := other[2]

	if first, err := strconv.Atoi(firstStr); err != nil {
		return data, err
	} else {
		data.FirstVersion = first
	}

	if last, err := strconv.Atoi(lastStr); err != nil {
		return data, err
	} else {
		data.LastVersion = last
	}

	if len(other) > 3 {
		designation := other[3]
		data.Designation = fipsName{
			"en": []string{designation},
		}
	}

	return data, nil
}

type fipsData struct {
	Records []fipsRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type fipsName map[string][]string

type fipsLine struct {
	Name         fipsName
	Code         string
	FirstVersion int
	LastVersion  int
	Designation  fipsName
}

type fipsRecord struct {
	Name         fipsName       `json:"name,omitempty" yaml:"name,omitempty"`
	Code         string         `json:"code,omitempty" yaml:"code,omitempty"`
	FirstVersion int            `json:"fips-first-version,omitempty" yaml:"fips-first-version,omitempty"`
	LastVersion  int            `json:"fips-last-version,omitempty" yaml:"fips-last-version,omitempty"`
	Designation  fipsName       `json:"designation,omitempty" yaml:"designation,omitempty"`
	Divisions    []fipsDivision `json:"divisions,omitempty" yaml:"divisions,omitempty"`
}

type fipsDivision struct {
	Name         fipsName `json:"name,omitempty" yaml:"name,omitempty"`
	Code         string   `json:"code,omitempty" yaml:"code,omitempty"`
	FirstVersion int      `json:"fips-first-version,omitempty" yaml:"fips-first-version,omitempty"`
	LastVersion  int      `json:"fips-last-version,omitempty" yaml:"fips-last-version,omitempty"`
	Designation  fipsName `json:"designation,omitempty" yaml:"designation,omitempty"`
}
