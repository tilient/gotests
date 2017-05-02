package list

import (
	"fmt"
	"github.com/joeshaw/gengen/generic"
)

//---------------------------------------------------------------------

type т generic.T

type тList struct {
	value т
	next  *тList
}

//---------------------------------------------------------------------

func (lst *тList) isEmpty() bool {
	return lst == nil
}

func (lst *тList) Add(val т) *тList {
	return &тList{val, lst}
}

//---------------------------------------------------------------------

func (lst *тList) String() string {
	return "[" + lst.innerString() + "]"
}

func (lst *тList) innerString() string {
	if lst.isEmpty() {
		return ""
	}
	if lst.next.isEmpty() {
		return fmt.Sprint(lst.value)
	}
	return fmt.Sprint(lst.value) + "," + lst.next.innerString()
}

//---------------------------------------------------------------------
