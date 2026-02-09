package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Docs struct {
	Examples  map[string]string
	Libraries map[string]string
}

func main() {
	docs := Docs{
		Examples:  make(map[string]string),
		Libraries: make(map[string]string),
	}
	for _, dir := range []string{"examples", "libraries"} {
		m := map[string]string{}
		filepath.Walk(filepath.Join("docs", dir), func(path string, info os.FileInfo, err error) error {
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
				m[filename] = string(content)
			}
			return nil
		})
		switch dir {
		case "examples":
			docs.Examples = m
		case "libraries":
			docs.Libraries = m
		}
	}
	b, _ := json.MarshalIndent(docs, "", "  ")
	f, _ := os.Create("website/src/data/docs.json")
	f.Write(b)
	f.Close()
}
