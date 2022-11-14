package gocopy

import (
    "reflect"
)

// NewFrom get a copy from src
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
