package main

import "fmt"

type seq struct {
	from int
	to   int
}

func expandAll(ch chan string, format string, seqs []seq, args ...interface{}) {
	if len(seqs) == 0 {
		ch <- fmt.Sprintf(format, args...)
	} else {
		for v := seqs[0].from; v <= seqs[0].to; v++ {
			expandAll(ch, format, seqs[1:], append(args, v)...)
		}
	}
}

func expand(format string, seqs ...seq) chan string {
	ch := make(chan string)
	go func() {
		expandAll(ch, format, seqs)
		close(ch)
	}()
	return ch
}

func main() {
	filenames := expand("file_%c_%02d.txt", seq{'d', 'f'}, seq{3, 5})
	for filename := range filenames {
		fmt.Println(filename)
	}
}
