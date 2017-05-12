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
	"github.com/theopticians/optician-api/core/store"
	"github.com/theopticians/optician-api/core/structs"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/cases", addCaseHandler).Methods("POST")
	r.HandleFunc("/batches", getBatchsHandler).Methods("GET")
	r.HandleFunc("/batches/{id}", getResultsByBatchHandler).Methods("GET")
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

func getBatchsHandler(rw http.ResponseWriter, req *http.Request) {
	batchs, err := core.Batchs()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	trJSON, err := json.Marshal(batchs)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(trJSON)
}

func getResultsByBatchHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	tests, err := core.ResultsByBatchs(id)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	trJSON, err := json.Marshal(tests)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(trJSON)
}

func getResultsHandler(rw http.ResponseWriter, req *http.Request) {
	tests, err := core.Results()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	trJSON, err := json.Marshal(tests)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(trJSON)
}

func addCaseHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var c ApiCase
	err := decoder.Decode(&c)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	defer req.Body.Close()

	results, err := core.AddCase(structs.Case(c))

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Header().Set("Location", "/results/"+results.ID)
	rw.WriteHeader(http.StatusCreated)
}

func getResultHandler(rw http.ResponseWriter, req *http.Request) {
	var results structs.Result
	var err error
	vars := mux.Vars(req)
	id := vars["id"]

	if id != "" {
		results, err = core.GetTest(id)
		if err != nil {
			if err == store.NotFoundError {
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

	apiResult := ApiResult(results)
	trJSON, err := json.Marshal(&apiResult)

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
		if err == store.NotFoundError {
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
	m := &structs.Mask{}
	err := decoder.Decode(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	// TODO return new results
	_, err = core.MaskTest(id, []image.Rectangle(*m))

	if err != nil {
		if err == store.NotFoundError {
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
