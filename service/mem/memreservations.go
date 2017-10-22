package memservice

import (
	"fmt"
	"strconv"
	"resys/service"
)

type MemReservationService struct {
	reservations []service.Reservation
}


func (rs *MemReservationService) nextReservationId() int {
	n := 0
	for _, r := range(rs.reservations) {
		if r.Id > n {
			n = r.Id
		}
	}
	return n+1
}


func (rs *MemReservationService) Create(reservation *service.Reservation) (int, error) {
	fmt.Println(reservation)

	// TODO: Check for overlaps

	reservation.Id = rs.nextReservationId()
	rs.reservations = append(rs.reservations, *reservation)

	return reservation.Id, nil
}


func (rs *MemReservationService) GetAll(begin func(err error),
 									next func(reservation *service.Reservation, err error),
									end func(err error)) {
	begin(nil)
	for _, r := range(rs.reservations) {
		next(&r, nil)
	}
	end(nil)
}


func (rs *MemReservationService) Get(id string) *service.Reservation {
	i, err := strconv.Atoi(id);
	if err != nil {
		return nil
	}

	for _, r := range(rs.reservations) {
		if r.Id == i {
			return &r
		}
	}
	return nil
}
