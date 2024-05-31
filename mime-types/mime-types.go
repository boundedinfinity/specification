package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/boundedinfinity/go-commoner/idiomatic/extentioner"
	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

var (
	cacheDir             = "cache"
	debug                = false
	ianaMediaTypesUrl, _ = url.Parse("https://www.iana.org/assignments/media-types")
)

func main() {
	outputDir := os.Args[1]

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourcesFn := filepath.Join(wd, "iana-media-types-sources.yaml")
	fmt.Printf("Processing %v\n", sourcesFn)

	recordsMap := map[string]*mimeRecord{}

	var sourcesData ianaMediaTypesSources
	if err := unmarshal(sourcesFn, &sourcesData); err != nil {
		panic(err)
	}

	var csvSourceContents []string
	for _, source := range sourcesData.MimeTypes {
		csvSourceContent, err := downloadFile(wd, "", source)
		if err != nil {
			panic(err)
		}

		csvSourceContents = append(csvSourceContents, csvSourceContent)
	}

	regex := regexp.MustCompile(`\[(.*?)\]`)

	for _, sourceContent := range csvSourceContents {
		reader := csv.NewReader(strings.NewReader(sourceContent))
		csvRecords, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}

		for _, csvRecord := range csvRecords {
			name := csvRecord[0]
			template := csvRecord[1]
			reference := csvRecord[2]

			if name == "Name" {
				continue
			}

			record := mimeRecord{
				MimeType:   template,
				References: []mimeReference{},
			}

			if strings.HasPrefix(reference, "[") {
				matches := regex.FindAllStringSubmatch(reference, -1)

				for _, match := range matches {
					if strings.HasPrefix(match[1], "RFC") {
						record.References = append(record.References, mimeReference{
							Name: match[1],
						})
					}
				}
			}

			recordsMap[record.MimeType] = &record
		}
	}

	for _, source := range sourcesData.FileExtentions {
		_, err := downloadFile(wd, "", source)
		if err != nil {
			panic(err)
		}

	}

	nonIanaDataFn := filepath.Join(wd, "non-iana.yaml")
	fmt.Printf("Processing %v\n", nonIanaDataFn)

	var nonIanaDatas []mimeRecord
	if err := unmarshal(nonIanaDataFn, &nonIanaDatas); err != nil {
		panic(err)
	}

	for _, nonIanaData := range nonIanaDatas {
		record, ok := recordsMap[nonIanaData.MimeType]

		if !ok {
			record = &mimeRecord{
				Description: nonIanaData.Description,
			}

			recordsMap[nonIanaData.MimeType] = record
		}

		record.Description = nonIanaData.Description
		record.MimeTypeAlternative = append(record.MimeTypeAlternative, nonIanaData.MimeTypeAlternative...)
		record.FileExtentions = append(record.FileExtentions, nonIanaData.FileExtentions...)
	}

	var data mimeData
	for _, record := range recordsMap {
		record.MimeTypeAlternative = slicer.Uniq(record.MimeTypeAlternative...)
		data.Records = append(data.Records, *record)

		detailsContent, err := getDetails(wd, record)
		if err != nil {
			panic(err)
		}

		fmt.Println(detailsContent)
	}

	if err := marshal(filepath.Join(outputDir, "mime-type.yaml"), data); err != nil {
		panic(err)
	}

	if debug {
		spew.Dump(data)
	}

	fmt.Printf("Records: %v", len(data.Records))
}

type mimeData struct {
	Records []mimeRecord `json:"records,omitempty" yaml:"records,omitempty"`
}

type mimeRecord struct {
	Description                    string          `json:"description,omitempty" yaml:"description,omitempty"`
	MimeType                       string          `json:"mime-type,omitempty" yaml:"mime-type,omitempty"`
	MimeTypeAlternative            []string        `json:"mime-type-alternative,omitempty" yaml:"mime-type-alternative,omitempty"`
	References                     []mimeReference `json:"references,omitempty" yaml:"references,omitempty"`
	FileExtentions                 []string        `json:"file-extentions,omitempty" yaml:"file-extentions,omitempty"`
	EncodingConsiderations         string          `json:"encoding-considerations,omitempty" yaml:"encoding-considerations,omitempty"`
	SecurityConsiderations         string          `json:"security-considerations,omitempty" yaml:"security-considerations,omitempty"`
	InteroperabilityConsiderations string          `json:"interoperability-considerations,omitempty" yaml:"interoperability-considerations,omitempty"`
	PublishedSpecification         string          `json:"published-specification,omitempty" yaml:"published-specification,omitempty"`
}

type mimeReference struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Url  string `json:"url,omitempty" yaml:"url,omitempty"`
}

type ianaMediaTypesSources struct {
	MimeTypes      []string `json:"mime-types,omitempty" yaml:"mime-types,omitempty"`
	FileExtentions []string `json:"file-extentions,omitempty" yaml:"file-extentions,omitempty"`
}

type ianaDetails struct {
	FileExtentions                 string `json:"file-extentions,omitempty" yaml:"file-extentions,omitempty"`
	EncodingConsiderations         string `json:"encoding-considerations,omitempty" yaml:"encoding-considerations,omitempty"`
	SecurityConsiderations         string `json:"security-considerations,omitempty" yaml:"security-considerations,omitempty"`
	InteroperabilityConsiderations string `json:"interoperability-considerations,omitempty" yaml:"interoperability-considerations,omitempty"`
	PublishedSpecification         string `json:"published-specification,omitempty" yaml:"published-specification,omitempty"`
}

func getDetails(wd string, record *mimeRecord) (string, error) {
	if record.MimeType == "" {
		return "", nil
	}

	detailsUrl, err := url.Parse(ianaMediaTypesUrl.String())

	if err != nil {
		return "", err
	}

	detailsUrl.Path = path.Join(detailsUrl.Path, record.MimeType)

	return downloadFile(wd, record.MimeType, detailsUrl.String())
}

func downloadFile(wd, relPath, source string) (string, error) {
	var content string

	fn, err := cacheName(wd, relPath, source)
	if err != nil {
		return content, err
	}

	if !pather.Files.Exists(fn) {
		resp, err := http.Get(source)
		if err != nil {
			return content, err
		}

		defer resp.Body.Close()

		of, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, fs.FileMode(0755))
		if err != nil {
			return content, err
		}

		defer of.Close()

		_, err = io.Copy(of, resp.Body)
		if err != nil {
			return content, err
		}
	}

	bs, err := os.ReadFile(fn)
	if err != nil {
		return content, err
	}

	content = string(bs)

	return content, nil
}

func cacheName(wd, relpath, source string) (string, error) {
	var fn string

	if relpath != "" {
		fn = path.Join(wd, cacheDir, relpath)
	} else {
		fn = source

		if strings.HasPrefix(fn, "http") {
			parsed, err := url.Parse(source)
			if err != nil {
				return fn, err
			}

			fn = parsed.Path
		}

		if strings.Contains(fn, "/") {
			fn = filepath.Base(fn)
		}

		fn = filepath.Join(wd, cacheDir, fn)
	}

	if err := os.MkdirAll(filepath.Dir(fn), os.FileMode(0755)); err != nil {
		return fn, err
	}

	return fn, nil
}

func marshal(fn string, v any) error {
	yamlFn := fn
	jsonFn := extentioner.Swap(fn, "yaml", "json")

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
