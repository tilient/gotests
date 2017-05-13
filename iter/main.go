package main

import "fmt"

// -----------------------------------------------------------

func main() {
	filenames := filenames("file_%c_%02d_%c.txt",
		intRanges{{'d', 'f'}, {3, 5}, {'x', 'z'}})
	for _, filename := range filenames {
		fmt.Println(filename)
	}
}

// --- format filenames --------------------------------------

func filenames(format string, ranges intRanges) []string {
	result := []string{}
	for _, args := range combinations(ranges) {
		filename := fmt.Sprintf(format, args.asInterfaceList()...)
		result = append(result, filename)
	}
	return result
}

// --- ranges ------------------------------------------------

type (
	intRange  [2]int
	intRanges []intRange
	intList   []int
	intLists  []intList
)

func combinations(ranges intRanges) intLists {
	return combine(ranges, intList{})
}

func combine(ranges intRanges, args intList) intLists {
	if len(ranges) == 0 {
		return intLists{args}
	}
	lst := intLists{}
	first := ranges[0]
	rest := ranges[1:]
	for v := first[0]; v <= first[1]; v++ {
		lst = append(lst, combine(rest, append(args, v))...)
	}
	return lst
}

// --- tools -------------------------------------------------

func (iLst intList) asInterfaceList() []interface{} {
	lst := make([]interface{}, len(iLst))
	for ix, i := range iLst {
		lst[ix] = i
	}
	return lst
}

// -----------------------------------------------------------
