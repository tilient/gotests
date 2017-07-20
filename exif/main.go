package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	//removeDNGsWithJPGsInPhotoKashbah()
	//removeJPGsWithDNGsInOriginals()
	//removeJPGsWithRawInOriginals()
	//findPictureFilesWithoutExifDate()
	findCandidateDirectoriesToConvert()
}

func findCandidateDirectoriesToConvert() {
	rootDir := "/mnt/c/wiffel/pictures/originals/"
	subDirs, _ := ioutil.ReadDir(rootDir)
	for _, subDir := range subDirs {
		dir := rootDir + subDir.Name() + "/"
		files, _ := ioutil.ReadDir(dir)
		for _, f := range files {
			d := dir + f.Name()
			if dirContainsCR2(d) && !dirContainsXmp(d) {
				fmt.Println("**", d)
			}
		}
	}

}

func dirContainsCR2(dirName string) bool {
	result := false
	filepath.Walk(dirName,
		func(path string, f os.FileInfo, err error) error {
			//fmt.Println("--", path)
			if isCR2(path) {
				result = true
			}
			return nil
		})
	return result
}

func dirContainsXmp(dirName string) bool {
	result := false
	filepath.Walk(dirName,
		func(path string, f os.FileInfo, err error) error {
			//fmt.Println("--", path)
			if isXmp(path) {
				result = true
			}
			return nil
		})
	return result
}

func findPictureFilesWithoutExifDate() {
	dirName := "/mnt/c/wiffel/public/photoKashbah"
	filepath.Walk(dirName,
		func(path string, f os.FileInfo, err error) error {
			if isPicture(path) {
				if !hasExifDate(path) {
					fmt.Println(path)
				}
			}
			return nil
		})
}

func removeJPGsWithRawInOriginals() {
	dirName := "/mnt/c/wiffel/pictures/originals"
	filepath.Walk(dirName,
		func(path string, f os.FileInfo, err error) error {
			if isRawPicture(path) {
				hasCompanionJpeg, jpgPath := companionJpeg(path)
				if hasCompanionJpeg {
					fmt.Println(jpgPath)
					os.Remove(jpgPath)
				}
			}
			return nil
		})
}

func removeJPGsWithDNGsInOriginals() {
	dirName := "/mnt/c/wiffel/pictures/originals"
	filepath.Walk(dirName,
		func(path string, f os.FileInfo, err error) error {
			if isDngPicture(path) {
				hasCompanionJpeg, jpgPath := companionJpeg(path)
				if hasCompanionJpeg {
					fmt.Println(jpgPath)
					os.Remove(jpgPath)
				}
			}
			return nil
		})
}

func removeDNGsWithJPGsInPhotoKashbah() {
	dirName := "/mnt/c/wiffel/public/photoKashbah"
	filepath.Walk(dirName,
		func(path string, f os.FileInfo, err error) error {
			if isRawPicture(path) {
				hasCompanionJpeg, _ := companionJpeg(path)
				if hasCompanionJpeg {
					fmt.Println("delete", path)
					os.Remove(path)
				}
			}
			return nil
		})
}

func isCR2(path string) bool {
	ext := filepath.Ext(path)
	return ((ext == ".cr2") || (ext == ".CR2"))
}

func isXmp(path string) bool {
	ext := filepath.Ext(path)
	return ((ext == ".xmp") || (ext == ".XMP"))
}

func isDngPicture(path string) bool {
	ext := filepath.Ext(path)
	return ((ext == ".dng") || (ext == ".DNG"))
}

func isRawPicture(path string) bool {
	rawExtensions := map[string]bool{
		".dng": true,
		".DNG": true,
		".cr2": true,
		".CR2": true,
	}
	ext := filepath.Ext(path)
	return rawExtensions[ext]
}

func isPicture(path string) bool {
	rawExtensions := map[string]bool{
		".jpg":  true,
		".JPG":  true,
		".jpeg": true,
		".JPEG": true,
		".png":  true,
		".PNG":  true,
		".dng":  true,
		".DNG":  true,
		".cr2":  true,
		".CR2":  true,
	}
	ext := filepath.Ext(path)
	return rawExtensions[ext]
}

func companionJpeg(path string) (bool, string) {
	jpgExtensions := []string{
		".jpg", ".JPG", ".jpeg", ".JPEG",
		"-1.jpg", "-1.JPG", "-1.jpeg", "-1.JPEG",
		"-2.jpg", "-2.JPG", "-2.jpeg", "-2.JPEG"}
	basename := fileBasename(path)
	for _, ext := range jpgExtensions {
		companionFilename := basename + ext
		if fileExists(companionFilename) {
			return true, companionFilename
		}
	}
	return false, ""
}

func fileBasename(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasExifDate(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	d, err := exif.Decode(f)
	if err != nil {
		return false
	}
	_, err = d.DateTime()
	if err != nil {
		return false
	}
	return true
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
