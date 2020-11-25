package image

import "encoding/csv"
import "strings"
import "fmt"

func ReadContent(user string) *ImageManager {
	reader := csv.NewReader(strings.NewReader(""))
	reader.Comma = '|'

	images, err := reader.ReadAll()
	if err != nil {
		// Error stuff
		fmt.Println(err)
	}

	for _, img := range images {
		fmt.Println("Name", img[0])
		fmt.Println("Datum", img[1])
		fmt.Println("Hash", img[1])
	}

	return &ImageManager{}
}

func WriteContent(user string, imgs *ImageManager) error {
	return nil
}