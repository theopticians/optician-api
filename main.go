package main

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/theopticians/optician-api/core"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/cases", addCaseHandler).Methods("POST")
	r.HandleFunc("/results", getResultsHandler).Methods("GET")
	r.HandleFunc("/results/{id}", getResultHandler).Methods("GET")
	r.HandleFunc("/results/{id}/accept", acceptHandler).Methods("POST")
	r.HandleFunc("/results/{id}/mask", maskHandler).Methods("POST")
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

func getResultsHandler(rw http.ResponseWriter, req *http.Request) {
	testList := core.TestList()

	trJSON, err := json.Marshal(testList)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(trJSON)
}

func addCaseHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var c Case
	err := decoder.Decode(&c)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	defer req.Body.Close()

	results, err := core.AddCase(core.Case(c))

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	trJSON, err := json.Marshal(results)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(trJSON)
}

func getResultHandler(rw http.ResponseWriter, req *http.Request) {
	var results core.Result
	var err error
	vars := mux.Vars(req)
	id := vars["id"]

	if id != "" {
		results, err = core.GetTest(id)
		if err != nil {
			if err == core.NotFoundError {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
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
		if err == core.NotFoundError {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func maskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	decoder := json.NewDecoder(r.Body)
	var m []image.Rectangle
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	// TODO return new results
	_, err = core.MaskTest(id, m)

	if err != nil {
		if err == core.NotFoundError {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
