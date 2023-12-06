package m5cli

import (
	// "fmt"
	L "github.com/fbaube/mlog"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var sWebPortNr string

// Notes:
//  - Does not launch a goroutine.
//  - Call with env.WebPort .
//  - Used Gorilla mux, but now use pilot version of new Go ServeMux.
//
// Dedicated servers:
//  - /events :: SSE
//  - /ws :: websocket upgrader
//  - /rest :: REST server (on a different port nr)
//  - /htmx :: htmx responder (??; https://github.com/angelofallars/htmx-go)
// Main server (with h2c):
//  - /all :: list of all endpoints 
//  - / :: something bland
//  - /app :: the wasm SPA
//  - /static :: static content 
//  - /about
//  - /contact
//  - /health (or?)
// - /db
// .
func RunWeb(portNr int) error {
	if portNr == 0 { // env.WebPort
		return nil
	}
	sWebPortNr = strconv.Itoa(portNr)
	println("==> Running WEB server on port:", sWebPortNr)
	r := mux.NewRouter()

	// mux docs say:
	// Routes are tested in the order they were added to
	// the router. If two routes match, the first one wins.

	/* OOPS

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

	*/

	// go func() {
	if err := http.ListenAndServe(":"+sWebPortNr, r); err != nil {
		L.L.Error("WEB server failed: " + err.Error())
	}
	return nil
}

/*

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("home vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)
		fmt.Fprintf(w, s)
	* /
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
	* /
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
	* /
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
	* /
	L.L.Info(s)

	// This will serve files under http://localhost:8000/static/<filename>
	// r.PathPrefix("/s/").Handler(http.StripPrefix("/s/", http.FileServer(http.Dir(dir))))
}

*/
