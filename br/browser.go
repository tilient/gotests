package main

import (
	"os"
	"os/exec"
	"runtime"
)

func main() {
	Open("www.google.be")
}

func Commands() [][]string {
	var cmds [][]string
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, []string{exe})
	}
	switch runtime.GOOS {
	case "darwin":
		cmds = append(cmds, []string{"/usr/bin/open"})
	case "windows":
		cmds = append(cmds, []string{"cmd", "/c", "start"})
	default:
		cmds = append(cmds, []string{"xdg-open"})
	}
	cmds = append(cmds,
		[]string{"chrome"},
		[]string{"google-chrome"},
		[]string{"firefox"})
	return cmds
}

func Open(url string) bool {
	for _, args := range Commands() {
		println(args[0])
		cmd := exec.Command(args[0], append(args[1:], url)...)
		if cmd.Start() == nil {
			return true
		}
	}
	return false
}
