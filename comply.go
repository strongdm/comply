package main

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/strongdm/comply/internal/cli"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	godotenv.Load(fmt.Sprintf("%s/.env", basepath))
	cli.Main()
}
