package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	http.HandleFunc("/filelist", enableCORS(FileList))
	http.HandleFunc("/file", enableCORS(GetFile))
	http.HandleFunc("/save", enableCORS(SaveSamples))
	http.HandleFunc("/samples", enableCORS(GetSamples))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type FileListing struct {
	Library     string              `json:"library"`
	SamplesByDB map[string][]string `json:"samples_by_db"`
}

func FileList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	files := []string{
		"bob",
		"bob2",
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
	fileLists := make([]FileListing, 0)
	for _, file := range files {
		fileList := FileListing{
			Library:     file,
			SamplesByDB: make(map[string][]string),
		}
		filepath.Walk("../store/mysql/"+file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				name := path
				name = strings.TrimPrefix(name, "../store/mysql/"+file+"/")
				fileList.SamplesByDB["mysql"] = append(fileList.SamplesByDB["mysql"], name)
			}
			return nil
		})
		filepath.Walk("../store/postgres/"+file+"/*.*", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				name := path
				name = strings.TrimPrefix(name, "../store/postgres/"+file+"/")
				fileList.SamplesByDB["postgres"] = append(fileList.SamplesByDB["postgres"], name)
			}
			return nil
		})
		filepath.Walk("../store/sqlite/"+file+"/*.*", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				name := path
				name = strings.TrimPrefix(name, "../store/sqlite/"+file+"/")
				fileList.SamplesByDB["sqlite"] = append(fileList.SamplesByDB["sqlite"], name)
			}
			return nil
		})
		fileLists = append(fileLists, fileList)
	}
	json.NewEncoder(w).Encode(fileLists)
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	library := r.URL.Query().Get("library")
	db := r.URL.Query().Get("db")
	file := r.URL.Query().Get("file")

	f, err := os.Open("../store/" + db + "/" + library + "/" + file)
	if err != nil {
		w.Write([]byte("{}"))
		return
	}
	defer f.Close()

	content := &bytes.Buffer{}
	io.Copy(content, f)
	lines := strings.Split(content.String(), "\n")
	json.NewEncoder(w).Encode(lines)
}

func GetSamples(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	lib := r.URL.Query().Get("library")
	f, err := os.Open("../docs/libraries/" + lib + "/samples.json")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("{}"))
		return
	}
	defer f.Close()
	data := make(map[string]interface{})
	json.NewDecoder(f).Decode(&data)
	json.NewEncoder(w).Encode(data)
}

func SaveSamples(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.Write([]byte("{}"))
		return
	}
	lib := r.URL.Query().Get("library")
	f, err := os.Create("../docs/libraries/" + lib + "/samples.json")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("{}"))
		return
	}
	defer f.Close()
	b, _ := json.MarshalIndent(data, "", "  ")
	f.Write(b)
	fmt.Println(string(b))
}
