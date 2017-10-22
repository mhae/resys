package memservice

import (
	"resys/service"
	"strconv"
	"fmt"
)



type MemResourceService struct {
	resources []*service.Resource
}


func DemoResources() service.ResourceService {
	rs := &MemResourceService{ resources: make([]*service.Resource, 0)}
	for i:=0; i<10; i++ {
		r := &service.Resource{Id:i, Name:"Resource"+strconv.Itoa(i), Tags:[]string{"A tag"}}
		rs.resources = append(rs.resources, r)
	}
	return rs
}


func (rs *MemResourceService) GetAll(begin func(err error),
 									next func(resource *service.Resource, err error),
									end func(err error)) {
	fmt.Println("Getting all")
	begin(nil)
	for i:=0; i<len(rs.resources); i++ {
		next(rs.resources[i], nil)
	}
	end(nil)
}


// TODO: Error handling
func (rs *MemResourceService) Get(id string) *service.Resource {
	fmt.Println("Getting: ", id)
	i, _ := strconv.Atoi(id)
	return rs.resources[i]
}

