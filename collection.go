package eel

import "reflect"

type IIDable interface {
	GetId(idType string) int
}

//// Definition of Map function
//type mapFunc func(interface{}) interface{}
//
//// Map maps the function onto the array
//func MapAll(fn mapFunc, array interface{}) []interface{} {
//	val := reflect.ValueOf(array)
//	outputArray := make([]interface{}, val.Len())
//	for i := 0; i < val.Len(); i++ {
//		outputArray[i] = fn(val.Index(i).Interface())
//	}
//	return outputArray
//}
//
//// Definition of Filter function
//type filterFunc func(interface{}) bool
//
//// Filter filters the array based on the predicate
//func Filter(fn filterFunc, array interface{}) []interface{} {
//	val := reflect.ValueOf(array)
//	var outputArray []interface{}
//	for i := 0; i < val.Len(); i++ {
//		if fn(val.Index(i).Interface()) {
//			outputArray = append(outputArray, val.Index(i).Interface())
//		}
//	}
//	return outputArray
//}
//
//// Definition of Foldl function
//type foldlFunc func(interface{}, interface{}) interface{}
//
//// Folds left the array values (reduction) based on the function
//func Foldl(fn foldlFunc, array interface{}, accumulator interface{}) interface{} {
//	val := reflect.ValueOf(array)
//	var result = accumulator
//	for i := 0; i < val.Len(); i++ {
//		result = fn(val.Index(i).Interface(), result)
//	}
//	return result
//}

type IntMapFunc func(interface{}) int

// Map maps the function onto the array
func MapInt(fn IntMapFunc, array interface{}) []int {
	val := reflect.ValueOf(array)
	outputArray := make([]int, val.Len())
	for i := 0; i < val.Len(); i++ {
		outputArray[i] = fn(val.Index(i).Interface())
	}
	return outputArray
}

//ExtractIds 从objs中获取id，组成id的集合
func ExtractIds(array interface{}, idType string) []int {
	val := reflect.ValueOf(array)
	outputArray := make([]int, val.Len())
	for i := 0; i < val.Len(); i++ {
		obj := val.Index(i).Interface()
		id := obj.(IIDable).GetId(idType)
		outputArray[i] = id
	}
	return outputArray
}

