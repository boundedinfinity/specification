package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"gopkg.in/yaml.v3"
)

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

	outputYamlFn := outputPath(".yaml")
	outputJsonFn := outputPath(".json")

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

	var data isoData

	for _, record := range alpha2Map {
		for lang, list := range record.Name {
			record.Name[lang] = slicer.Uniq(list...)
		}

		record.Lang = slicer.Uniq(record.Lang...)
		record.Lang = slicer.Sort(record.Lang...)
		record.Divisions = slicer.SortFn(func(r isoDivision) string { return r.Code }, record.Divisions...)

		data.Records = append(data.Records, *record)
	}

	data.Records = slicer.SortFn(func(r isoRecord) string { return r.Alpha2 }, data.Records...)

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

var alpha2Map = map[string]*isoRecord{}

type isoData struct {
	Records []isoRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type isoName map[string][]string

type isoRecord struct {
	Name        isoName       `json:"name,omitempty" yaml:"name,omitempty"`
	Alpha2      string        `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Alpha3      string        `json:"alpha-3,omitempty" yaml:"alpha-3,omitempty"`
	Numeric     int           `json:"numeric,omitempty" yaml:"numeric,omitempty"`
	Independent bool          `json:"independent,omitempty" yaml:"independent,omitempty"`
	Lang        []string      `json:"lang,omitempty" yaml:"lang,omitempty"`
	Status      string        `json:"status,omitempty" yaml:"status,omitempty"`
	Territory   []string      `json:"territory,omitempty" yaml:"territory,omitempty"`
	Divisions   []isoDivision `json:"divisions,omitempty" yaml:"divisions,omitempty"`
	Sources     []isoSource   `json:"sources,omitempty" yaml:"sources,omitempty"`
}

func newIsoRecord() *isoRecord {
	return &isoRecord{
		Name: isoName{},
	}
}

type isoDivision struct {
	Code               string   `json:"code,omitempty" yaml:"code,omitempty"`
	Category           string   `json:"category,omitempty" yaml:"category,omitempty"`
	Parent             string   `json:"parent,omitempty" yaml:"parent,omitempty"`
	Lang               []string `json:"lang,omitempty" yaml:"lang,omitempty"`
	Name               isoName  `json:"name,omitempty" yaml:"name,omitempty"`
	RomanizationSystem string   `json:"romanization-system,omitempty" yaml:"romanization-system,omitempty"`
}

type isoSource struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Url  string `json:"url,omitempty" yaml:"url,omitempty"`
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
		record.Alpha2 = normStr(fields[2])
		record.Alpha3 = normStr(fields[3])

		numericCode, err := strconv.Atoi(normStr(fields[4]))

		if err != nil {
			return err
		}

		record.Numeric = numericCode

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

	// fmt.Printf("Processed: %v %v, (divisions: %v)\n", record.Alpha2, record.Name["en"][0], len(record.Divisions))
	return nil
}

var normStrReplacer = strings.NewReplacer(
	"*", "",
	",", "",
	"†", "",
)

func normStr(s string) string {
	s = normStrReplacer.Replace(s)
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
