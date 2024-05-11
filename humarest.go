package m5cli

import (
	"fmt"
	// L "github.com/fbaube/mlog"
	// "github.com/gorilla/mux"
	"net/http"
	"strconv"
	CTX "context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

var sRestPortNr string

type HelloReq struct {
     Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type HelloRsp struct {
	  Body struct {
	Message string `json:"message" example:"Hello, world!" doc:"Hello msg"`
	}
}

type StaticReq struct {
     Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type HtmlRsp struct {
	  Body struct {
	Message string `json:"message" example:"Hello, world!" doc:"Hello msg"`
	}
}

type HumaHandler[I,O any] func(CTX.Context, *I) (*O, error)

func DoHello/*[HelloReq, HelloRsp]*/(ctx CTX.Context, pReq *HelloReq) (pRsp *HelloRsp, e error) {
     println("GET!-DoHello")
     pRsp = new(HelloRsp)
     pRsp.Body.Message = fmt.Sprintf("Hello, %s!", pReq.Name)
     return pRsp, nil
}

type HumaOpSpec[I,O any] struct {
     HttpVerb string
     UrlPatrn string
     InStruct *I
     OutStruct *O
     // HH HumaHandler
}

/*
var OpSpecs = []HumaOpSpec {
    /*HumaOpSpec[HelloReq, HelloRsp]* / {
    	"GET", "/hello/{name}", *HelloReq, *HelloRsp, DoHello },
}
*/

func RunRest(portNr int) error {
	if portNr == 0 { // env.RestPort
		return nil
	}
	sRestPortNr = strconv.Itoa(portNr)
	println("==> Running Huma-REST server on port:", sRestPortNr)

	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Derf API", "0.0.1"))

	// VERB + URL-PATTERN + IN-STRUCT + OUT+STRUCT + FUNC
	huma.Get(api, "/hello/{name}",
		func(ctx CTX.Context, I *HelloReq) (*HelloRsp, error) {
		println("GET!")
		pRsp := new(HelloRsp)
		pRsp.Body.Message = fmt.Sprintf("Hello, %s!", I.Name)
		return pRsp, nil
	})
	// OR just load this "ABOUT" into the mux as a normal HTTP Handler ???
	huma.Get(api, "/about",
		func(ctx CTX.Context, I *StaticReq) (*HtmlRsp, error) {
		println("GET-STATIC")
		pRsp := new(HtmlRsp)
		pRsp.Body.Message = 
			"<!DOCTYPE html>\n<html>\n<body>\nABOUT!\n" +
			"</body></html>"
		// pRsp
		return pRsp, nil
	})

	// fmt.Printf("API: %+v \n", api)

	// Start the server!
	// http.ListenAndServe("127.0.0.1:8888", mux)
	http.ListenAndServe("localhost:8888", mux)
	return nil
}

/*

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

* /

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

