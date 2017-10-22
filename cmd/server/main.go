package main

import (
	"resys/rest"
	"log"
	"net/http"
	"resys/service/mem"
)


func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func tracer(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		h.ServeHTTP(w, r)
		log.Println("After")
	})
}


func main() {

	// In memory DB
	resourcesService := memservice.DemoResources()
	usersService := &memservice.MemUserService{}
	reservationsService := &memservice.MemReservationService{}

	http.HandleFunc(rest.ResourcesPath, use(rest.NewResourcesHandler(resourcesService), tracer))
	http.HandleFunc(rest.UsersPath, rest.NewUsersHandler(usersService))
	http.HandleFunc(rest.ReservationsPath, rest.NewReservationsHandler(reservationsService, resourcesService))

	http.HandleFunc("/login", loginHandler) // TODO: SSL for login
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("templates"))))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

