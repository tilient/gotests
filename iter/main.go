package main

import "fmt"

func interval(from, to int) chan int {
	ch := make(chan int)
	go func() {
		for i := from; i <= to; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func filenames() <-chan string {
	ch := make(chan string)
	go func() {
		for c := range interval('a', 'e') {
			for i := range interval(0, 2) {
				ch <- fmt.Sprintf("file_%c_%02d.txt", c, i)
			}
		}
		close(ch)
	}()
	return ch
}

func main() {
	for filename := range filenames() {
		go fmt.Println(filename)
	}
}
