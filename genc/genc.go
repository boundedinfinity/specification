package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"gopkg.in/yaml.v1"
)

// https://gobyexample.com/xml

func main() {
	inputFn := os.Args[1]
	outputDir := os.Args[2]
	outputYamlFn := filepath.Join(outputDir, "genc.yaml")
	outputJsonFn := filepath.Join(outputDir, "genc.json")

	fmt.Printf("Processing %v\n", inputFn)

	content, err := os.ReadFile(inputFn)
	if err != nil {
		log.Fatal(err)
	}

	var baseline GencStandardBaseline
	if err = xml.Unmarshal(content, &baseline); err != nil {
		log.Fatal(err)
	}

	var data gencData

	for _, country := range baseline.Country {
		numericCode, err := strconv.Atoi(country.Encoding.NumericCode)

		if err != nil {
			log.Fatal(err)
		}

		record := gencRecord{
			Name: gencName{
				"en": []string{
					country.ShortName,
					country.FullName,
				},
			},
			Alpha2:         country.Encoding.Char2Code,
			Alpha3:         country.Encoding.Char3Code,
			Numeric:        numericCode,
			AdditionalInfo: country.AdditionalInfo,
		}

		for _, localShortName := range country.LocalShortName {
			lang := localShortName.NameLanguage2Char

			if _, ok := record.Name[lang]; !ok {
				record.Name[lang] = []string{}
			}

			record.Name[lang] = append(record.Name[lang], localShortName.Name)
		}

		for lang, list := range record.Name {
			record.Name[lang] = slicer.Uniq(list...)
			record.Name[lang] = slicer.Sort(list...)
		}

		data.Records = append(data.Records, record)
	}

	data.Records = slicer.SortFn(func(r gencRecord) string { return r.Alpha2 }, data.Records...)

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

type gencData struct {
	Records []gencRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type gencRecord struct {
	Name           gencName `json:"name,omitempty" yaml:"name,omitempty"`
	Alpha2         string   `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Alpha3         string   `json:"alpha-3,omitempty" yaml:"alpha-3,omitempty"`
	Numeric        int      `json:"numeric,omitempty" yaml:"numeric,omitempty"`
	AdditionalInfo string   `json:"additionalInfo,omitempty" yaml:"additionalInfo,omitempty"`
}

type gencName map[string][]string

type GencStandardBaseline struct {
	XMLName          xml.Name              `xml:"GENCStandardBaseline"`
	Authority        string                `xml:"authority"`
	Baseline         string                `xml:"baseline"`
	PromulgationDate string                `xml:"promulgationDate"`
	Country          []GencISOCountryEntry `xml:"ISOCountryEntry"`
}

type GencISOCountryEntry struct {
	XMLName                  xml.Name             `xml:"ISOCountryEntry"`
	Encoding                 GencEncoding         `xml:"encoding"`
	Char3CodeElementStatus   string               `xml:"char3CodeElementStatus"`
	Char2CodeElementStatus   string               `xml:"char2CodeElementStatus"`
	NumericCodeElementStatus string               `xml:"numericCodeElementStatus"`
	Name                     string               `xml:"name"`
	ShortName                string               `xml:"shortName"`
	FullName                 string               `xml:"fullName"`
	UnLegalStatus            string               `xml:"unLegalStatus"`
	EntryDate                string               `xml:"entryDate"`
	LocalShortName           []GencLocalShortName `xml:"localShortName"`
	BgnShortNameVariance     string               `xml:"bgnShortNameVariance"`
	BgnFullNameVariance      string               `xml:"bgnFullNameVariance"`
	UsRecognition            string               `xml:"usRecognition"`
	AdditionalInfo           string               `xml:"gencAdditionalInfo"`
}

type GencEncoding struct {
	XMLName           xml.Name           `xml:"encoding"`
	Char3Code         string             `xml:"char3Code"`
	Char3CodeURISet   GencCharCodeURISet `xml:"char3CodeURISet"`
	Char2Code         string             `xml:"char2Code"`
	Char2CodeURISet   GencCharCodeURISet `xml:"char2CodeURISet"`
	NumericCode       string             `xml:"numericCode"`
	NumericCodeURISet GencCharCodeURISet `xml:"numericCodeURISet"`
}

type GencCharCodeURISet struct {
	CodespaceURL           string `xml:"codespaceURL"`
	CodespaceURN           string `xml:"codespaceURN"`
	CodespaceURNBased      string `xml:"codespaceURNBased"`
	CodespaceURNBasedShort string `xml:"codespaceURNBasedShort"`
}

type GencLocalShortName struct {
	NameLanguage2Char string `xml:"nameLanguage2Char"`
	NameLanguage3Char string `xml:"nameLanguage3Char"`
	Name              string `xml:"name"`
	Iso6393Char3Code  string `xml:"iso6393Char3Code"`
}
