package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	dirName := "/mnt/c/wiffel/public/photoKashbah"
	filepath.Walk(dirName, visit)
}

func visit(path string, f os.FileInfo, err error) error {
	if isRawPicture(path) {
		exists := companionJpeg(path)
		if exists {
			fmt.Println("delete", path)
			os.Remove(path)
		}
	}
	return nil
}

func isRawPicture(path string) bool {
	rawExtensions := map[string]bool{
		".CR2": true,
		".cr2": true,
		".dng": true,
		".DNG": true,
	}
	ext := filepath.Ext(path)
	return rawExtensions[ext]
}

func companionJpeg(path string) bool {
	jpgExtensions := []string{".jpg", ".JPG", ".jpeg", ".JPEG"}
	basename := fileBasename(path)
	for _, ext := range jpgExtensions {
		companionFilename := basename + ext
		if fileExists(companionFilename) {
			return true
		}
	}
	return false
}

func fileBasename(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func exifStuff() {
	fname := "xxx"
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	tm, _ := x.DateTime()
	fmt.Println(fname, "@", tm)
}
