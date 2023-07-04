package main

func main() {
	infile := "input.txt"
	outfile := "output.txt"
	dfile := "decoded.txt"

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
