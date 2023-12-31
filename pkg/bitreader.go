/*
Copyright © 2023 Martin Coder <martincoder1@gmail.com>

Use of this source code is governed by an MIT-style
license that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package huffmyfile

import (
	"errors"
	"io"
)

type BitReader struct {
	reader   io.Reader // Underlying reader
	buffer   byte      // Buffer to read individual bits from
	bitCount uint8     // Keeping track of what bits have been read
	err      error
}

func NewBitReader(reader io.Reader) *BitReader {
	return &BitReader{
		reader:   reader,
		buffer:   0,
		bitCount: 0,
		err:      nil,
	}
}

/*
*	ReadBit(): Wrapper function for readBit() to return a bool rather than uint8
 */
func (br *BitReader) ReadBit() (bit bool, err error) {
	if br.err != nil {
		return false, br.err
	}

	b, err := br.readBit()
	if err != nil {
		return false, err
	}

	if b == 1 {
		return true, nil
	} else {
		return false, nil
	}

}

/*	readBit(): Reads the next bit from the buffer. Reads a new byte to the buffer when
*	bitCount reaches 0.
 */
func (br *BitReader) readBit() (bit uint8, err error) {
	if br.err != nil {
		return 0, br.err
	}

	if br.bitCount == 0 {
		br.buffer, err = br.ReadByte()
		if err != nil {
			return 0, err
		}
		br.bitCount = 8
	}

	bit = (br.buffer >> (br.bitCount - 1)) & 1
	br.bitCount--
	return bit, nil
}

/*	ReadByte(): Wrapper function for readByte(). Throws an error if ReadByte() is called
*	while the buffer is still full.
 */
func (br *BitReader) ReadByte() (b byte, err error) {
	if br.err != nil {
		return 0, br.err
	}

	if br.bitCount > 0 {
		return 0, errors.New("cannot read new byte with bits remaining in buffer")
	}

	return br.readByte()
}

/*	readByte(): Uses the underlying reader to read a single byte into a buffer, then
*	returns that byte.
 */
func (br *BitReader) readByte() (b byte, err error) {
	buff := make([]byte, 1)
	_, err = br.reader.Read(buff)
	if err != nil {
		br.err = err
		return 0, err
	}
	return buff[0], nil
}
