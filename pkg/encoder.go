/*
Copyright © 2023 Martin Coder <martincoder1@gmail.com>

Use of this source code is governed by an MIT-style
license that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package huffmyfile

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type Encoder struct {
	frequencyMap   map[int]int
	codeMap        map[int]string
	reverseCodeMap map[string]int
}

func (e *Encoder) EncodeToDefaultOutputFile(inputFileName string) {

	extension := path.Ext(inputFileName)
	nameWithoutExtension := inputFileName[:len(inputFileName)-len(extension)]
	outputFileName := nameWithoutExtension + ".huff"

	Encode(inputFileName, outputFileName, e)

}

func Encode(inputFileName, outputFileName string, e *Encoder) {
	//Open input file
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	//close inputFile on exit & check for its returned error
	defer func() {
		if err := inputFile.Close(); err != nil {
			panic(err)
		}
	}()

	//Open output file
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	//close outputFile on exit
	defer outputFile.Close()

	//Create Reader & Writer
	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)
	bitWriter := NewBitWriter(outputFile)

	//Generate Frequency Map
	frequencyMap := makeFrequencyMap(inputFileName)

	//Generate Huffman Tree using frequency map
	huffmanTree := HuffTree{}
	huffmanTree.MakeHuffmanTree(frequencyMap)

	//Generate Code Map from Huffman Tree
	codeMap := huffmanTree.CodeMap()

	//Set instance variables to generated maps
	e.codeMap = codeMap
	e.frequencyMap = frequencyMap

	//Write code table to output file as first line
	for k, v := range e.codeMap {
		_, err := writer.WriteString(fmt.Sprint(k) + " " + v + " ")
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()

	// Track number of bits written to file (not including code table)
	bodyLength := 0

	// Write encoded body of file to output file
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			code := e.codeMap[int(c)]
			WriteEncodedRune(code, bitWriter)
			bodyLength += len(code)
		}
	}
	// Write pseudo-EOF to end of file
	code := e.codeMap[pseudoEOF]
	WriteEncodedRune(code, bitWriter)
	bodyLength += len(code)

	// Fill remaining zeroes in byte
	for i := 0; i < 8-bodyLength%8; i++ {
		bitWriter.WriteBit(false)
	}

	// Write zero byte to end of file to prevent io.UnexpectedEOF error
	// writer.WriteByte(0)
	// writer.WriteByte(0)

	writer.Flush()
}

func WriteEncodedRune(code string, bitWriter *BitWriter) {
	for i := 0; i < len(code); i++ {
		if code[i] == '1' {
			bitWriter.WriteBit(true)
		} else {
			bitWriter.WriteBit(false)
		}
	}
}

func (e *Encoder) DecodeToDefaultOutputFile(inputFileName string) {

	extension := path.Ext(inputFileName)
	if extension != ".huff" {
		log.Fatal("Can only decode .huff files")
	}
	nameWithoutExtension := inputFileName[:len(inputFileName)-len(extension)]
	outputFileName := nameWithoutExtension + "_decoded.txt"

	Decode(inputFileName, outputFileName, e)

}

func Decode(inputFileName, outputFileName string, e *Encoder) {
	//Open encoded file
	encodedFile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	//close inputFile on exit & check for its returned error
	defer func() {
		if err := encodedFile.Close(); err != nil {
			panic(err)
		}
	}()

	//Open output file
	decodedFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	//close outputFile on exit
	defer decodedFile.Close()

	//Create Reader & Writer
	reader := bufio.NewReader(encodedFile)
	bitReader := NewBitReader(reader)
	writer := bufio.NewWriter(decodedFile)

	//Generate Code Map from printed code table on encoded file
	e.codeMap = make(map[int]string)
	scanner := bufio.NewScanner(encodedFile)
	scanner.Scan()
	frequencyMapString := scanner.Text()
	words := strings.Fields(frequencyMapString)

	var last string
	for i, word := range words {
		var k, v string
		if i%2 == 1 {
			k = last
			v = word
			intK, err := strconv.Atoi(k)
			if err != nil {
				log.Fatal(err)
			}
			e.codeMap[intK] = v
		}
		last = word
	}

	fmt.Println("\nNew Code Map: ", e.codeMap)

	//Create Reversed Code Map for reverse lookups
	e.reverseCodeMap = reverseMap(e.codeMap)

	/*
	*	Print encoded body to file:
	*	Read each bit
	*	Append it to a string
	*	Check if it matches a key in the codemap
	*	If not, start again with the next bit
	*	If it does, write that character to the decoded file, reset the string, then start again
	 */

	var code string

	//Reset file cursor to beginning of file
	_, err = encodedFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	//	Skip code table printed at top of encoded file
	reader.ReadString('\n') //Skip code table printed at top of encoded file
	for {

		if b, err := bitReader.ReadBit(); err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		} else {
			if b {
				code = code + "1"
			} else {
				code = code + "0"
			}
		}
		if asciiVal, exists := e.reverseCodeMap[code]; exists {
			if asciiVal == pseudoEOF {
				break
			}
			writer.WriteString(string(rune(asciiVal)))
			code = ""
			writer.Flush()
		}

	}

}

func reverseMap(m map[int]string) map[string]int {
	n := make(map[string]int, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
