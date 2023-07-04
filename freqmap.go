package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

//Tracks the frequency of each character in the input text
type FreqMap struct {
	fmap map[int]int
}

func makeFreqMap(infile string) map[int]int {
	m := make(map[int]int)

	//Open input file
	inf, err := os.Open(infile)
	check(err)
	//close inf on exit & check for its returned error
	defer func() {
		if err := inf.Close(); err != nil {
			panic(err)
		}
	}()

	//Create Reader
	r := bufio.NewReader(inf)
	//Read through input file one rune at a time
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			//Adds 1 to value of c in map
			m[int(c)] += 1
		}
	}

	return m
}
