package cli

import (
	"fmt"
	L "github.com/fbaube/mlog"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var sRestPortNr string

// Does not launch a goroutine.
// Call with env.RestPort .
// Uses Gorilla mux, mainly cos now it's in archive mode.
// Instructions for usage are found
// [here](https://github.com/gorilla/mux#examples)
// or in alternate format [github.com/gorilla/mux]
// or in alternate format [mux]
// or in alternate format "[muxdox]"
//
// [muxdox]: https://github.com/gorilla/mux#examples
// .
func RunRest(portNr int) error {
	if portNr == 0 { // env.RestPort
		return nil
	}
	sRestPortNr = strconv.Itoa(portNr)
	println("==> Running REST server on port:", sRestPortNr)
	r := mux.NewRouter()

	// mux docs say:
	// Routes are tested in the order they were added to
	// the router. If two routes match, the first one wins.

	// ADMIN
	r.HandleFunc("/stc", hdlStcRoot)
	// TOPICS, MAPS, DATABASE, STATIC CONTENT
	rtrTpc := r.PathPrefix("/tpc").Subrouter()
	rtrMap := r.PathPrefix("/map").Subrouter()
	rtrDb := r.PathPrefix("/db").Subrouter()
	// HOME (incl. "About", etc.)
	r.HandleFunc("/", HomeHandler)

	// TOPICS
	rtrTpc.HandleFunc("/{id}/meta", hdlTopicMeta)
	rtrTpc.HandleFunc("/{id}/links", hdlTopicLinks)
	rtrTpc.HandleFunc("/{id}", hdlTopic)
	rtrTpc.HandleFunc("/", hdlTopicRoot)

	// MAPS
	rtrMap.HandleFunc("/{id}/meta", hdlMapMeta)
	rtrMap.HandleFunc("/{id}/links", hdlMapLinks)
	rtrMap.HandleFunc("/{id}", hdlMap)
	rtrMap.HandleFunc("/", hdlMapRoot)

	// DB (schemas? stats?)
	rtrDb.HandleFunc("/fld/{name}", hdlDbField)
	rtrDb.HandleFunc("/{name}", hdlDbTable)
	rtrDb.HandleFunc("/", hdlDbRoot)

	// go func() {
	if err := http.ListenAndServe(":"+sRestPortNr, r); err != nil {
		L.L.Error("REST server failed: " + err.Error())
	}
	return nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("home vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)
		fmt.Fprintf(w, s)
	*/
	L.L.Info(s)
}

func TopicRootHdlr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("topic vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)'

		fmt.Fprintf(w, s)
	*/
	L.L.Info(s)
}
func MapRootHdlr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("map vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)
		fmt.Fprintf(w, s)
	*/
	L.L.Info(s)
}

func StcRootHdlr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("static vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)
		fmt.Fprintf(w, s)
	*/
	L.L.Info(s)

	// This will serve files under http://localhost:8000/static/<filename>
	// r.PathPrefix("/s/").Handler(http.StripPrefix("/s/", http.FileServer(http.Dir(dir))))
}
