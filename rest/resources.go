package rest

import (
	"net/http"
	"fmt"
	"html"
	"resys/service"
	"encoding/json"
)

type Resource struct {
	Uri string // read-only
	Name string
	Tags []string
}

const ResourcesPath = "/resources/"


// TODO Inject?
var resourceService service.ResourceService

func NewResourcesHandler(rs service.ResourceService) func(w http.ResponseWriter, r *http.Request) {
	resourceService = rs
	return ResourceHandler
}

func ResourceHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "POST %q", html.EscapeString(r.URL.Path))
		break
	case http.MethodGet:
		id := r.URL.Path[len(ResourcesPath):]

		w.Header().Set("Content-Type", "application/json")

		if id == "" {
			getAllResources(w)
		} else {
			getResource(id, w)
		}

		break
	}
}

func getResource(id string, w http.ResponseWriter) {
	sr := resourceService.Get(id)
	res := convertResourceFromService(sr)

	resJson, err := json.Marshal(res)
	if err != nil {
		panic(err) // TODO: error handler
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}


func getAllResources(w http.ResponseWriter) {

	// begin
	b := func(err error) {
		// TODO: check error
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"Resources\":["))
	}

	// next
	first := true
	n := func(sr *service.Resource, err error) {
		res := convertResourceFromService(sr)

		resJson, err := json.Marshal(res)
		if err != nil {
			panic(err) // TODO: error handler
		}

		fmt.Println(res)

		if !first {
			w.Write([]byte(","))
		}
		w.Write(resJson)
		first = false
	}

	// end
	e := func(err error) {
		w.Write([]byte("]}"))
	}

	resourceService.GetAll(b, n, e)
}


// Conversion Service to REST object
func convertResourceFromService(sr *service.Resource) *Resource {
	return &Resource{
		Uri: MakeIdPath(ResourcesPath, sr.Id),
		Name: sr.Name,
		Tags: sr.Tags,
	}
}