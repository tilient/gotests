package main

import (
	"fmt"
	. "github.com/tilient/gotests/genlist/list"
)

func main() {
	var lst *IntList
	lst = lst.Add(123).Add(321).Add(777)
	fmt.Println(lst)
}
