// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"fmt"
	"runtime"
)

const AppName = "go-bindata"

var AppVersion = "dev"

func Version() string {
	return fmt.Sprintf("%s %s (Go runtime %s).\nCopyright (c) 2010-2013, Jim Teeuwen.",
		AppName, AppVersion, runtime.Version())
}
