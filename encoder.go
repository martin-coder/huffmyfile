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
	freqMap    map[int]int
	codeMap    map[int]string
	revCodeMap map[string]int
}

func (e *Encoder) Encode(infile, outfile string) {
	//Open input file
	inf, err := os.Open(infile)
	check(err)
	//close inf on exit & check for its returned error
	defer func() {
		if err := inf.Close(); err != nil {
			panic(err)
		}
	}()

	//Open output file
	outf, err := os.Create(outfile)
	check(err)
	//close outf on exit
	defer outf.Close()

	//Create Reader & Writer
	r := bufio.NewReader(inf)
	w := bufio.NewWriter(outf)

	//Generate Frequency Map
	fmap := makeFreqMap(infile)
	fmt.Println("\nFrequency Map: ", fmap)

	//Generate Huffman Tree using frequency map
	h := huffmantree.HuffTree{}
	h.MakeHuffmanTree(fmap)
	fmt.Println()
	fmt.Print("Tree: ")
	h.Print()
	fmt.Println()

	//Generate Code Map from Huffman Tree
	cmap := h.CodeMap()

	//Set instance variables to generated maps
	e.codeMap = cmap
	e.freqMap = fmap

	//Write code table to output file as first line
	for k, v := range e.codeMap {
		_, err := w.WriteString(fmt.Sprint(k) + " " + v + " ")
		check(err)
	}
	_, err = w.WriteString("\n")
	check(err)
	w.Flush()

	//Write encoded body of file to output file
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			w.WriteString(e.codeMap[int(c)])
		}
	}
	w.Flush()
}

func (e *Encoder) Decode(enfile, dfile string) {
	//Open encoded file
	enf, err := os.Open(enfile)
	check(err)
	//close inf on exit & check for its returned error
	defer func() {
		if err := enf.Close(); err != nil {
			panic(err)
		}
	}()

	//Open output file
	df, err := os.Create(dfile)
	check(err)
	//close outf on exit
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
	e.revCodeMap = reverseMap(e.codeMap)
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
			if asciiVal, exists := e.revCodeMap[sb.String()]; exists {
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
