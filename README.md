# HuffMyFile

A CLI tool written in Go to compress and decompress text files!

## Installation

### Prerequisites

* Ensure you have Go installed. If not, you can download and install it from the official [Go website](https://golang.org/).

### Install via `go install`
```
$ go install github.com/martin-coder/huffmyfile
```


## Getting Started

### Usage
```
huffmyfile [OPTIONS] [FILE]
```

### Compress a .txt file
```
$ huffmyfile huff [FILE]
```

### Decompress a .huff file
```
$ huffmyfile unhuff [FILE]
```

## Description

HuffMyFile is a command-line tool written in Go that enables you to losslessly compress and decompress text files using Huffman coding. This tool is designed to reduce the size of text files by efficiently encoding characters based on their frequency in the input text.

### What is Huffman coding?
Huffman coding is a popular algorithm used in data compression to encode data, particularly text, in a way that reduces its size while ensuring that no information is lost. It's a variable-length prefix coding technique developed by David A. Huffman in 1952 while he was a Ph.D. student at MIT. The primary idea behind Huffman coding is to assign shorter codes to more frequently occurring symbols and longer codes to less frequent symbols, resulting in an efficient representation of the data.

### Basic Concepts
Huffman coding works by constructing a binary tree called a Huffman tree. This tree is built in a way that ensures the most frequent characters are closer to the root of the tree, allowing for shorter codes. The Huffman tree consists of nodes, each containing a character (or a group of characters) and its associated frequency (how often it appears in the data).

The algorithm involves the following steps:

1. **Character Frequency Analysis:** Count the frequency of each character in the data to be compressed. This is typically done using techniques like hash tables or arrays.
2. **Priority Queue (Min Heap):** Create a priority queue (also known as a min heap) containing the characters and their frequencies. Each character is considered as a single-node tree initially.
3. **Building the Huffman Tree:** While there's more than one node in the priority queue, repeatedly remove the two nodes with the lowest frequencies and create a new internal node (non-leaf node) with their combined frequency. Insert this new node back into the priority queue.
4. **Assigning Codes:** Traverse the Huffman tree from the root to each leaf, assigning a '0' for a left branch and a '1' for a right branch. The path from the root to a leaf node represents the binary code for that character.

### How Huffman Coding Works
Let's take an example to illustrate how Huffman coding works:

Suppose we have the following text: "ABRACADABRA".

1. Calculate the frequency of each character:
`A: 5 | B: 2 | R: 2 | C: 1 | D: 1`

2. Build the priority queue:
```
C:1
D:1
B:2
R:2
A:5
```
3. Build the Huffman tree:
```
      11
     /  \
    A:5  6
        / \
      R:2  4
          / \
        B:2  2         
            / \
          C:1 D:1
```
4. Assign codes by traversing the tree, adding a 0 for each left, and a 1 for each right:
```
A: 0
R: 10
B: 110
C: 1110
D: 1111
```
The compressed representation of "ABRACADABRA" using Huffman coding is:
`01101001110011110110100`

## Limitations

huffmyfile is designed for text files and might not work optimally for binary files.
Large text files might consume substantial memory during compression and decompression.

## Authors

### Contributors names and contact info

**Martin Coder** - [linkedin.com/in/martin-coder](https://www.linkedin.com/in/martin-coder) - [martincoder1@gmail.com](mailto:martincoder1@gmail.com)

## Version History

* 0.1
    * Initial Release
 
## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details


