package address

import (
	"encoding/csv"
	"os"
)

type Region struct {
	SubRegion
	SubRegions []SubRegion
}

type SubRegion struct {
	ID           string
	Latitude     string
	Longitude    string
	JapaneseName string
	EnglishName  string
	GermanName   string
	FrenchName   string
	SpanishName  string
	ItalianName  string
	DutchName    string
}

func loadCSVFile() []Region {
	f, err := os.Open("regions.csv")
	checkError(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	checkError(err)

	var region []Region
	currentCountryIndex := -1
	for i, record := range records {
		if i == 0 {
			// Headers
			continue
		}

		// We are guaranteed to have the coordinates of the country be (0, 0)
		if record[3] == "0" {
			currentCountryIndex++
			region = append(region, Region{
				SubRegion: SubRegion{
					ID:           record[0],
					Latitude:     record[3],
					Longitude:    record[4],
					JapaneseName: record[5],
					EnglishName:  record[6],
					GermanName:   record[7],
					FrenchName:   record[8],
					SpanishName:  record[9],
					ItalianName:  record[10],
					DutchName:    record[11],
				},
				SubRegions: nil,
			})
		} else {
			region[currentCountryIndex].SubRegions = append(region[currentCountryIndex].SubRegions, SubRegion{
				ID:           record[2],
				Latitude:     record[3],
				Longitude:    record[4],
				JapaneseName: record[5],
				EnglishName:  record[6],
				GermanName:   record[7],
				FrenchName:   record[8],
				SpanishName:  record[9],
				ItalianName:  record[10],
				DutchName:    record[11],
			})
		}
	}

	return region
}
