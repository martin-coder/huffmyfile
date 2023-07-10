package main

func main() {
	infile := "mobydick.txt"
	outfile := "mobydick_huffed.txt"
	dfile := "mobydick_unhuffed.txt"

	//Create Encoder and use to write encoded output file
	e := Encoder{}
	e.Encode(infile, outfile)

	//Decode outfile and write results to new decoded file
	e.Decode(outfile, dfile)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
