package main

import (
	"bytes"
	"encoding/json"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/gerardabello/optician/core"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/test", testHandler)
	r.HandleFunc("/results/{id}", resultHandler)
	r.HandleFunc("/accept/{id}", acceptHandler)
	r.HandleFunc("/image/{id}", imageHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func testHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t Test
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	defer req.Body.Close()

	results, err := core.TestImage(t.Image, t.ProjectID, t.Branch, t.Target, t.Browser)

	if err != nil {
		panic(err)
	}

	trJSON, err := json.Marshal(results)

	if err != nil {
		panic(err)
	}

	rw.Write(trJSON)
}

func resultHandler(rw http.ResponseWriter, req *http.Request) {
	var results core.Results
	var err error

	idParam := req.URL.Query().Get("id")
	if idParam != "" {
		results, err = core.GetResults(idParam)
		if err != nil {
			panic(err)
		}
	}

	trJSON, err := json.Marshal(results)

	if err != nil {
		panic(err)
	}

	rw.Write(trJSON)
}

func acceptHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := core.AcceptTest(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	img := core.GetImage(id)

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}

}
