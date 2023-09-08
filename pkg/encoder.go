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

/* EncodeToDefaultOutputFile():
* Wrapper for Encode() so an output file name doesn't need to be specified.
* Creates an output file name based on the input file name.
 */
func (e *Encoder) EncodeToDefaultOutputFile(inputFileName string) {

	extension := path.Ext(inputFileName)
	nameWithoutExtension := inputFileName[:len(inputFileName)-len(extension)]
	outputFileName := nameWithoutExtension + ".huff"

	Encode(inputFileName, outputFileName, e)
}

/* Encode(): Encodes a text file to a .huff file. */
func Encode(inputFileName, outputFileName string, e *Encoder) {
	//Open input file
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
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
	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()

	//Create Reader & Writer
	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)
	bitWriter := NewBitWriter(outputFile)

	//Generate Frequency Map
	println("Building frequency map...")
	frequencyMap := makeFrequencyMap(inputFileName)

	//Stop encoding if input file is empty
	if len(frequencyMap) == 0 {
		println("Input is empty.")
		println("Encoding complete.")
		return
	}
	//Generate Huffman Tree using frequency map
	println("Generating Huffman tree...")
	huffmanTree := HuffTree{}
	huffmanTree.MakeHuffmanTree(frequencyMap)

	//Generate Code Map from Huffman Tree
	codeMap := huffmanTree.CodeMap()

	//Set instance variables to generated maps
	e.codeMap = codeMap
	e.frequencyMap = frequencyMap

	//Write code table to output file as first line
	println("Writing to file...")
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
			writeEncodedRune(code, bitWriter)
			bodyLength += len(code)
		}
	}
	// Write pseudo-EOF to end of file
	code := e.codeMap[pseudoEOF]
	writeEncodedRune(code, bitWriter)
	bodyLength += len(code)

	// Flush the writers
	bitWriter.Flush()
	writer.Flush()

	println("Compression complete.")
}

/* GetCompressionRatio(): Compares the sizes of the original and the compressed
* files, returns the ratio as a float64.
 */
func GetCompressionRatio(originalFileName, compressedFileName string) float64 {
	oFile, err := os.Open(originalFileName)
	if err != nil {
		log.Fatal(err)
	}
	oFileInfo, err := oFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	cFile, err := os.Open(compressedFileName)
	if err != nil {
		log.Fatal(err)
	}
	cFileInfo, err := cFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	oFileSize := oFileInfo.Size()
	cFileSize := cFileInfo.Size()

	compressionRatio := float64(oFileSize) / float64(cFileSize)

	return compressionRatio
}

/* WriteEncodedRune(): Writes the binary encoding of each rune to the compressed file using
* the BitWriter.
 */
func writeEncodedRune(code string, bitWriter *BitWriter) {
	for i := 0; i < len(code); i++ {
		if code[i] == '1' {
			bitWriter.WriteBit(true)
		} else {
			bitWriter.WriteBit(false)
		}
	}
}

/* DecodeToDefaultOutputFile():
* Wrapper function for Decode(). Allows for decoding without specifying an
* output file. Creates an output file based on the name for the input file.
 */
func (e *Encoder) DecodeToDefaultOutputFile(inputFileName string) {

	extension := path.Ext(inputFileName)
	if extension != ".huff" {
		log.Fatal("Can only decompress .huff files")
	}
	nameWithoutExtension := inputFileName[:len(inputFileName)-len(extension)]
	outputFileName := nameWithoutExtension + "_decoded.txt"

	Decode(inputFileName, outputFileName, e)

}

/* Decode(): Takes an encoded .huff file, decodes and writes decoded text to
* an output file.
 */
func Decode(inputFileName, outputFileName string, e *Encoder) {
	//	Open encoded file
	encodedFile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	//	Close inputFile on exit & check for its returned error
	defer func() {
		if err := encodedFile.Close(); err != nil {
			panic(err)
		}
	}()

	//	Open output file
	decodedFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	//	Close outputFile on exit
	defer func() {
		if err := decodedFile.Close(); err != nil {
			panic(err)
		}
	}()

	//	Create Reader & Writer
	reader := bufio.NewReader(encodedFile)
	bitReader := NewBitReader(reader)
	writer := bufio.NewWriter(decodedFile)

	//	Generate Code Map from printed code table on encoded file
	println("Generating code map...")
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

	//	Exit if code map (and thus the encoded file) is empty
	if len(e.codeMap) == 0 {
		println("Compressed file is empty.")
		println("Decompression complete.")
		return
	}

	//	Create Reversed Code Map for reverse lookups
	e.reverseCodeMap = reverseMap(e.codeMap)

	//	Reset file cursor to beginning of file
	_, err = encodedFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	//	Skip code table printed at top of encoded file
	reader.ReadString('\n')

	//	Read encoded file one bit at a time until a sequence of bits matches a code
	//	in the code table. Once it does, write the corresponding rune to the
	//	decoded .txt file and start again. This is done until EOF is reached.
	println("Decoding file...")
	var code string
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
	println("Decoding complete.")
}

/* reverseMap(): Takes a map, returns the same map but in reverse. */
func reverseMap(m map[int]string) map[string]int {
	n := make(map[string]int, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
