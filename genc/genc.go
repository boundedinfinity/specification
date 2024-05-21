package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

// https://gobyexample.com/xml

func main() {
	fn := os.Args[1]
	fmt.Printf("Processing %v\n", fn)
}

type GencStandardBaseline struct {
	XMLName          xml.Name              `xml:"genc:GENCStandardBaseline"`
	Authority        string                `xml:"genc:authority"`
	Baseline         string                `xml:"genc:baseline"`
	PromulgationDate string                `xml:"genc:promulgationDate"`
	Country          []GencISOCountryEntry `xml:"genc:ISOCountryEntry"`
}

type GencISOCountryEntry struct {
	XMLName xml.Name `xml:"genc:ISOCountryEntry"`
}

type GencEncoding struct {
	XMLName                  xml.Name             `xml:"genc:encoding"`
	Char3Code                string               `xml:"genc:char3Code"`
	Char3CodeURISet          GencCharCodeURISet   `xml:"genc:char3CodeURISet"`
	Char2Code                string               `xml:"genc:char2Code"`
	Char2CodeURISet          GencCharCodeURISet   `xml:"genc:char2CodeURISet"`
	NumericCode              string               `xml:"genc:numericCode"`
	NumericCodeURISet        GencCharCodeURISet   `xml:"genc:numericCodeURISet"`
	Char3CodeElementStatus   string               `xml:"genc:char3CodeElementStatus"`
	Char2CodeElementStatus   string               `xml:"genc:char2CodeElementStatus"`
	NumericCodeElementStatus string               `xml:"genc:numericCodeElementStatus"`
	Name                     string               `xml:"genc:name"`
	ShortName                string               `xml:"genc:shortName"`
	FullName                 string               `xml:"genc:fullName"`
	UnLegalStatus            string               `xml:"genc:unLegalStatus"`
	EntryDate                string               `xml:"genc:entryDate"`
	LocalShortName           []GencLocalShortName `xml:"genc:localShortName"`
	BgnShortNameVariance     string               `xml:"genc:bgnShortNameVariance"`
	BgnFullNameVariance      string               `xml:"genc:bgnFullNameVariance"`
	UsRecognition            string               `xml:"genc:usRecognition"`
}

type GencCharCodeURISet struct {
	CodespaceURL           string `xml:"genc:codespaceURL"`
	CodespaceURN           string `xml:"genc:codespaceURN"`
	CodespaceURNBased      string `xml:"genc:codespaceURNBased"`
	CodespaceURNBasedShort string `xml:"genc:codespaceURNBasedShort"`
}

type GencLocalShortName struct {
	NameLanguage2Char string `xml:"genc:nameLanguage2Char"`
	NameLanguage3Char string `xml:"genc:nameLanguage3Char"`
	Name              string `xml:"genc:name"`
	Iso6393Char3Code  string `xml:"genc:iso6393Char3Code"`
}
