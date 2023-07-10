package main

import (
	"HuffmanCoder/huffmantree"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Encoder struct {
	frequencyMap   map[int]int
	codeMap        map[int]string
	reverseCodeMap map[string]int
}

func (e *Encoder) Encode(inputFileName, outputFileName string) {
	//Open input file
	inputFile, err := os.Open(inputFileName)
	check(err)
	//close inputFile on exit & check for its returned error
	defer func() {
		if err := inputFile.Close(); err != nil {
			panic(err)
		}
	}()

	//Open output file
	outputFile, err := os.Create(outputFileName)
	check(err)
	//close outputFile on exit
	defer outputFile.Close()

	//Create Reader & Writer
	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	//Generate Frequency Map
	frequencyMap := makeFrequencyMap(inputFileName)
	fmt.Println("\nFrequency Map: ", frequencyMap)

	//Generate Huffman Tree using frequency map
	huffmanTree := huffmantree.HuffTree{}
	huffmanTree.MakeHuffmanTree(frequencyMap)
	fmt.Println()
	fmt.Print("Tree: ")
	huffmanTree.Print()
	fmt.Println()

	//Generate Code Map from Huffman Tree
	codeMap := huffmanTree.CodeMap()

	//Set instance variables to generated maps
	e.codeMap = codeMap
	e.frequencyMap = frequencyMap

	//Write code table to output file as first line
	for k, v := range e.codeMap {
		_, err := writer.WriteString(fmt.Sprint(k) + " " + v + " ")
		check(err)
	}
	_, err = writer.WriteString("\n")
	check(err)
	writer.Flush()

	//Write encoded body of file to output file
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			writer.WriteString(e.codeMap[int(c)])
		}
	}
	writer.Flush()
}

func (e *Encoder) Decode(enfile, dfile string) {
	//Open encoded file
	enf, err := os.Open(enfile)
	check(err)
	//close inputFile on exit & check for its returned error
	defer func() {
		if err := enf.Close(); err != nil {
			panic(err)
		}
	}()

	//Open output file
	df, err := os.Create(dfile)
	check(err)
	//close outputFile on exit
	defer df.Close()

	//Create Reader & Writer
	r := bufio.NewReader(enf)
	w := bufio.NewWriter(df)

	//Generate Code Map from printed code table on encoded file
	scanner := bufio.NewScanner(enf)
	scanner.Scan()
	fmapstr := scanner.Text()
	words := strings.Fields(fmapstr)

	var last string
	for i, word := range words {
		var k, v string
		if i%2 == 1 {
			k = last
			v = word
			intK, err := strconv.Atoi(k)
			check(err)
			e.codeMap[intK] = v
		}
		last = word
	}

	fmt.Println("\nNew Code Map: ", e.codeMap)

	//Create Reversed Code Map for reverse lookups
	e.reverseCodeMap = reverseMap(e.codeMap)
	fmt.Println("\nRev Code Map: ", reverseMap(e.codeMap))
	fmt.Println()

	sb := strings.Builder{}

	enf.Seek(0, 0)              //Reset file cursor to beginning of file
	_, err = r.ReadString('\n') //Skip code table printed at top of encoded file
	check(err)

	//Print encoded body to file
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			sb.WriteRune(c)
			if asciiVal, exists := e.reverseCodeMap[sb.String()]; exists {
				w.WriteString(string(rune(asciiVal)))
				sb.Reset()
				w.Flush()
			}
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
