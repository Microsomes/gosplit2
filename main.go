package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v3"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

type FileSplit struct {
	FileName       string // name of the file to be split
	SplitBy        uint64 //how many bytes to split by
	ResultsDirName string // directroy name to be created
}

type Part struct {
	Hash     string
	PrevHash string
	PartNo   uint
	PartData string
}

type Block struct {
	Connection *badger.DB
	Header     FileSplit
	Parts      []Part
	PrevHash   string
	Time       uint64
	Hash       string
}

func (block *Block) HashIt() {
	hash := sha256.New()

	for _, x := range block.Parts {
		hash.Write([]byte(x.Hash))
	}

	block.Hash = hex.EncodeToString(hash.Sum(nil))

}

func (part *Part) HashIt() {
	hash := sha256.New()
	hash.Write(append([]byte(part.PartData), byte(part.PartNo)))
	part.Hash = hex.EncodeToString((hash.Sum(nil)))
}

func (t FileSplit) Split() *Block {
	var prevHash string

	fmt.Println("splitting")

	//read entire file in bytes
	by, err := ioutil.ReadFile(t.FileName)
	HandleError(err)

	//convert file to hex string
	toS := hex.EncodeToString(by)
	totalSize := len(toS)
	//calc total parts
	totalParts := totalSize / int(t.SplitBy)

	time := time.Now().Unix()

	var allParts []Part

	for i := 0; i < totalParts; i++ {

		start := i * int(t.SplitBy)
		end := start + int(t.SplitBy)

		part := Part{
			PartNo:   uint(i),
			PartData: "",
		}

		if i >= totalParts-1 {
			fmt.Println("last")
			part.PartData = toS[start:]
		} else {
			part.PartData = toS[start:end]
		}

		part.HashIt()

		if prevHash != "" {
			part.PrevHash = prevHash
			prevHash = part.Hash

		} else {
			prevHash = part.Hash
		}

		allParts = append(allParts, part)

	}
	block := Block{
		Header: t,
		Parts:  allParts,
	}
	block.HashIt()
	fmt.Println(block.Hash)

}

func (t FileSplit) SplitAndSaveDebug() *Block {

	b, err := json.Marshal(block)

	f, err := os.Create(fmt.Sprintf("%s/MANIFEST.json", generatedDirName))
	HandleError(err)
	f.Write(b)

	return &block
}

func main() {

	fs := FileSplit{
		FileName:       "vid.mp4",
		SplitBy:        1000000, //237542379
		ResultsDirName: "vid6",
	}
	fs.SplitAndSaveDebug()

	return

	//recover two halfs

	parts := []string{
		"half1.txt",
		"half2.txt",
	}

	var grabBoth string = ""

	for _, e := range parts {
		by, err := ioutil.ReadFile(e)
		HandleError(err)
		inString := string(by)
		grabBoth = grabBoth + inString
	}

	by, err := hex.DecodeString(grabBoth)
	HandleError(err)

	f, err := os.Create("split2.mp4")
	HandleError(err)

	f.Write(by)

	return
	fmt.Println("hello")

	bytes, err := ioutil.ReadFile("vid.mp4")
	HandleError(err)

	str := hex.EncodeToString(bytes)

	fmt.Println(len(str))

	firstHalf := str[0 : len(str)/2]

	f, err = os.Create("half1.txt")
	HandleError(err)

	f.WriteString(firstHalf)

	secondHalf := str[len(str)/2:]

	f, err = os.Create("half2.txt")
	HandleError(err)

	f.WriteString(secondHalf)

	// f, err := os.Create("dump.txt")
	// HandleError(err)

	// f.WriteString(str)

	// bytes, err := ioutil.ReadFile("dump.txt")
	// HandleError(err)

	// b, err := hex.DecodeString(string(bytes))

	// f, err := os.Create("new.mp4")
	// HandleError(err)
	// f.Write(b)

}
