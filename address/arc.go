package address

import (
	_ "embed"
	"github.com/wii-tools/arclib"
)

//go:embed address.u8
var archive []byte

// rewriteArc rewrites the regionData.js file of a preexisting U8 archive with what we generated.
func rewriteArc(regionData []byte) []byte {
	arc, err := arclib.Load(archive)
	checkError(err)

	file, err := arc.OpenFile("regionData.js")
	checkError(err)

	file.Write(regionData)

	data, err := arc.Save()
	checkError(err)

	return data
}
