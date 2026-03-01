package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var samples = []string{
	"bob",
	/*"bob2",
	"bun",
	"ent",
	"gorm",
	"jet",
	"sqlboiler",
	"sqlc",
	"sqlx",
	"upper db",
	"xorm",*/
}

type Samples struct {
	Store StoreSamples `json:"store"`
}
type StoreSamples struct {
	MySQL    *StoreDBSamples `json:"mysql"`
	Postgres *StoreDBSamples `json:"postgres"`
	SQLite   *StoreDBSamples `json:"sqlite"`
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
	Go  SampleFile `json:"go"`
	Sql SampleFile `json:"sql"`
}
type SampleFile struct {
	File       string   `json:"file"`
	Lines      [][2]int `json:"lines"`
	Highlights [][2]int `json:"highlights"`
}
type SubSamples map[string]SampleFile

func main() {
	for _, lib := range samples {
		sampleData := Samples{}
		f, err := os.Open("docs/libraries/" + lib + "/samples.json")
		if err != nil {
			log.Fatal("Open Samples.json: ", err)
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(&sampleData)
		if err != nil {
			log.Fatal("Parse samples.json: ", err)
		}
		processSamples(lib, sampleData)
	}
}

func processSamples(lib string, samples Samples) {
	if samples.Store.MySQL != nil {
		processDB(lib, "mysql", *samples.Store.MySQL)
	}
	if samples.Store.Postgres != nil {
		processDB(lib, "postgres", *samples.Store.Postgres)
	}
	if samples.Store.SQLite != nil {
		processDB(lib, "sqlite", *samples.Store.SQLite)
	}
}

func processDB(lib string, db string, samples StoreDBSamples) {
	filePrefix := fmt.Sprintf("website/public/samples/%s-%s", lib, db)
	saveSample(filePrefix, "create-sale", processSample(lib, db, samples.CreateSale))
}

func saveSample(filePrefix string, name string, data []byte) {
	f, err := os.Create(filePrefix + "_" + name + ".json")
	if err != nil {
		log.Fatal("Create file: ", err)
	}
	defer f.Close()
	f.Write(data)
}

type SampleOutput struct {
	GoFile        string   `json:"go_file"`
	SqlFile       string   `json:"sql_file"`
	GoSrc         string   `json:"go_src"`
	SqlSrc        string   `json:"sql_src"`
	GoLines       []string `json:"go_lines"`
	SqlLines      []string `json:"sql_lines"`
	GoHighlights  []string `json:"go_highlights"`
	SqlHighlights []string `json:"sql_highlights"`
}

func processSample(lib string, db string, sample SampleRefs) []byte {
	out := SampleOutput{
		GoFile:  sample.Go.File,
		SqlFile: sample.Sql.File,
	}
	srcPrefix := fmt.Sprintf("store/%s/%s", db, lib)
	if sample.Go.File != "" {
		out.GoSrc, out.GoLines, out.GoHighlights = getLines(srcPrefix+"/"+sample.Go.File, sample.Go.Lines, sample.Go.Highlights)
	}
	if sample.Sql.File != "" {
		out.SqlSrc, out.SqlLines, out.SqlHighlights = getLines(srcPrefix+"/"+sample.Sql.File, sample.Sql.Lines, sample.Sql.Highlights)
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent: ", err)
	}
	return b
}
func getLines(file string, lines [][2]int, highlights [][2]int) (string, []string, []string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Open file: ", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	srcLines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		srcLines = append(srcLines, line)
	}

	fileLines := make([]string, 0)
	fileHighlights := make([]string, 0)

	if len(lines) == 0 {
		lines = append(lines, [2]int{1, len(srcLines)})
	}
	for _, rng := range lines {
		start := rng[0]
		end := rng[1]
		curr := []string{}
		for i := start; i <= end; i++ {
			curr = append(curr, srcLines[i-1])
		}
		fileLines = append(fileLines, strings.Join(curr, "\n"))
	}
	if len(highlights) == 0 {
		fileHighlights = fileLines
	}
	for _, rng := range highlights {
		start := rng[0]
		end := rng[1]
		curr := []string{}
		for i := start; i <= end; i++ {
			curr = append(curr, srcLines[i-1])
		}
		fileHighlights = append(fileHighlights, strings.Join(curr, "\n"))
	}
	return strings.Join(srcLines, "\n"), fileLines, fileHighlights
}
