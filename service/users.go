package service

import (

)


type User struct {
	Id int
	Name string
}

type UserService interface  {
	GetAll(begin func(err error), next func(user *User, err error), end func(err error))
	Get(id string) *User
}
