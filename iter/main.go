package main

import "fmt"

// ----------------------------------------------------------

func main() {
	filenames := formatString("file_%c%c_%02d_%c.txt").expand(
		intervals{{'D', 'E'}, {'A', 'B'}, {0, 1}, {'a', 'c'}})
	for _, fn := range filenames {
		fmt.Println(fn)
	}
}

// --- template ---------------------------------------------

type formatString string

func (fs formatString) expand(intervals intervals) []string {
	list2string := func(l list) string {
		return fmt.Sprintf(string(fs), l.asInterfaceList()...)
	}
	return intervals.combinations().mapIt(list2string)
}

// --- intervals --------------------------------------------

type (
	interval struct {
		low  int
		high int
	}
	intervals []interval
)

func (r interval) collect(f func(i int) lists) lists {
	result := lists{}
	for i := r.low; i <= r.high; i++ {
		result = append(result, f(i)...)
	}
	return result
}

func (intervals intervals) combinations(args ...int) lists {
	if len(intervals) == 0 {
		return lists{args}
	}
	head := intervals[0]
	tail := intervals[1:]
	tailCombinations := func(i int) lists {
		return tail.combinations(append(args, i)...)
	}
	return head.collect(tailCombinations)
}

// --- lists ------------------------------------------------

type (
	list  []int
	lists []list
)

func (lst list) asInterfaceList() []interface{} {
	result := make([]interface{}, len(lst))
	for ix, v := range lst {
		result[ix] = v
	}
	return result
}

func (lsts lists) mapIt(f func(list) string) []string {
	result := []string{}
	for _, lst := range lsts {
		result = append(result, f(lst))
	}
	return result
}

// ----------------------------------------------------------
