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

/*
https://huma.rocks/features/request-inputs/
Requests can have parameters and/or a body as input to the handler function.
Inputs use standard Go structs with special fields and/or tags.
Here are the available tags:

Tag	Description	Example
path	Name of the path parameter	path:"thing-id"
query	Name of the query string parm	query:"q"
header	Name of the header parameter	header:"Authorization"
cookie	Name of the cookie parameter	cookie:"session"
required Mark query/hdr parm as req'd	required:"true"

Request Body
The special struct field Body will be treated as the input request
body and can refer to any other type or you can embed a struct or
slice inline. If the body is a pointer, then it is optional. All
doc & validation tags are allowed on the body plus these tags:

Tag	Description	Example
contentType Override the content type	contentType:"application/my-type+json"
required    Mark the body as required	required:"true"

Response Headers
Headers are set by fields on the response struct. Available tags:

Tag     	Description	Example
header     Name of response header header:"Authorization"
timeFormat Format of time.Time     timeFormat:"Mon, 02 Jan 2006 15:04:05 GMT"

Resoonse Body
The special struct field Body will be treated as the response body and
can refer to any other type or you can embed a struct or slice inline.
A default Content-Type header will be set if none is present, selected
via client-driven content negotiation with the server based on the
registered serialization types.

*/

type HelloReq struct {
     Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
     	/* // ADDED STUFF
	ID      string `path:"id"`
	Detail  bool   `query:"detail" doc:"Show full details"`
	Auth    string `header:"Authorization"`
	Body    MyBody
	RawBody []byte */
}

type HelloRsp struct {
	  Body struct {
	Message string `json:"message" example:"Hello, world!" doc:"Hello msg"`
	}
}

type StaticReq struct {
     // Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type HtmlRsp struct {
	ContentType  string    `header:"Content-Type"`
	// LastModified time.Time `header:"Last-Modified"`
	// MyHeader     int       `header:"My-Header"`
	Body []byte 
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
		func(ctx CTX.Context, Z *struct{}) (*HtmlRsp, error) {
		println("GET-STATIC")
		pRsp := new(HtmlRsp)
		pRsp.ContentType = "text/html" 
		pRsp.Body = []byte(
			"<!DOCTYPE html>\n<html>\n<body>\nABOUT!\n" +
			"</body></html>")
		return pRsp, nil
	})

	// fmt.Printf("API: %+v \n", api)

	// Start the server!
	http.ListenAndServe("127.0.0.1:8888", mux)
	// http.ListenAndServe("localhost:8888", mux)
	println("==> Running Huma-REST server on port:", sRestPortNr)
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

