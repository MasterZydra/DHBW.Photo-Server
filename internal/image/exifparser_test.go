package image

import (
	"os"
	"testing"
)

func TestReadDate(t *testing.T) {
	filename := "./Testbild.jpg"
	file, err := os.Open(filename)
	if err != nil {
		t.Errorf("Error reading image: %v", err)
	}

	exifData, err := parseRawExifDataFromFile(file)
	if err != nil {
		t.Error(err)
	}

	date := getDateFromData(exifData)

	wantedDate := "2020-08-20 14:23:45"

	if string(date) != wantedDate {
		t.Errorf("Error parsing date from Exif-data. Wanted: %v, Got: %v", wantedDate, string(date))
	}
}

func TestGetDate(t *testing.T) {
	testdata := [][]byte{[]byte("a8Q2o81O7GwBD0GMhBTNxWN9MPp7f4wogshEqnEpmOHTZm0n2020:06:08 03:49:39q19AVUR3TPftEXcu1G903NYpLUatW9VvPvakgsFvBMBTuVhQ1iEFar1C116sIbEMJiOqyS6rXQMFPcwiRogqSxkPmddFIk0fwL0Wx9MLEXrDhB6F"),
						[]byte("J8i9H5LBbg1jDezNjVIA2lbQVfF05FEg0Pk6U5BCc0zhe0nN41MlKG1yNNZRkUlZyP0AKBmz4Ys0aTyPoZY0WGvac732z6z91s8E9StlajfzZBI2l8EHfZNcEIi8MCLIVUh9UU2018:03:29 13:19:58dMPJTum8IbSYOn4TfX8t2oMvKd"),
						[]byte("Kpl4IG2S95YDwSV58vlOlquh6eo6CongO0BOwp0xM9AzvMCGRtMovBr4BIPQOgYATSopAyrH1Xk1FNi8irdIxfHLT172UHBqI1pRku52018:12:28 05:06:16eWwEWOnpoiSuoRmgM2RWsgmqw2RaaBYJemCwuzWUni8uKj9wIphYoTMTL"),
						[]byte("6T3VG4nXGETcRqWYiY5fjeUS2018:05:21 04:42:22euQzdQm55bjigHY9643PCj456vZII23QQ9bA3ZD26SCx8MiJQhugbIltLdpU2tyIr9mOIePMuMeoGbTG3UBZlT90ti04iY0UesiSIu7gO5ySCDtdmBtxqa3ilFlpNbC6Sg3HGZzA"),
						[]byte("TSIy2020:07:10 04:33:56wbJMNmvTX4xBIGvfQPMEpGiRQ0epQqGj1aDTe7JkMZzyyjDQU8HjjFnXqn0TSjhZf8ki9LmA9wWO3exRcxPYYQb3cVEv2oQqIygmzp4d5DMyRmTStQnQgWL4eokXszXoSMG8WZvGcz0rRvAXI8EVmen2pZUl"),
						[]byte("DDc95zi01JhpEHfShmLShfeWrS4TiS4616MU22wh8s4I6XpfI3arPpD2h4077FHGKUTl8ALkWhDWky3O1yUumRMgLXMxIzUcykaCOCQ7N2XgqBxJ1lS6EiieBcgn3gzezygxB62019:10:12 23:46:53xT6ejn3aqr5VKnlMXY4QL526lU"),
						[]byte("EdTg11lNlbrhB5zH7fn31U75coaWHl7UsUNjEt9OuEHB4kuway23YdNOl4y1aUPhwNEF0YaszOG6moZY7gj2019:07:06 20:32:16eBlzULh8yedPmWKqQavmqEH2yyE0CfPpj6wMghQQOfwSHJHI0eYt1hkWiTOiO0SoGb8FBX2YwAqA9"),
						[]byte("Jl38RVr1WiJGndibH2020:01:06 01:10:50QUs5STpE3QdJIVn7bKIiDh8YKLvN7f8rXYuwAncJqDiX2TIG89DdxvST0WYZBcnAlHuaIEDBidWUjm1LnTGSg8cQAtyjka6xMsIyICDWnzWoROogashKMKPY1PtgJGleA0799MdpyQLkcUX"),
						[]byte("9zzVW9vmA2020:05:21 01:59:32HuH213ynSRkG92ikUOieWE7Xa4c7C6cJRV0yOZz0ft03gKzhkDr8v7W6Sp6YbtTjqO2W2vUb3yjVKKgGogNOeabNlLOciuuziUMhOKakWkaQXF89YN01IXQxLf0ieR3HhoK6PkdVOVPz9cbC8RpF1TE"),
						[]byte("37ttP8BLcU1ABwwxBWBmrafBfWENHDDRHfjk7v00eaK7Xwp4d3czzOaHSEO3FUYGWclm33kMXbzlbrFVjCADHNETzZy9GXFCB8lbbRE2020:10:30 09:28:310veoWq7rakwXtJ8YcsuuyzGQeLer2hyO7ehg6EtDKk2gsXwhMjas3JE4Q"),
						[]byte("kICzBtBtEWPtTs6h7ch3Y3b62fNhPC6hd8XJ2019:06:08 07:55:00mnVFr7dGF7scyIXJJnytGh4nWpTkJOqYsFP8RzP2DjrleDywF6oLBSiGSwcuufKkBilmkdu08tKKDxbFrBObukodLtHjTDRt18F9XVthXNVQD8y5LGLZktnvU4sO"),
						[]byte("Hxw34Sc3RSeDs1H00LioURQHWZ2FNf6ebaeKDb6qVZ8LiCPzFjdfg0OYVAutzx1jenyNfvLnSno48QXtKkyEni76rx90mVnYrMHCCQZwiqQ7SNSs8uqQt6bCD10GlwPzWvspPwl7xp2018:09:02 06:23:04RVSghS6Nnpm8TT7n8w7irt")}
	expectedDates := []string {"2020-06-08 03:49:39",
						"2018-03-29 13:19:58",
						"2018-12-28 05:06:16",
						"2018-05-21 04:42:22",
						"2020-07-10 04:33:56",
						"2019-10-12 23:46:53",
						"2019-07-06 20:32:16",
						"2020-01-06 01:10:50",
						"2020-05-21 01:59:32",
						"2020-10-30 09:28:31",
						"2019-06-08 07:55:00",
						"2018-09-02 06:23:04"}

	for index, data := range testdata {
		date := getDateFromData(data)
		if date != expectedDates[index] {
			t.Errorf("Error parsing date from Exif-data. Wanted: %v, Got: %v", expectedDates[index], string(date))
		}
	}
}