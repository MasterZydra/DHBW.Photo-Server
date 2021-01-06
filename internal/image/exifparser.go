package image

import (
	"bufio"
	"encoding/binary"
	"io"
	"regexp"
	"strings"
)

func parseRawExifDataFromFile(file io.Reader) ([]byte, error) {
	// EXIF-Data is in APP1-Section
	// Marker for APP1-Section: 0xFFE1
	// read App1-Section data
	app1Section, err := readSectionData(0xE1, bufio.NewReader(file))

	if err != nil {
		return nil, err
	}

	return app1Section, nil
}

func readSectionData (marker byte, br *bufio.Reader) ([]byte, error) {
	// find marker in reader and go there
	var dataLength int = 0
	var markerFound bool = false
	for !markerFound {
		c, err := br.ReadByte()
		if err != nil {
			return nil, err
		} else if c == marker {
			// marker found
			markerFound = true
			// first two bytes after marker are size of APP1-Section-Data
			bytes := make([]byte, 2)
			for k, _ := range bytes {
				c, err := br.ReadByte()
				if err != nil {
					return nil, err
				}
				bytes[k] = c
			}
			// set dataLength
			dataLength = int(binary.BigEndian.Uint16(bytes)) - 2
		}
	}

	// read data

	alreadyRead := 0
	data := make([]byte, 0)

	for alreadyRead < dataLength {
		read := make([]byte, dataLength-alreadyRead)
		n, err := br.Read(read)
		alreadyRead += n
		if err != nil && n < dataLength {
			return nil, err
		}
		data = append(data, read[:n]...)
	}

	return data, nil
}

func getDateFromData (data []byte) string {
	// DateTime Format: YYYY:MM:DD HH:MM:SS
	dateRegex := regexp.MustCompile(`\d{4}\:(0[1-9]|1[012])\:(0[1-9]|[12][0-9]|3[01]) ([01][0-9]|2[0-3]):([0-6][0-9]):([0-6][0-9])`)
	// There are three Tags with this format in EXIF-Data: DateTime, DateTimeOriginal, DateTimeDigitized
	// all three Tags usually contain the same value
	var date string = string(dateRegex.Find(data))
	// wanted format: yyyy-mm-dd hh:mm:ss
	date = strings.Replace(date, ":", "-", 2)
	return date
}