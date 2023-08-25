package address

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/wii-tools/lzx/lz10"
	"log"
	"os"
	"sync"
)

var (
	key      = []byte{0xff, 0x4c, 0x1a, 0xe3, 0xd4, 0xff, 0xd2, 0x36, 0x71, 0x2e, 0x25, 0x8a, 0x1f, 0x0b, 0x91, 0xe7, 0x2c, 0x91, 0x25, 0xb0, 0xdf, 0x94, 0xc1, 0x69, 0x1b, 0xce, 0xf1, 0x30, 0x11, 0xf1, 0x6c, 0x0f}
	someData = []byte{0x8D, 0x2D, 0x7D, 0x86, 0x76, 0xA6, 0x30, 0xA8, 0x29, 0x72, 0xAB, 0x97, 0x35, 0xE1, 0xA5, 0xCE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("CFH Address file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func MakeAddresses() {
	regions := loadCSVFile()

	wg := sync.WaitGroup{}
	wg.Add(len(regions))
	for _, region := range regions {
		go func(region Region) {
			defer wg.Done()
			regionData := writeRegionData(region)

			arcData := rewriteArc(regionData)

			buffer := new(bytes.Buffer)
			buffer.Write(someData)

			compressed, err := lz10.Compress(arcData)
			checkError(err)

			buffer.Write(compressed)

			for buffer.Len()%aes.BlockSize != 0 {
				buffer.WriteByte(aes.BlockSize)
			}

			block, err := aes.NewCipher(key)
			checkError(err)

			encrypted := make([]byte, buffer.Len())
			mode := cipher.NewCBCEncrypter(block, make([]byte, 16))
			mode.CryptBlocks(encrypted, buffer.Bytes())

			var finalBuffer []byte
			finalBuffer = make([]byte, 64)
			finalBuffer = append(finalBuffer, SignFile(encrypted)...)
			finalBuffer = append(finalBuffer, encrypted...)

			err = os.WriteFile(fmt.Sprintf("%s.alas", ZFill(region.ID, 3)), finalBuffer, 0666)
			checkError(err)
		}(region)
	}
	wg.Wait()
}

func ZFill(str string, size int) string {
	temp := ""

	for i := 0; i < size-len(str); i++ {
		temp += "0"
	}

	return temp + str
}

func SignFile(contents []byte) []byte {
	rsaData, err := os.ReadFile("Private.pem")
	checkError(err)

	rsaBlock, _ := pem.Decode(rsaData)

	parsedKey, err := x509.ParsePKCS1PrivateKey(rsaBlock.Bytes)
	checkError(err)

	// Hash our data then sign
	hash := sha1.New()
	_, err = hash.Write(contents)
	checkError(err)

	contentsHashSum := hash.Sum(nil)

	reader := rand.Reader
	signature, err := rsa.SignPKCS1v15(reader, parsedKey, crypto.SHA1, contentsHashSum)
	checkError(err)

	return signature
}
