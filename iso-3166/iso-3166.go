package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	inputFn := os.Args[1]
	outputFn := os.Args[2]
	fmt.Printf("Processing %v\n", inputFn)

	file, err := os.Open(inputFn)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "=====") {
			break
		}

		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	if err := process1(lines); err != nil {
		log.Fatal(err)
	}

	lines = []string{}

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "=====") {
			if err = process2(lines); err != nil {
				log.Fatal(err)
			}

			lines = []string{}
		}

		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	var records []isoRecord

	for _, record := range alpha2Map {
		for lang, list := range record.Name {
			record.Name[lang] = uniq(list)
		}

		record.Lang = uniq(record.Lang)

		records = append(records, *record)
	}

	bs, err := yaml.Marshal(records)

	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(outputFn, bs, os.FileMode(0755)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processed %v records\n", len(records))
}

var alpha2Map = map[string]*isoRecord{}

type isoName map[string][]string

type isoRecord struct {
	Name        isoName       `json:"name" yaml:"name,omitempty"`
	Alpha2      string        `json:"alpha-2" yaml:"alpha-2,omitempty"`
	Alpha3      string        `json:"alpha-3" yaml:"alpha-3,omitempty"`
	Numeric     string        `json:"numeric" yaml:"numeric,omitempty"`
	Independent bool          `json:"independent" yaml:"independent,omitempty"`
	Lang        []string      `json:"lang" yaml:"lang,omitempty"`
	Status      string        `json:"status" yaml:"status,omitempty"`
	Territory   []string      `json:"territory" yaml:"territory,omitempty"`
	Divisions   []isoDivision `json:"divisions" yaml:"divisions,omitempty"`
	Sources     []isoSource   `json:"sources,omitempty" yaml:"sources,omitempty"`
}

func newIsoRecord() *isoRecord {
	return &isoRecord{
		Name: isoName{},
	}
}

type isoDivision struct {
	Code               string   `json:"code" yaml:"code,omitempty"`
	Category           string   `json:"category" yaml:"category,omitempty"`
	Parent             string   `json:"parent,omitempty" yaml:"parent,omitempty"`
	Lang               []string `json:"lang" yaml:"lang,omitempty"`
	Name               isoName  `json:"name" yaml:"name,omitempty"`
	RomanizationSystem string   `json:"romanization-system" yaml:"romanization-system,omitempty"`
}

type isoSource struct {
	Name string `json:"name" yaml:"name,omitempty"`
	Url  string `json:"url" yaml:"url,omitempty"`
}

func newIsoDivision() *isoDivision {
	return &isoDivision{
		Name: isoName{},
	}
}

func process1(lines []string) error {
	for _, line := range lines {
		fields := getFields(line)
		record := newIsoRecord()

		record.Name["en"] = normName(fields[0])
		record.Name["fr"] = normName(fields[1])
		record.Alpha2 = fields[2]
		record.Alpha3 = fields[3]
		record.Numeric = fields[4]

		if _, ok := alpha2Map[record.Alpha2]; !ok {
			alpha2Map[record.Alpha2] = record
		}
	}

	return nil
}

var errEOR = errors.New("end of record")

