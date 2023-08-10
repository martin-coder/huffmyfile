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
