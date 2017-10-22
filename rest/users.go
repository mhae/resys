package rest

import (
	"net/http"
	"fmt"
	"html"
	"resys/service"
	"encoding/json"
)

type User struct {
	Uri string
	Name string
}

const UsersPath = "/users/"


// TODO Inject?
var usersService service.UserService

func NewUsersHandler(us service.UserService) func(w http.ResponseWriter, r *http.Request) {
	usersService = us
	return UsersHandler
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "POST %q", html.EscapeString(r.URL.Path))
		break
	case http.MethodGet:
		id := r.URL.Path[len(UsersPath):]

		if id == "" {
			getAllUsers(w)
		} else {
			getUser(id, w)
		}

		break
	}
}

func getUser(id string, w http.ResponseWriter) {
	us := usersService.Get(id)
	ur := convertUserFromService(us)

	userJson, err := json.Marshal(ur)
	if err != nil {
		panic(err) // TODO: error handler
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}


func getAllUsers(w http.ResponseWriter) {

	// begin
	b := func(err error) {
		// TODO: check error
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"Users\":["))
	}

	// next
	first := true
	n := func(us *service.User, err error) {
		ur := convertUserFromService(us)

		resJson, err := json.Marshal(ur)
		if err != nil {
			panic(err) // TODO: error handler
		}

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

	usersService.GetAll(b, n, e)
}


// Conversion Service to REST object
func convertUserFromService(us *service.User) *User {
	return &User{
		Uri: MakeIdPath(UsersPath,us.Id),
		Name: us.Name,
	}
}