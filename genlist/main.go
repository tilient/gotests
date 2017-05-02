package main

import (
	"fmt"
	. "github.com/tilient/gotests/genlist/list"
)

//---------------------------------------------------------------------
// go:generate genny -in=list/list.go -out=list/intlist.go gen "т=int"
//go:generate gengen github.com/tilient/gotests/genlist/list/list2.go int
//---------------------------------------------------------------------

func main() {
	var lst *IntList
	lst = lst.Add(123).Add(321).Add(777)
	fmt.Println(lst)
}
