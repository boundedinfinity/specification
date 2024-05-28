package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"gopkg.in/yaml.v3"
)

func main() {
	rootDir := os.Args[1]
	countryDir := filepath.Join(rootDir, "country")
	countryMapFn := filepath.Join(countryDir, "country-map.yaml")
	iso3166DataFn := filepath.Join(rootDir, "iso-3166.yaml")
	fipsDataFn := filepath.Join(rootDir, "fips.yaml")
	gencDataFn := filepath.Join(rootDir, "genc.yaml")
	flagDataFn := filepath.Join(rootDir, "country-flags.yaml")
	iso639DataFn := filepath.Join(rootDir, "iso-639.yaml")

	fmt.Printf("Processing %v\n", countryDir)

	var iso639Data iso639Data
	if err := unmarshal(iso639DataFn, &iso639Data); err != nil {
		log.Fatal(err)
	}

	iso639map := map[string]iso639Record{}
	for _, record := range iso639Data.Records {
		if _, ok := iso639map[record.Set1]; !ok {
			iso639map[record.Set1] = record
		}
	}

	var iso3166Data iso3166Data
	if err := unmarshal(iso3166DataFn, &iso3166Data); err != nil {
		log.Fatal(err)
	}

	iso3166map := map[string]iso3166Record{}
	for _, record := range iso3166Data.Records {
		if _, ok := iso3166map[record.Alpha2]; !ok {
			iso3166map[record.Alpha2] = record
		}
	}

	var countryMapData countryMapData
	if err := unmarshal(countryMapFn, &countryMapData); err != nil {
		log.Fatal(err)
	}

	var fipsData fipsData
	if err := unmarshal(fipsDataFn, &fipsData); err != nil {
		log.Fatal(err)
	}

	fipsMap := map[string]fipsRecord{}
	for _, record := range fipsData.Records {
		if _, ok := fipsMap[record.Code]; !ok {
			fipsMap[record.Code] = record
		}
	}

	var gencData gencData
	if err := unmarshal(gencDataFn, &gencData); err != nil {
		log.Fatal(err)
	}

	gencMap := map[string]gencRecord{}
	for _, record := range gencData.Records {
		if _, ok := gencMap[record.Alpha2]; !ok {
			gencMap[record.Alpha2] = record
		}
	}

	var flagData flagData
	if err := unmarshal(flagDataFn, &flagData); err != nil {
		log.Fatal(err)
	}

	flagMap := map[string]flagRecord{}
	for _, record := range flagData.Records {
		if _, ok := flagMap[record.Alpha2]; !ok {
			flagMap[record.Alpha2] = record
		}
	}

	var countryData countryData

	for _, mapRecord := range countryMapData.Records {
		record := countryRecord{
			Name: map[string][]string{},
			Lang: []countryIso639{},
		}

		if specRecord, ok := iso3166map[mapRecord.Iso3166]; ok {
			mergeName(record.Name, specRecord.Name)
			record.Iso3166 = countryIso3166{
				Alpha2:  specRecord.Alpha2,
				Alpha3:  specRecord.Alpha3,
				Numeric: specRecord.Numeric,
			}
		}

		if specRecord, ok := fipsMap[mapRecord.Fips]; ok {
			mergeName(record.Name, specRecord.Name)
			record.Fips = countryFips{
				Code: specRecord.Code,
			}
		}

		if specRecord, ok := gencMap[mapRecord.Genc]; ok {
			mergeName(record.Name, specRecord.Name)
			record.Genc = countryGenc{
				Alpha2:  specRecord.Alpha2,
				Alpha3:  specRecord.Alpha3,
				Numeric: specRecord.Numeric,
			}
		}

		countryData.Records = append(countryData.Records, record)
	}

	if err := marshal(rootDir, "county", countryData); err != nil {
		log.Fatal(err)
	}
}

func mergeName(dst map[string][]string, src map[string][]string) {
	for lang := range src {
		if _, ok := dst[lang]; !ok {
			dst[lang] = []string{}
		}

		dst[lang] = append(dst[lang], src[lang]...)
		dst[lang] = slicer.Uniq(dst[lang]...)
	}
}

func marshal(rootDir string, name string, v any) error {
	yamlFn := filepath.Join(rootDir, fmt.Sprintf("%v.%v", name, "yaml"))
	jsonFn := filepath.Join(rootDir, fmt.Sprintf("%v.%v", name, "json"))

	bs, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	err = os.WriteFile(yamlFn, bs, os.FileMode(0755))
	if err != nil {
		return err
	}

	bs, err = json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(jsonFn, bs, os.FileMode(0755))
	if err != nil {
		return err
	}

	return nil
}

func unmarshal(fn string, v any) error {
	content, err := os.ReadFile(fn)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, v)

	return err
}

type nameList map[string][]string

