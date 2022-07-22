package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type FrontModule struct {
	Name string `json:"name"`
}

func modulesHandler(w http.ResponseWriter, r *http.Request) {
	var modulesArray []FrontModule

	baseDir := os.Getenv("baseDir")

	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			subFiles, err := ioutil.ReadDir(baseDir + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, subFile := range subFiles {
				if f.IsDir() {
					var item FrontModule
					item.Name = f.Name() + "/" + subFile.Name()
					modulesArray = append(modulesArray, item)
				}
			}
		}
	}

	result, _ := json.Marshal(modulesArray)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func main() {
	baseDir := os.Getenv("baseDir")

	r := mux.NewRouter()
	r.HandleFunc("/list", modulesHandler)
	pr := "/modules/"
	r.PathPrefix(pr).Handler(http.StripPrefix(pr, http.FileServer(http.Dir(baseDir))))
	http.Handle("/", r)

	host := os.Getenv("server.Host")
	port := os.Getenv("server.Port")
	addr := host + ":" + port

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
