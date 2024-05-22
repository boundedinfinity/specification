package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	var data isoData

	for _, country := range baseline.Country {
		numericCode, err := strconv.Atoi(country.Encoding.NumericCode)

		if err != nil {
			log.Fatal(err)
		}

		record := isoRecord{
			Name: isoName{
				"en": []string{
					country.ShortName,
					country.FullName,
				},
			},
			Alpha2:  country.Encoding.Char2Code,
			Alpha3:  country.Encoding.Char3Code,
			Numeric: numericCode,
		}

		for _, localShortName := range country.LocalShortName {
			lang := localShortName.NameLanguage2Char

			if _, ok := record.Name[lang]; !ok {
				record.Name[lang] = []string{}
			}

			record.Name[lang] = append(record.Name[lang], localShortName.Name)
		}

		for lang, list := range record.Name {
			record.Name[lang] = uniq(list)
		}

		data.Records = append(data.Records, record)
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

type isoData struct {
	Records []isoRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type isoRecord struct {
	Name    isoName `json:"name,omitempty" yaml:"name,omitempty"`
	Alpha2  string  `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Alpha3  string  `json:"alpha-3,omitempty" yaml:"alpha-3,omitempty"`
	Numeric int     `json:"numeric,omitempty" yaml:"numeric,omitempty"`
}

type isoName map[string][]string

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

func uniq(vs []string) []string {
	o := []string{}
	m := map[string]bool{}

	for _, v := range vs {
		if _, ok := m[v]; !ok {
			m[v] = true
		}
	}

	for v := range m {
		v = strings.ReplaceAll(v, "*", "")
		o = append(o, v)
	}

	return o
}
