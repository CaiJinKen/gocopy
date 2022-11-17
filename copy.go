package gocopy

import (
	"errors"
	"reflect"
)

// NewFrom get a deep copy from src
func NewFrom(src interface{}) (dst interface{}) {
	if src == nil {
		return nil
	}

	srcT := reflect.TypeOf(src)
	srcV := reflect.ValueOf(src)

	dstVal := reflect.New(srcV.Type()).Elem()
	switch srcT.Kind() {
	case reflect.Bool:
		dstVal.SetBool(srcV.Bool())
	case reflect.String:
		dstVal.SetString(srcV.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dstVal.SetInt(srcV.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		dstVal.SetUint(srcV.Uint())
	case reflect.Float32, reflect.Float64:
		dstVal.SetFloat(srcV.Float())
	case reflect.Complex64, reflect.Complex128:
		dstVal.SetComplex(srcV.Complex())
	case reflect.Slice:
		if !srcV.IsNil() {
			// dstVal.Set(reflect.AppendSlice(dstVal, srcV)) // it`s same in underline
			for i := 0; i < srcV.Len(); i++ {
				dstVal = reflect.Append(dstVal, reflect.ValueOf(NewFrom(srcV.Index(i).Interface())))
			}
		}
	case reflect.Array:
		if srcV.Len() > 0 {
			reflect.Copy(dstVal, srcV)
		}
	case reflect.Map:
		iterator := srcV.MapRange()
		for iterator.Next() {
			if dstVal.IsNil() {
				mt := reflect.MapOf(iterator.Key().Type(), iterator.Value().Type())
				mp := reflect.MakeMapWithSize(mt, dstVal.Len())
				// mp := reflect.MakeMap(mt)
				dstVal.Set(mp)
			}
			dstVal.SetMapIndex(iterator.Key(), reflect.ValueOf(NewFrom(iterator.Value().Interface())))
		}
	case reflect.Ptr:
		dstVal = reflect.New(srcV.Elem().Type()).Elem()
		v := NewFrom(srcV.Elem().Interface())
		dstVal.Set(reflect.ValueOf(v))
		dstVal = dstVal.Addr()
	case reflect.Interface:
		dstVal = reflect.New(srcV.Elem().Type()).Elem()
		v := NewFrom(srcV.Elem())
		dstVal.Set(reflect.ValueOf(v))
	case reflect.Struct:
		for i := 0; i < srcT.NumField(); i++ {
			srcField := srcV.Field(i)
			if srcField.Interface() == nil {
				continue
			}
			dstVal.Field(i).Set(reflect.ValueOf(NewFrom(srcField.Interface())))
		}

	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		dstVal.Set(srcV)

	}
	dst = dstVal.Interface()
	return
}

// Update assign value from src to dst when src and dst has same field name and data type
// src must be a struct or a pointer of struct
// dst must be a pointer of struct
func Update(src, dst interface{}) error {
	if src == nil || dst == nil {
		return errors.New("both src and dst can`t be nil value")
	}
	dstV := reflect.ValueOf(dst)
	dstT := reflect.TypeOf(dst)
	if dstT.Kind() != reflect.Ptr || reflect.TypeOf(dstV.Elem().Interface()).Kind() != reflect.Struct {
		return errors.New("dst should be a pointer of struct")
	}

	srcT := reflect.TypeOf(src)
	srcV := reflect.ValueOf(src)
	if srcV.Kind() == reflect.Ptr {
		srcT = reflect.TypeOf(srcV.Elem().Interface())
		srcV = reflect.ValueOf(srcV.Elem().Interface())
	}

	if srcT.Kind() != reflect.Struct {
		return errors.New("dst should be a struct or pointer of struct")
	}

	for i := 0; i < dstT.Elem().NumField(); i++ {
		dstF := dstT.Elem().Field(i)
		srcF, ok := srcT.FieldByName(dstF.Name)
		if !ok {
			continue
		}

		if dstF.Type.Kind().String() != srcF.Type.Kind().String() {
			continue
		}
		dstV.Elem().Field(i).Set(srcV.FieldByName(dstF.Name))
	}

	return nil
}
