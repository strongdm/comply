package main

import (
	"fmt"
	"os"
	"time"
)

/* oracle nvl, return first non-empty string */
func nvls(xs ...string) string {
	for _, s := range xs {
		if s != "" {
			return s
		}
	}

	return ""
}

func vprintln(a ...interface{}) (int, error) {
	if VERBOSITY > 0 {
		return fmt.Fprintln(os.Stderr, a...)
	}

	return 0, nil
}

func vprintf(format string, a ...interface{}) (int, error) {
	if VERBOSITY > 0 {
		return fmt.Fprintf(os.Stderr, format, a...)
	}

	return 0, nil
}

// formats time `t` as `fmt` if it is not nil, otherwise returns `def`
func timeFmtOr(t *time.Time, fmt, def string) string {
	if t == nil {
		return def
	}
	return t.Format(fmt)
}

// isCharDevice returns true if f is a character device (panics if f can't
// be stat'ed).
func isCharDevice(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
