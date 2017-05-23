package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"sort"
)

const pieceLen = 128

type (
	charOcc struct {
		char byte
		occ  float64
	}
	charOccs          []charOcc
	stringCharOccsMap map[[16]byte]charOccs
)

func main() {
	content, _ := ioutil.ReadFile("captmidn.txt")
	text := string(content)
	scomap := make(stringCharOccsMap)
	for t := 0; t < 4; t++ {
		for ix := 0; ix < len(text)-pieceLen; ix++ {
			piece := text[ix : ix+pieceLen]
			md5piece := md5.Sum([]byte(piece))
			realNextChar := text[ix+pieceLen]
			scomap.learn(md5piece, realNextChar)
		}
	}
	str := "He had dropped out of a management engineering course at Worcester Polytechnical Institute in Massachusetts after two years, but his first job was installing satellite TV dishes."
	str = str[0:pieceLen]
	fmt.Println(str, "...")
	fmt.Println("---------------------------------")
	for t := 0; t < 1000; t++ {
		md5str := md5.Sum([]byte(str))
		nextChar := scomap.predict(md5str)
		fmt.Printf("%c", nextChar)
		str = str[1:] + string(nextChar)
	}
	fmt.Println("\n---------------------------------")
}

func (m *stringCharOccsMap) predict(str [16]byte) byte {
	chOccs := (*m)[str]
	if len(chOccs) > 0 {
		return chOccs[0].char
	}
	return '?'
}

func (m *stringCharOccsMap) learn(str [16]byte, b byte) {
	chars := (*m)[str]
	if len(chars) == 0 {
		chars = charOccs{}
	}
	if len(chars) > 8 {
		chars = chars[0:7]
	}
	pos, found := chars.charPos(b)
	if found {
		chars[pos].occ += 1.0
	} else {
		chars = append(chars, charOcc{b, 1.0})
	}
	sort.Slice(chars, func(i, j int) bool {
		return chars[i].occ > chars[j].occ
	})
	(*m)[str] = chars
}

func (co *charOccs) charPos(ch byte) (int, bool) {
	for ix, c := range *co {
		if c.char == ch {
			return ix, true
		}
	}
	return 0, false
}
