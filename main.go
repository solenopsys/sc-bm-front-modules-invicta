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
	var modulesArray = make([]FrontModule, 0)

	log.Println("REQUEST ", r.GetBody)

	baseDir := os.Getenv("baseDir")

	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Println("DIR NOT FOUND ", baseDir)
	} else {
		for _, f := range files {
			if f.IsDir() {
				subDir := baseDir + "/" + f.Name()
				subFiles, err := ioutil.ReadDir(subDir)
				if err != nil {
					log.Println("SUB DIR NOT FOUND ", subDir)
				} else {
					for _, subFile := range subFiles {
						if f.IsDir() {
							var item FrontModule
							item.Name = f.Name() + "/" + subFile.Name()
							modulesArray = append(modulesArray, item)
						}
					}
				}

			}
		}
	}

	result, _ := json.Marshal(modulesArray)

	w.Header().Set("Content-Type", "application/json")
	setCorsHeaders(w)
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do your cors stuff
		// return if you do not want the FileServer handle a specific request
		setCorsHeaders(w)
		fs.ServeHTTP(w, r)
	}
}

func main() {
	baseDir := os.Getenv("baseDir")

	r := mux.NewRouter()
	r.HandleFunc("/fm/list", modulesHandler)
	pr := "/fm/modules/"
	r.PathPrefix(pr).Handler(http.StripPrefix(pr, cors(http.FileServer(http.Dir(baseDir)))))
	http.Handle("/", r)

	host := os.Getenv("server.Host")
	port := os.Getenv("server.Port")

	log.Println("START SERVER host,port ", host, " ", port)
	log.Println("DIR ", baseDir)
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
