package memservice

import (
	"fmt"
	"strconv"
	"resys/service"
)


type MemUserService struct {}


func (rs *MemUserService) GetAll(begin func(err error),
 									next func(user *service.User, err error),
									end func(err error)) {
	fmt.Println("Getting all")
	begin(nil)
	for i:=0; i<10; i++ {
		r := &service.User{Id:i, Name:"U"+strconv.Itoa(i),}
		next(r, nil)
	}
	end(nil)
}


func (rs *MemUserService) Get(id string) *service.User {
	fmt.Println("Getting: ", id)
	return &service.User{Id:1, Name:"U1",}
}
