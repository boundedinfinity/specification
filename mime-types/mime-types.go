package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourcesFn := filepath.Join(wd, "iana-media-types-sources.yaml")
	fmt.Printf("Processing %v", sourcesFn)

	var sourcesData ianaMediaTtypesSources
	if err := unmarshal(sourcesFn, &sourcesData); err != nil {
		panic(err)
	}

	for _, source := range sourcesData.MimeTypes {
		if err := downloadFile(wd, source); err != nil {
			panic(err)
		}
	}
}

var (
	cacheDir = "cache"
)

type ianaMediaTtypesSources struct {
	MimeTypes []string `json:"mime-types,omitempty" yaml:"mime-types,omitempty"`
}

func downloadFile(wd, source string) error {
	resp, err := http.Get(source)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fn, err := cacheName(wd, source)
	if err != nil {
		return err
	}

	of, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, fs.FileMode(0755))
	if err != nil {
		return err
	}

	defer of.Close()

	_, err = io.Copy(of, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func cacheName(wd, source string) (string, error) {
	fn := source

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

	absCacheDir := filepath.Join(wd, cacheDir)

	if err := os.MkdirAll(absCacheDir, os.FileMode(0755)); err != nil {
		return fn, err
	}

	fn = filepath.Join(absCacheDir, fn)

	return fn, nil
}

// func marshal(rootDir string, name string, v any) error {
// 	yamlFn := filepath.Join(rootDir, fmt.Sprintf("%v.%v", name, "yaml"))
// 	jsonFn := filepath.Join(rootDir, fmt.Sprintf("%v.%v", name, "json"))

// 	bs, err := yaml.Marshal(v)
// 	if err != nil {
// 		return err
// 	}

// 	err = os.WriteFile(yamlFn, bs, os.FileMode(0755))
// 	if err != nil {
// 		return err
// 	}

// 	bs, err = json.MarshalIndent(v, "", "    ")
// 	if err != nil {
// 		return err
// 	}

// 	err = os.WriteFile(jsonFn, bs, os.FileMode(0755))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func unmarshal(fn string, v any) error {
	content, err := os.ReadFile(fn)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, v)

	return err
}
