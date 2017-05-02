package list

import (
	"fmt"
	"github.com/cheekybits/genny/generic"
)

//---------------------------------------------------------------------
//go:generate genny -in=$GOFILE -out=builtins$GOFILE gen "Type=BUILTINS"
//---------------------------------------------------------------------

type Type generic.Type

type TypeList struct {
	value Type
	next  *TypeList
}

//---------------------------------------------------------------------

func (lst *TypeList) isEmpty() bool {
	return lst == nil
}

func (lst *TypeList) Add(val Type) *TypeList {
	return &TypeList{val, lst}
}

//---------------------------------------------------------------------

func (lst *TypeList) String() string {
	return "[" + lst.innerString() + "]"
}

func (lst *TypeList) innerString() string {
	if lst.isEmpty() {
		return ""
	}
	if lst.next.isEmpty() {
		return fmt.Sprint(lst.value)
	}
	return fmt.Sprint(lst.value) + "," + lst.next.innerString()
}

//---------------------------------------------------------------------
