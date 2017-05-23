package main

import "fmt"

// ------------------------------------------------------------

func main() {
	filenameTemplate := "%v_%c%02d_T%04dF%03d%vZ%02dC%02d.tif"
	barcode := "110000124800"
	wellRows := fromTo('A', 'P')
	wellColumns := fromTo(1, 24)
	times := fromTo(1, 1)
	sites := fromTo(0, 2)
	zoomLevels := fromTo(0, 2)
	imageSets := allImageSetsFor(barcode,
		wellRows, wellColumns, times, sites, zoomLevels)

	channels := fromTo(1, 3)
	actions := []string{"L01A01", "L01A02", "L01A01"}
	filenames := []string{}
	for _, channel := range channels {
		fn := fmt.Sprintf(filenameTemplate,
			barcode, row, column, time, site,
			actions[channel-1], zoom, channel)
		filenames = append(filenames, fn)
	}
	if ok, missing := allFilesExist(filenames); ok {
		fmt.Println(filenames)
	} else {
		fmt.Println("missing:", missing)
	}
}

// --- ranges  ------------------------------------------------

func fromTo(from, to int) []int {
	result := make([]int, 1+to-from)
	for i := from; i <= to; i++ {
		result[i-from] = i
	}
	return result
}

// --- imageSets ----------------------------------------------

type imageSet struct {
	barcode string
	row     int
	column  int
	time    int
	site    int
	zoom    int
}
type imageSets []imageSet

func allImageSetsFor(barcode string,
	wellRows, wellColumns, times, sites, zoomLevel []int) {
	result := []imageSet{}
	for _, row := range wellRows {
		for _, column := range wellColumns {
			for _, time := range times {
				for _, site := range sites {
					for _, zoom := range zoomLevels {
						result = append(result, imageSet{
							barcode, row, column, time, site, zoom})
					}
				}
			}
		}
	}
	return result
}

func allFilesExist(filenames []string) (bool, []string) {
	result := true
	missing := []string{}
	for _, filename := range filenames {
		if filename == "110000124800_P24_T0001F002L01A01Z02C01.tif" {
			missing = append(missing, filename)
			result = false
		}
	}
	return result, missing
}

func (sets imageSets) checkValidity(actions []string) (imageSets, imageSets) {
	valid := imageSets{}
	invalid := imageSets{}

	return valid, invalid
}

// ------------------------------------------------------------
