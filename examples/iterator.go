package main

import (
	"fmt"
	"reflect"
)

type Iterable interface {
	Len() error
}

type Iterator struct {
	collection interface{}
	iterN      int
}

func (it *Iterator) next() interface{} {
	coll := reflect.ValueOf(it.collection)

	if coll.Kind() == reflect.Ptr {
		coll = coll.Elem()
	}

	if coll.Kind() != reflect.Slice {
		panic("Must be slice")
	}

	if it.iterN >= coll.Len() { return nil }

	return coll.Index(it.iterN).Interface()

}

func forEach(iterator Iterator, f func (a ...interface{})) {
	for  {
		v := iterator.next()
		if v == nil { break }
		iterator.iterN++
		f(v)
	}

}


func main() {
	iter := Iterator{[]int{1, 2, 3}, 0}

	f1 := func(x ...interface{}) {
		fmt.Println(x)
	}

	forEach(iter, f1)
}