func process2(lines []string) error {
	index := 0

	moveTo := func(s string) error {
		for i := index; i < len(lines); i++ {
			if strings.Contains(lines[i], s) {
				index = i
				break
			}
		}

		return nil
	}

	captureTo := func(s string) ([]string, error) {
		var captured []string

		for i := index; i < len(lines); i++ {
			if strings.Contains(lines[i], s) {
				index = i
				break
			}

			captured = append(captured, lines[i])
		}

		return captured, nil
	}

	if err := moveTo("Alpha-2 code"); err != nil {
		return err
	}

	alpha2 := strings.ReplaceAll(lines[index], "Alpha-2 code", "")
	record, ok := alpha2Map[alpha2]

	if !ok {
		return fmt.Errorf("%v not found", alpha2)
	}

	if err := moveTo("Independent"); err != nil {
		return err
	}

	if strings.Contains(lines[index], "Yes") {
		record.Independent = true
	}

	if err := moveTo("Territory name"); err != nil {
		return err
	}

	territoryLine := strings.TrimSpace(strings.ReplaceAll(lines[index], "Territory name", ""))

	if territoryLine != "" {
		for _, territory := range strings.Split(territoryLine, ",") {
			record.Territory = append(record.Territory, normStr(territory))
		}
	}

	if err := moveTo("Status"); err != nil {
		return err
	}

	record.Status = strings.ReplaceAll(lines[index], "Status", "")

	if err := moveTo("Administrative language"); err != nil {
		return err
	}

	index += 1
	langLines, err := captureTo("Subdivisions")

	if err != nil {
		return err
	}

	for _, langLine := range langLines {
		fields := getFields(langLine)
		if fields[0] != "-" {
			record.Lang = append(record.Lang, fields[0])
			record.Name[fields[0]] = append(record.Name[fields[0]], fields[2])
		}
	}

	if err := moveTo("List source"); err != nil {
		return err
	}

	index += 1
	sourcesLines, err := captureTo("Code source")

	if err != nil {
		return err
	}

	regex := regexp.MustCompile(`(.*) \((.*)\)`)

	for _, line1 := range sourcesLines {
		if line1 != "" {
			sourcesItems := strings.Split(line1, ";")

			for _, sourceItem := range sourcesItems {
				sourceItem = strings.TrimSpace(sourceItem)

				if sourceItem == "" {
					continue
				}

				found := regex.FindStringSubmatch(sourceItem)
				var source isoSource

				if len(found) > 1 && strings.HasPrefix(found[2], "http") {
					if !strings.Contains(found[1], "jsessionid") {
						source.Name = found[1]
					}

					source.Url = found[2]
				} else if strings.HasPrefix(sourceItem, "http") {
					source.Url = sourceItem
				} else {
					if !strings.Contains(sourceItem, "jsessionid") {
						source.Name = sourceItem
					}
				}

				if source.Name != "" || source.Url != "" {
					record.Sources = append(record.Sources, source)
				}
			}
		}
	}

	if err := moveTo("Subdivision category"); err != nil {
		return err
	}

	index += 1
	divisionLines, err := captureTo("Change history of country code")

	if err != nil {
		return err
	}

	divisionMap := map[string]*isoDivision{}

	for _, divisionLine := range divisionLines {
		fields := strings.Split(divisionLine, "\t")
		category := normStr(fields[0])
		code := normStr(fields[1])
		name := normStr(fields[2])
		localVariant := normStr(fields[3])
		lang := normStr(fields[4])
		romanizationSystem := normStr(fields[5])
		parent := normStr(fields[6])

		division, ok := divisionMap[code]

		if !ok {
			division = newIsoDivision()
			divisionMap[code] = division
		}

		division.Category = category
		division.Code = code

		if lang != "-" {
			division.Name[lang] = append(division.Name[lang], name)

			if localVariant != "" {
				division.Name[lang] = append(division.Name[lang], localVariant)
			}
		}

		division.RomanizationSystem = romanizationSystem
		division.Parent = parent
	}

	for _, division := range divisionMap {
		record.Divisions = append(record.Divisions, *division)
	}

	if err := moveTo("========"); err != nil {
		return err
	}

	fmt.Printf("Processed: %v %v, (divisions: %v)\n", record.Alpha2, record.Name["en"][0], len(record.Divisions))
	return nil
}

func normStr(s string) string {
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "†", "")
	s = strings.TrimSpace(s)
	return s
}

func normName(s string) []string {
	output := []string{}
	output = append(output, s)
	regex := regexp.MustCompile(`(.*) \((.*)\)`)
	found := regex.FindStringSubmatch(s)

	if len(found) > 1 {
		var name string

		if strings.HasSuffix(found[2], "'") {
			name = fmt.Sprintf("%v%v", found[2], found[1])
		} else {
			name = fmt.Sprintf("%v %v", found[2], found[1])
		}

		output = append(output, normStr(name))
	}

	return output
}

func getFields(s string) []string {
	fields := strings.Split(s, "\t")
	var output []string

	for _, field := range fields {
		if field != "\t" && field != "" {
			output = append(output, field)
		}
	}

	return output
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
