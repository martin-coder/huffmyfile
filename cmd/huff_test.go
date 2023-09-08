package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func TestHuff(t *testing.T) {
	testFileName := "testfile.txt"
	compressedTestFileName := "testfile.huff"
	decodedTestFileName := "testfile_decoded.txt"
	testFile, err := os.Create(testFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Test on empty file
	huffCmd := NewHuffCmd(testFileName)
	huffCmd.Execute()

	unhuffCmd := NewUnhuffCmd(compressedTestFileName)
	unhuffCmd.Execute()

	if !deepCompare(testFileName, decodedTestFileName) {
		t.Errorf("Test Case 1 failed. Input file not equal to decoded file.")
	}

	testContent := "ABRACADABRA\nalakazam\n! : åßˆ\n\n"
	testFile.WriteString(testContent)

	huffCmd = NewHuffCmd(testFileName)
	huffCmd.Execute()

	unhuffCmd = NewUnhuffCmd(compressedTestFileName)
	unhuffCmd.Execute()

	if !deepCompare(testFileName, decodedTestFileName) {
		t.Errorf("Test Case 2 failed. Input file not equal to decoded file.")
	}

}

const chunkSize = 64000

func deepCompare(file1, file2 string) bool {
	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