type countryData struct {
	Records []countryRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type countryRecord struct {
	Name    nameList        `json:"name,omitempty" yaml:"name,omitempty"`
	Iso3166 countryIso3166  `json:"iso-3166,omitempty" yaml:"iso-3166,omitempty"`
	Fips    countryFips     `json:"fips,omitempty" yaml:"fips,omitempty"`
	Genc    countryGenc     `json:"genc,omitempty" yaml:"genc,omitempty"`
	Lang    []countryIso639 `json:"lang,omitempty" yaml:"lang,omitempty"`
	Flag    countryFlag     `json:"flag,omitempty" yaml:"flag,omitempty"`
}

type countryIso639 struct {
	Set1  string `json:"set-1,omitempty" yaml:"set-1,omitempty"`
	Set2t string `json:"set-2t,omitempty" yaml:"set-2t,omitempty"`
	Set2b string `json:"set-2b,omitempty" yaml:"set-2b,omitempty"`
	Set3  string `json:"set-3,omitempty" yaml:"set-3,omitempty"`
}

type countryIso3166 struct {
	Alpha2  string `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Alpha3  string `json:"alpha-3,omitempty" yaml:"alpha-3,omitempty"`
	Numeric int    `json:"numeric,omitempty" yaml:"numeric,omitempty"`
}

type countryFips struct {
	Code string `json:"code,omitempty" yaml:"code,omitempty"`
}

type countryFlag struct {
	Svg map[string]string `json:"svg,omitempty" yaml:"svg,omitempty"`
}

type countryGenc struct {
	Alpha2  string `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Alpha3  string `json:"alpha-3,omitempty" yaml:"alpha-3,omitempty"`
	Numeric int    `json:"numeric,omitempty" yaml:"numeric,omitempty"`
}

type countryMapData struct {
	Records []mapRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type mapRecord struct {
	Iso3166 string `json:"iso-3166,omitempty" yaml:"iso-3166,omitempty"`
	Fips    string `json:"fips,omitempty" yaml:"fips,omitempty"`
	Genc    string `json:"genc,omitempty" genc:"fips,omitempty"`
}

type iso3166Data struct {
	Records []iso3166Record `json:"records,omitempty" yaml:"records,omitempty"`
}

type iso3166Record struct {
	Name        nameList          `json:"name,omitempty" yaml:"name,omitempty"`
	Alpha2      string            `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Alpha3      string            `json:"alpha-3,omitempty" yaml:"alpha-3,omitempty"`
	Numeric     int               `json:"numeric,omitempty" yaml:"numeric,omitempty"`
	Independent bool              `json:"independent,omitempty" yaml:"independent,omitempty"`
	Lang        []string          `json:"lang,omitempty" yaml:"lang,omitempty"`
	Status      string            `json:"status,omitempty" yaml:"status,omitempty"`
	Territory   []string          `json:"territory,omitempty" yaml:"territory,omitempty"`
	Divisions   []iso3166Division `json:"divisions,omitempty" yaml:"divisions,omitempty"`
	Sources     []iso3166Source   `json:"sources,omitempty" yaml:"sources,omitempty"`
}

type iso3166Division struct {
	Code               string   `json:"code,omitempty" yaml:"code,omitempty"`
	Category           string   `json:"category,omitempty" yaml:"category,omitempty"`
	Parent             string   `json:"parent,omitempty" yaml:"parent,omitempty"`
	Lang               []string `json:"lang,omitempty" yaml:"lang,omitempty"`
	Name               nameList `json:"name,omitempty" yaml:"name,omitempty"`
	RomanizationSystem string   `json:"romanization-system,omitempty" yaml:"romanization-system,omitempty"`
}

type iso3166Source struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Url  string `json:"url,omitempty" yaml:"url,omitempty"`
}

type fipsData struct {
	Records []fipsRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type fipsRecord struct {
	Name         nameList       `json:"name,omitempty" yaml:"name,omitempty"`
	Code         string         `json:"code,omitempty" yaml:"code,omitempty"`
	FirstVersion int            `json:"fips-first-version,omitempty" yaml:"fips-first-version,omitempty"`
	LastVersion  int            `json:"fips-last-version,omitempty" yaml:"fips-last-version,omitempty"`
	Designation  nameList       `json:"designation,omitempty" yaml:"designation,omitempty"`
	Divisions    []fipsDivision `json:"divisions,omitempty" yaml:"divisions,omitempty"`
}

type fipsDivision struct {
	Name         nameList `json:"name,omitempty" yaml:"name,omitempty"`
	Code         string   `json:"code,omitempty" yaml:"code,omitempty"`
	FirstVersion int      `json:"fips-first-version,omitempty" yaml:"fips-first-version,omitempty"`
	LastVersion  int      `json:"fips-last-version,omitempty" yaml:"fips-last-version,omitempty"`
	Designation  nameList `json:"designation,omitempty" yaml:"designation,omitempty"`
}

type gencData struct {
	Records []gencRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type gencRecord struct {
	Name           nameList `json:"name,omitempty" yaml:"name,omitempty"`
	Alpha2         string   `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Alpha3         string   `json:"alpha-3,omitempty" yaml:"alpha-3,omitempty"`
	Numeric        int      `json:"numeric,omitempty" yaml:"numeric,omitempty"`
	AdditionalInfo string   `json:"additionalInfo,omitempty" yaml:"additionalInfo,omitempty"`
}

type flagData struct {
	Records []flagRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type flagRecord struct {
	Alpha2 string            `json:"alpha-2,omitempty" yaml:"alpha-2,omitempty"`
	Svg    map[string]string `json:"svg,omitempty" yaml:"svg,omitempty"`
}

type iso639Data struct {
	Records []iso639Record `json:"records,omitempty" yaml:"records,omitempty"`
}

type iso639Record struct {
	Name  nameList `json:"name,omitempty" yaml:"name,omitempty"`
	Set1  string   `json:"set-1,omitempty" yaml:"set-1,omitempty"`
	Set2t string   `json:"set-2t,omitempty" yaml:"set-2t,omitempty"`
	Set2b string   `json:"set-2b,omitempty" yaml:"set-2b,omitempty"`
	Set3  string   `json:"set-3,omitempty" yaml:"set-3,omitempty"`
}
