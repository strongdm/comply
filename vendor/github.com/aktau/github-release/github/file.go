package github

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

func GetFileSize(f *os.File) (int64, error) {
	/* first try stat */
	off, err := fsizeStat(f)
	if err != nil {
		/* if that fails, try seek */
		return fsizeSeek(f)
	}

	return off, nil
}

func fsizeStat(f *os.File) (int64, error) {
	fi, err := f.Stat()

	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

func fsizeSeek(f *os.File) (int64, error) {
	off, err := f.Seek(0, 2)
	if err != nil {
		return 0, fmt.Errorf("seeking did not work, stdin is not" +
			"supported yet because github doesn't support chunking" +
			"requests (and I haven't implemented detecting stdin and" +
			"buffering yet")
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return 0, fmt.Errorf("could not seek back in the file")
	}
	return off, nil
}

// materializeFile takes a physical file or stream (named pipe, user input,
// ...) and returns an io.Reader and the number of bytes that can be read
// from it.
func materializeFile(f *os.File) (io.Reader, int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}

	// If the file is actually a char device (like user typed input)
	// or a named pipe (like a streamed in file), buffer it up.
	//
	// When uploading a file, you need to either explicitly set the
	// Content-Length header or send a chunked request. Since the
	// github upload server doesn't accept chunked encoding, we have
	// to set the size of the file manually. Since a stream doesn't have a
	// predefined length, it's read entirely into a byte buffer.
	if fi.Mode()&(os.ModeCharDevice|os.ModeNamedPipe) == 1 {
		vprintln("input was a stream, buffering up")

		var buf bytes.Buffer
		n, err := buf.ReadFrom(f)
		if err != nil {
			return nil, 0, errors.New("req: could not buffer up input stream: " + err.Error())
		}
		return &buf, n, err
	}

	// We know the os.File is most likely an actual file now.
	n, err := GetFileSize(f)
	return f, n, err
}
