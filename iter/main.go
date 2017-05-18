package main

import "fmt"

// ------------------------------------------------------------

func main() {
	filenameParts := formatString("img_%c_%02d_").expand(
		intervals{{'A', 'C'}, {0, 2}})
	for _, fn := range filenameParts {
		filenames := stringPart(fn).combinations(
			[]string{"A01B_C01.img", "A02B_C02.img", "A01C_C03.img"})
		if ok, missing := allFilesExist(filenames); ok {
			fmt.Println(filenames)
		} else {
			fmt.Println("missing:", missing)
		}
	}
}

func allFilesExist(filenames []string) (bool, []string) {
	result := true
	missing := []string{}
	for _, filename := range filenames {
		if filename == "img_C_01_A02B_C02.img" {
			missing = append(missing, filename)
			result = false
		}
	}
	return result, missing
}

// --- strings ------------------------------------------------

type stringPart string

func (str stringPart) combinations(strings []string) []string {
	res := []string{}
	for _, ext := range strings {
		res = append(res, string(str)+ext)
	}
	return res
}

// --- template -----------------------------------------------

type formatString string

func (fs formatString) expand(intervals intervals) []string {
	list2string := func(l list) string {
		return fmt.Sprintf(string(fs), l.asInterfaceList()...)
	}
	return intervals.combinations().mapIt(list2string)
}

// --- intervals ----------------------------------------------

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

// --- lists --------------------------------------------------

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

// ------------------------------------------------------------
