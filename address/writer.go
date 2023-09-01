package address

import (
	"bytes"
	"fmt"
)

const (
	ArrayStart   = "var RegionInfo = new Array(\n"
	ArrayContent = "    new Array(%s, %s, %s, new Array(\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"))%s\n"
	ArrayEnd     = ");\n"
)

func writeRegionData(region Region) []byte {
	buffer := new(bytes.Buffer)

	// The region data is a Javascript Array in a file named `regionData.js`.
	buffer.WriteString(ArrayStart)

	// It is possible for a country to not have any subregions. Take for example, any Caribbean Island.
	separator := ","
	if len(region.SubRegions) == 0 {
		separator = ""
	}

	buffer.WriteString(fmt.Sprintf(ArrayContent, "1", region.Latitude, region.Longitude,
		region.JapaneseName, region.EnglishName, region.GermanName,
		region.FrenchName, region.SpanishName, region.ItalianName, region.DutchName, separator))

	for i, subRegion := range region.SubRegions {
		separator := ","
		if i == len(region.SubRegions)-1 {
			separator = ""
		}

		buffer.WriteString(fmt.Sprintf(ArrayContent, subRegion.ID, subRegion.Latitude, subRegion.Longitude,
			subRegion.JapaneseName, subRegion.EnglishName, subRegion.GermanName,
			subRegion.FrenchName, subRegion.SpanishName, subRegion.ItalianName, subRegion.DutchName, separator))
	}

	// And fin.
	buffer.WriteString(ArrayEnd)
	return buffer.Bytes()
}
