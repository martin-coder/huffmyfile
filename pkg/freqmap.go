/*
Copyright © 2023 Martin Coder <martincoder1@gmail.com>

Use of this source code is governed by an MIT-style
license that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

/*
* Tracks the frequency of each character in the input text
 */

package huffmyfile

import (
	"bufio"
	"io"
	"log"
	"os"
)

const pseudoEOF = int(^uint(0) >> 1) // MaxInt

func makeFrequencyMap(infile string) map[int]int {
	m := make(map[int]int)

	//Open input file
	inf, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
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

	//Add a pseudo-EOF character to aid in decompression (MaxInt)
	m[pseudoEOF] = 1

	return m
}
