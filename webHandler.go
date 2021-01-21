package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

// Struct representing objects sent by the user, containing the BPMs to add to the blockchain
type Message struct {
	BPM int `json:"BPM"`
}

func runServer() error {
	mux := router()

	server := &http.Server{
		// Addr        : fmt.Sprintf(":%s", os.Getenv("PORT")),
		Addr        : fmt.Sprintf("%s:%s",os.Getenv("HOST") , os.Getenv("PORT")),
		Handler     : mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout : 15 * time.Second,
	}

	fmt.Printf("Server listening on port %s\n", os.Getenv("PORT"))

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func router() http.Handler {
	router := mux.NewRouter()

	// Initialize the router
	router.HandleFunc("/", getBlockchainHandler).Methods("GET")
	router.HandleFunc("/", postBlockchainHandler).Methods("POST")

	return router
}

func getBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	/*
	bytes, err := json.Marshal(Blockchain)
	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, "")
		return
	}
	 */
	respondWithJson(w, r, http.StatusOK, Blockchain)
}

func postBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	var input Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	newBlock, err := newBlock(Blockchain[len(Blockchain)-1], input.BPM)
	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, input)
		return
	}

	if blockIsValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
	}

	respondWithJson(w, r, http.StatusCreated, newBlock)
}

func respondWithJson(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")

	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}