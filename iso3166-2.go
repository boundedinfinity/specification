package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type tableData struct {
	category string
	code     string
	name     string
	local    string
	lang     string
	system   string
	parent   string
}

type isoRecord struct {
	SubDivisions []isoSubDivision `json:"subdivisions" yaml:"subdivisions"`
}

type isoSubDivision struct {
	Code     string              `json:"code" yaml:"code"`
	Category string              `json:"category" yaml:"category"`
	Parent   string              `json:"parent,omitempty" yaml:"parent,omitempty"`
	Lang     []string            `json:"lang" yaml:"lang"`
	Name     map[string][]string `json:"name" yaml:"name"`
}

func get(index int, fields []string, v *string) {
	if len(fields) > index {
		*v = fields[index]
	}
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

func main() {
	fn := os.Args[1]
	fmt.Printf("Processing %v\n", fn)

	content, err := os.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	fieldMap := map[string]*isoSubDivision{}

	for _, line := range lines {
		td := tableData{}
		fields := strings.Split(line, "\t")

		get(0, fields, &td.category)
		get(1, fields, &td.code)
		get(2, fields, &td.name)
		get(3, fields, &td.local)
		get(4, fields, &td.lang)
		get(5, fields, &td.system)
		get(6, fields, &td.parent)

		iso := isoSubDivision{
			Code:     td.code,
			Category: td.category,
			Parent:   td.parent,
			Lang:     []string{td.lang},
			Name: map[string][]string{
				"en": {td.name},
			},
		}

		if iso.Code == "" {
			continue
		}

		if _, ok := fieldMap[iso.Code]; !ok {
			fieldMap[iso.Code] = &iso
		} else {
			fieldMap[iso.Code].Lang = append(fieldMap[iso.Code].Lang, iso.Lang...)
			for lang := range fieldMap[iso.Code].Name {
				if _, ok := iso.Name[lang]; ok {
					fieldMap[iso.Code].Name[lang] = append(fieldMap[iso.Code].Name[lang], iso.Name[lang]...)
				}
			}
		}
	}

	isoDivisions := []isoSubDivision{}

	for _, v := range fieldMap {
		v.Code = strings.ReplaceAll(v.Code, "*", "")
		v.Lang = uniq(v.Lang)

		for lang := range v.Name {
			v.Name[lang] = uniq(v.Name[lang])
		}

		isoDivisions = append(isoDivisions, *v)
	}

	buf := &bytes.Buffer{}
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(4)

	records := isoRecord{
		SubDivisions: isoDivisions,
	}

	err = enc.Encode(records)

	if err != nil {
		log.Fatal(err)
	}

	rawY := buf.String()

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(rawY)
	fmt.Println(strings.Repeat("-", 50))
}
