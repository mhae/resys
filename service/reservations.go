package service

import (
	"time"
)


type Reservation struct {
	Id         int
	UserId     int
	ResourceId int
	ResourceName string
	Start      time.Time
	End        time.Time
}

type ReservationService interface {
	Create(reservation *Reservation) (int, error)
	Get(id string) *Reservation
	GetAll(begin func(err error), next func(reservation *Reservation, err error), end func(err error))
}

