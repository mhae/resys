package service

import (

)


type Resource struct {
	Id int
	Name string
	Tags []string
}


type ResourceService interface {
	GetAll(begin func(err error), next func(resource *Resource, err error), end func(err error))
	Get(id string) *Resource
}

