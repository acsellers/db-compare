package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	GenerateLibraries()
	GenerateExamples()
}

var Libraries = []string{
	"bob",
	"bun",
	"ent",
	"gorm",
	"jet",
	"sqlboiler",
	"sqlc",
	"sqlx",
	"upper db",
	"xorm",
}

type Features struct {
	Databases map[string]GradeResult `json:"databases"`
	Features  map[string]GradeResult `json:"features"`
	Other     map[string]GradeResult `json:"other"`
}

type GradeResult struct {
	Level string `json:"level"`
	Notes string `json:"notes"`
}

type Grades struct {
	GetSale          GradeResult `json:"get_sale"`
	CreateSale       GradeResult `json:"create_sale"`
	CustomerUpdate   GradeResult `json:"customer_update"`
	BasicGrouping    GradeResult `json:"basic_grouping"`
	AdvancedGrouping GradeResult `json:"advanced_grouping"`
	WithQueries      GradeResult `json:"with_queries"`
	SaleSearch       GradeResult `json:"sale_search"`
	BulkCustomers    GradeResult `json:"bulk_customers"`
	JSON             GradeResult `json:"json"`
}

type LibraryInfo struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	ShortDesc    string   `json:"short_description"`
	MarkdownDesc string   `json:"markdown_desc"`
	Website      string   `json:"website"`
	Repo         string   `json:"repo"`
	Databases    []string `json:"databases"`
	License      string   `json:"license"`
	Features     []string `json:"features"`
	Popularity   int      `json:"popularity"`
}

type Samples struct {
	Store StoreSamples `json:"store"`
}
type StoreSamples struct {
	MySQL    StoreDBSamples `json:"mysql"`
	Postgres StoreDBSamples `json:"postgres"`
	SQLite   StoreDBSamples `json:"sqlite"`
}
type StoreDBSamples struct {
	GetSale          SampleRefs `json:"get_sale"`
	CreateSale       SampleRefs `json:"create_sale"`
	CustomerUpdate   SubSamples `json:"customer_update"`
	BasicGrouping    SubSamples `json:"basic_grouping"`
	AdvancedGrouping SubSamples `json:"advanced_grouping"`
	WithQueries      SubSamples `json:"with_queries"`
	SaleSearch       SampleRefs `json:"sale_search"`
	BulkCustomers    SampleRefs `json:"bulk_customers"`
	JSON             SampleRefs `json:"json"`
}
type SampleRefs struct {
	File  string `json:"file"`
	Query string `json:"query"`
}
type SubSamples struct {
	SubExamples map[string]SampleRefs `json:"sub_examples"`
}

type FinalLibraryInfo struct {
	Info     LibraryInfo `json:"info"`
	Grades   Grades      `json:"grades"`
	Features Features    `json:"features"`
}

func GenerateLibraries() {
	libs := map[string]FinalLibraryInfo{}
	for _, lib := range Libraries {
		fli := FinalLibraryInfo{}
		infoFile, err := os.Open(fmt.Sprintf("docs/libraries/%s/info.json", lib))
		if err != nil {
			fmt.Println("Missing info.json for ", lib)
			log.Fatal(err)
		}
		defer infoFile.Close()
		json.NewDecoder(infoFile).Decode(&fli.Info)

		descFile, err := os.ReadFile(fmt.Sprintf("docs/libraries/%s/description.md", lib))
		if err != nil {
			fmt.Println("Missing description.md for ", lib)
			log.Fatal(err)
		}
		fli.Info.MarkdownDesc = string(descFile)

		gradesFile, err := os.Open(fmt.Sprintf("docs/libraries/%s/grades.json", lib))
		if err != nil {
			fmt.Println("Missing grades.json for ", lib)
			log.Fatal(err)
		}
		defer gradesFile.Close()
		json.NewDecoder(gradesFile).Decode(&fli.Grades)

		featuresFile, err := os.Open(fmt.Sprintf("docs/libraries/%s/features.json", lib))
		if err != nil {
			fmt.Println("Missing features.json for ", lib)
			log.Fatal(err)
		}
		defer featuresFile.Close()
		json.NewDecoder(featuresFile).Decode(&fli.Features)

		libs[lib] = fli

		// TODO: Create a json file with all the samples by opening the samples.json file
		// and using the references to open the actual files and save them to a new json file
		// in the public folder under website/samples/<lib>/<sample>.json
	}

	docsFile, err := os.Create("website/src/data/docs.json")
	if err != nil {
		fmt.Println("Missing docs.json for ", libs)
		log.Fatal(err)
	}
	defer docsFile.Close()
	jsonBytes, err := json.MarshalIndent(libs, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling docs.json for ", libs)
		log.Fatal(err)
	}
	docsFile.Write(jsonBytes)
}

type ExampleText struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	SubExamples []SubExampleText `json:"sub_examples"`
}
type SubExampleText struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

func GenerateExamples() {
	examples := map[string]ExampleText{}
	filepath.Walk(filepath.Join("docs", "examples"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".md" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			filename := filepath.Base(path)
			filename = filename[:len(filename)-len(filepath.Ext(filename))]
			examples[filename] = ParseExample(content)
		}
		return nil
	})

	examplesFile, err := os.Create("website/src/data/examples.json")
	if err != nil {
		fmt.Println("Missing examples.json for ", examples)
		log.Fatal(err)
	}
	defer examplesFile.Close()
	jsonBytes, err := json.MarshalIndent(examples, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling examples.json for ", examples)
		log.Fatal(err)
	}
	examplesFile.Write(jsonBytes)
}

func ParseExample(content []byte) ExampleText {
	et := ExampleText{}
	scanner := bufio.NewScanner(bytes.NewBuffer(content))
	desc := ""
	inSubExample := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			et.Title = strings.TrimPrefix(line, "# ")
		} else if strings.HasPrefix(line, "## ") {
			if !inSubExample {
				et.Description = strings.TrimSpace(desc)
				inSubExample = true
			} else {
				et.SubExamples[len(et.SubExamples)-1].Description = strings.TrimSpace(desc)
			}
			desc = ""
			et.SubExamples = append(et.SubExamples, SubExampleText{
				Title: strings.TrimPrefix(line, "## "),
			})
		} else {
			desc += line + "\n"
		}
	}
	if !inSubExample {
		et.Description = strings.TrimSpace(desc)
	} else {
		et.SubExamples[len(et.SubExamples)-1].Description = strings.TrimSpace(desc)
	}
	return et
}
