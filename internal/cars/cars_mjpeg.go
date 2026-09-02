package cars

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type MJPEGReader struct {
	reader *bufio.Reader
}

func NewMJPEGReader(r io.Reader) *MJPEGReader {
	return &MJPEGReader{
		reader: bufio.NewReader(r),
	}
}

func (r *MJPEGReader) ReadFrame() ([]byte, error) {
	for {
		line, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		if line == "--frame" {
			break
		}
	}

	var contentLength int

	for {
		line, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		const prefix = "Content-Length:"

		if strings.HasPrefix(line, prefix) {
			value := strings.TrimSpace(
				strings.TrimPrefix(line, prefix),
			)

			contentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid content length: %w",
					err,
				)
			}
		}
	}

	if contentLength <= 0 {
		return nil, fmt.Errorf("invalid content length: %d", contentLength)
	}

	frame := make([]byte, contentLength)

	if _, err := io.ReadFull(r.reader, frame); err != nil {
		return nil, err
	}

	if _, err := r.reader.ReadString('\n'); err != nil {
		return nil, err
	}

	return bytes.Clone(frame), nil
}
