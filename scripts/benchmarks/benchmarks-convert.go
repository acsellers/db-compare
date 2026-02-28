package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type BenchmarkRun struct {
	Key     string          `json:"key"`
	RunDate string          `json:"run_date"`
	Items   []BenchmarkItem `json:"items"`
}
type BenchmarkItem struct {
	Name     string  `json:"name"`
	Time     int     `json:"time"`
	VsStdlib float64 `json:"vs_stdlib"`
	Runs     []int   `json:"runs"`
	Rating   string  `json:"rating"`
	Notes    string  `json:"notes"`
}
type Benchmarks struct {
	RunDate string          `json:"run_date"`
	Items   []BenchmarkItem `json:"items"`
}

var benches = map[string]BenchmarkRun{}

func main() {

	LoadStdlib()
	stdlib := map[string]float64{}
	for _, item := range benches["stdlib"].Items {
		stdlib[item.Name] = float64(item.Time)
	}
	for _, lib := range Libraries {
		LoadLibrary(lib)
		for _, item := range benches[lib].Items {
			item.VsStdlib = float64(item.Time) / stdlib[item.Name]
		}
	}
	outputFile, err := os.Create("website/src/data/benchmarks.json")
	if err != nil {
		fmt.Println("Missing benchmarks.json for ", benches)
		log.Fatal(err)
	}
	defer outputFile.Close()
	jsonBytes, err := json.MarshalIndent(benches, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling benchmarks.json for ", benches)
		log.Fatal(err)
	}
	outputFile.Write(jsonBytes)
}
func LoadStdlib() {
	b := loadData("stdlib")
	benches["stdlib"] = b
}
func LoadLibrary(lib string) {
	b := loadData(lib)
	benches[lib] = b
}
func loadData(lib string) BenchmarkRun {
	file, err := os.Open(fmt.Sprintf("docs/libraries/%s/benchmarks.json", lib))
	if err != nil {
		fmt.Println("Missing benchmarks.json for ", lib)
		log.Fatal(err)
	}
	defer file.Close()
	bench := BenchmarkRun{}
	json.NewDecoder(file).Decode(&benches[lib])

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
