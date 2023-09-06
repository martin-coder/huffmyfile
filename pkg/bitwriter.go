/*
Copyright © 2023 Martin Coder <martincoder1@gmail.com>

Use of this source code is governed by an MIT-style
license that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package huffmyfile

import (
	"io"
)

type BitWriter struct {
	writer io.Writer // The underlying writer
	buffer byte      // Buffer to accumulate bits
	offset uint8     // Current bit offset within the buffer
}

func NewBitWriter(writer io.Writer) *BitWriter {
	return &BitWriter{
		writer: writer,
	}
}

/* WriteBit(): Sets bits into a buffer byte one at a time. When the byte is full (8 bits
*	have been entered), the accumulated byte is written to the underlying writer.
 */
func (bw *BitWriter) WriteBit(bit bool) error {
	if bit {
		bw.buffer |= 1 << (7 - bw.offset) // Set the bit in the buffer
	}

	bw.offset++

	if bw.offset == 8 {
		// Write the accumulated byte to the underlying writer
		_, err := bw.writer.Write([]byte{bw.buffer})
		if err != nil {
			return err
		}

		bw.buffer = 0
		bw.offset = 0
	}

	return nil
}

/* Flush(): Writes whatever bits are left to the underlying writer.
 */
func (bw *BitWriter) Flush() error {
	if bw.offset > 0 {
		// Write the remaining bits to the underlying writer
		_, err := bw.writer.Write([]byte{bw.buffer})
		if err != nil {
			return err
		}

		bw.buffer = 0
		bw.offset = 0
	}

	return nil
}
