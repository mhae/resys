package rest

import (
	"net/http"
	"resys/service"
	"encoding/json"
	"time"
	"strconv"
	"fmt"
)

// curl -H "Content-Type: application/json" -X POST -d '{"UserUri":"/users/1","ResourceUri":"/resources/1","Start":"2017-10-08T22:57:57+00:00","End":"2017-10-09T22:57:57+00:00"}' http://localhost:8080/reservations/
type NewReservation struct {
	UserUri string
	ResourceUri string
	Start time.Time
	End time.Time
}

// GET/PUT
type Reservation struct {
	Uri string // read-only
	UserUri string
	ResourceUri string
	ResourceName string // read-only
	Start time.Time
	End time.Time
}

const ReservationsPath = "/reservations/"


// TODO Inject?
var reservationsService service.ReservationService
var resourcesService service.ResourceService

func NewReservationsHandler(rv service.ReservationService, rs service.ResourceService) func(w http.ResponseWriter, r *http.Request) {
	reservationsService = rv
	resourcesService = rs
	return ReservationsHandler
}

func ReservationsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		createNewReservation(w, r)
		break
	case http.MethodGet:
		id := r.URL.Path[len(ReservationsPath):]

		if id == "" {
			getAllReservations(w)
		} else {
			getReservation(id, w)
		}
		break
	}
}

func createNewReservation(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var nr NewReservation
	err := decoder.Decode(&nr)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if nr.validate(w) != nil {
		return
	}

	snr := convertNewReservationToService(&nr)

	id, err := reservationsService.Create(snr)

	// TODO: reply ... potential reservation conflict

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Location", ReservationsPath+"/"+strconv.Itoa(id))
}

func convertNewReservationToService(nrv *NewReservation) *service.Reservation {
	userId, _ := strconv.Atoi(nrv.UserUri[len(UsersPath):])
	resourceId, _ := strconv.Atoi(nrv.ResourceUri[len(ResourcesPath):])

	// find resource name
	r := resourcesService.Get(nrv.ResourceUri[len(ResourcesPath):])
	fmt.Println(r)
	// TODO error check

	rv := &service.Reservation{UserId:userId, ResourceId:resourceId, ResourceName:r.Name }
	FlatMapper(nrv, rv)
	return rv
}

func convertReservationFromService(srv *service.Reservation) *Reservation {
	rrv := &Reservation{
		Uri: MakeIdPath(ReservationsPath, srv.Id),
		UserUri: MakeIdPath(UsersPath, srv.UserId),
		ResourceUri:MakeIdPath(ResourcesPath, srv.ResourceId) }

	FlatMapper(srv, rrv)

	return rrv
}

func (nr *NewReservation) validate(w http.ResponseWriter) error {
	// TODO
	return nil
}

func getReservation(id string, w http.ResponseWriter) {
	rs := reservationsService.Get(id)
	rr := convertReservationFromService(rs)

	reservationJson, err := json.Marshal(rr)
	if err != nil {
		panic(err) // TODO: error handler
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(reservationJson)
}


func getAllReservations(w http.ResponseWriter) {

	// begin
	b := func(err error) {
		// TODO: check error
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"Reservations\":["))
	}

	// next
	first := true
	n := func(rs *service.Reservation, err error) {
		rr := convertReservationFromService(rs)

		reservationJson, err := json.Marshal(rr)
		if err != nil {
			panic(err) // TODO: error handler
		}

		if !first {
			w.Write([]byte(","))
		}
		w.Write(reservationJson)
		first = false
	}

	// end
	e := func(err error) {
		w.Write([]byte("]}"))
	}

	reservationsService.GetAll(b, n, e)
}

