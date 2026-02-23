package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Docs struct {
	Examples    map[string]ExampleText
	Libraries   map[string]LibraryInfo
	ReportCards map[string]ReportCard
}

func main() {
	docs := Docs{
		Examples:    make(map[string]ExampleText),
		Libraries:   make(map[string]LibraryInfo),
		ReportCards: make(map[string]ReportCard),
	}
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
			docs.Examples[filename] = ParseExample(content)
		}
		return nil
	})

	markdowns := map[string]string{}
	filepath.Walk(filepath.Join("docs", "libraries"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch {
		case filepath.Ext(path) == ".md":
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			filename := filepath.Base(path)
			fmt.Println("Found: ", filename)
			filename = filename[:len(filename)-len(filepath.Ext(filename))]
			markdowns[filename] = string(content)
		case strings.HasSuffix(path, "_samples.json"):
			// todo: create sample data from files and lines
		case filepath.Ext(path) == ".json":
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			filename := filepath.Base(path)
			filename = filename[:len(filename)-len(filepath.Ext(filename))]
			var reportCard ReportCard
			err = json.Unmarshal(content, &reportCard)
			if err != nil {
				fmt.Println(path, err)
				return nil
			}
			if reportCard.Key == "" {
				reportCard.Key = filename
			}
			docs.ReportCards[filename] = reportCard
		}
		return nil
	})
	for name, markdown := range markdowns {
		rc := docs.ReportCards[name]
		docs.Libraries[name] = LibraryInfo{
			Key:          name,
			Name:         rc.Name,
			MarkdownDesc: markdown,
			Website:      rc.Website,
			Repo:         rc.Repo,
			Description:  rc.Description,
			Databases:    rc.Databases,
			License:      rc.License,
			Features:     rc.Features,
			Popularity:   rc.Popularity,
		}
	}
	b, _ := json.MarshalIndent(docs, "", "  ")
	f, _ := os.Create("website/src/data/docs.json")
	f.Write(b)
	f.Close()
}

type LibraryInfo struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	MarkdownDesc string   `json:"markdown_desc"`
	Website      string   `json:"website"`
	Repo         string   `json:"repo"`
	Description  string   `json:"description"`
	Databases    []string `json:"databases"`
	License      string   `json:"license"`
	Features     []string `json:"features"`
	Popularity   int      `json:"popularity"`
}
type ReportCard struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Website     string   `json:"website"`
	Repo        string   `json:"repo"`
	Description string   `json:"description"`
	Databases   []string `json:"databases"`
	License     string   `json:"license"`
	Features    []string `json:"features"`
	Popularity  int      `json:"popularity"`
	Matrix      Matrix   `json:"matrix"`
	Grades      Grades   `json:"grades"`
}
type Matrix struct {
	Databases map[string]GradeResult `json:"databases"`
	Features  map[string]GradeResult `json:"features"`
	Other     map[string]GradeResult `json:"other"`
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
type GradeResult struct {
	Level string `json:"level"`
	Notes string `json:"notes"`
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
