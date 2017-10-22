package main

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"net/http"
)

const SESSION = "session"

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		h.ServeHTTP(w, r)
	}
}

func mustAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(SESSION)
		if err == http.ErrNoCookie {
			// not authenticated
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		if err != nil {
			// some other error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Do authorization

		fmt.Println(getUserName(r))

		// success - call the next handler
		h.ServeHTTP(w, r)
	})
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode(SESSION, value); err == nil {
		cookie := &http.Cookie{
			Name:  SESSION,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func getUserName(request *http.Request) (userName string) {

	if cookie, err := request.Cookie(SESSION); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(SESSION, cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   SESSION,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	email := request.FormValue("inputEmail")
	pass := request.FormValue("inputPassword")
	redirectTarget := "/"
	if email != "" && pass != "" {
		//
		// TODO: Check credentials
		//
		setSession(email, response)
		redirectTarget = "/templates/main.html"
	} else {
		redirectTarget = "/templates/login.html"
	}

	http.Redirect(response, request, redirectTarget, 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/templates/login.html", 302)
}
