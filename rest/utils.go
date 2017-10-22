package rest

import (
	"strconv"
	"reflect"
)

func MakeIdPath(p string, id int) string {
	return p+strconv.Itoa(id)
}

// Simple (flat) struct mapper. Copies fields from src struct to dest struct that have the same name.
// Probably not very performant
func FlatMapper(src interface{}, dest interface{})  {

	st := reflect.TypeOf(src).Elem()
	dv := reflect.ValueOf(dest).Elem()
	sv := reflect.ValueOf(src).Elem()

	if st.Kind() != reflect.Struct || dv.Kind() != reflect.Struct {
		return
	}

	for i:=0; i<st.NumField(); i++ {

		df := dv.FieldByName(st.Field(i).Name)

		if !(df.IsValid() && df.CanSet()) {
			continue
		}

		if df.Kind() != st.Field(i).Type.Kind() {
			return
		}

		df.Set(sv.Field(i))
	}

}