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

	r.HandleFunc("/tests", getTestsHandler).Methods("GET")
	r.HandleFunc("/tests", runTestHandler).Methods("POST")
	r.HandleFunc("/tests/{id}", resultHandler).Methods("GET")
	r.HandleFunc("/tests/{id}/accept", acceptHandler).Methods("POST")
	r.HandleFunc("/image/{id}", imageHandler).Methods("GET")

	http.Handle("/", middleware(r))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(rw, req)
	})
}

func getTestsHandler(rw http.ResponseWriter, req *http.Request) {
	testList := core.TestList()

	trJSON, err := json.Marshal(testList)

	if err != nil {
		panic(err)
	}

	rw.Write(trJSON)
}

func runTestHandler(rw http.ResponseWriter, req *http.Request) {
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
	vars := mux.Vars(req)
	id := vars["id"]

	if id != "" {
		results, err = core.GetResults(id)
		if err != nil {
			if err == core.NotFoundError {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to encode image."))
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
