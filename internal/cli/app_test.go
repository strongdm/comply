package cli

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/strongdm/comply/internal/util"
	"github.com/urfave/cli"
)

type PandocMustExist struct{}

func TestPandocMustExist(t *testing.T) {
	util.ExecuteTests(t, reflect.TypeOf(PandocMustExist{}), beforeEach, nil)
}

func beforeEach() {
	util.MockConfig()
}

func (tg PandocMustExist) WhenBinaryExists(t *testing.T) {
	dockerPullCalled := false

	pandocBinaryMustExist = func(c *cli.Context) (e error, found, goodVersion, pdfLatex bool) {
		return nil, true, true, true
	}

	dockerMustExist = func(c *cli.Context) (e error, inPath bool, isRunning bool) {
		return errors.New("docker doesn't exist"), false, false
	}

	pandocImageExists = func(ctx context.Context) bool {
		return false
	}

	dockerPull = func(c *cli.Context) error {
		dockerPullCalled = true
		return nil
	}

	pandocMustExist(&cli.Context{})

	if dockerPullCalled != false {
		t.Fatal("Docker was pulled")
	}
}

func (tg PandocMustExist) WhenImageExists(t *testing.T) {
	dockerPullCalled := false

	pandocBinaryMustExist = func(c *cli.Context) (e error, found, goodVersion, pdfLatex bool) {
		return errors.New("binary doesn't exist"), false, true, true
	}

	dockerMustExist = func(c *cli.Context) (e error, inPath bool, isRunning bool) {
		return errors.New("docker doesn't exist"), false, false
	}

	pandocImageExists = func(ctx context.Context) bool {
		return true
	}

	dockerPull = func(c *cli.Context) error {
		dockerPullCalled = true
		return nil
	}

	pandocMustExist(&cli.Context{})

	if dockerPullCalled != false {
		t.Fatal("Docker was pulled")
	}
}

func (tg PandocMustExist) WhenBinaryAndImageDontExists(t *testing.T) {
	dockerPullCalled := false

	pandocBinaryMustExist = func(c *cli.Context) (e error, found, goodVersion, pdfLatex bool) {
		return errors.New("binary doesn't exist"), false, false, false
	}

	dockerMustExist = func(c *cli.Context) (e error, inPath bool, isRunning bool) {
		return nil, true, false
	}

	pandocImageExists = func(ctx context.Context) bool {
		return false
	}

	dockerPull = func(c *cli.Context) error {
		dockerPullCalled = true
		return nil
	}

	pandocMustExist(&cli.Context{})

	if dockerPullCalled != true {
		t.Fatal("Docker wasn't pulled")
	}
}

func (tg PandocMustExist) WhenCannotPullPandoc(t *testing.T) {
	dockerPullCalled := false

	pandocBinaryMustExist = func(c *cli.Context) (e error, found, goodVersion, pdfLatex bool) {
		return errors.New("binary doesn't exist"), false, false, false
	}

	dockerMustExist = func(c *cli.Context) (e error, inPath bool, isRunning bool) {
		return errors.New("docker doesn't exist"), false, false
	}

	pandocImageExists = func(ctx context.Context) bool {
		return false
	}

	dockerPull = func(c *cli.Context) error {
		dockerPullCalled = true
		return nil
	}

	pandocMustExist(&cli.Context{})

	if dockerPullCalled != false {
		t.Fatal("Docker was pulled")
	}
}

func (tg PandocMustExist) WhenMustUseLocalPandoc(t *testing.T) {
	os.Setenv("COMPLY_USE_LOCAL_PANDOC", "true")
	dockerPullCalled := false

	pandocBinaryMustExist = func(c *cli.Context) (e error, found, goodVersion, pdfLatex bool) {
		return errors.New("binary doesn't exist"), false, false, false
	}

	dockerMustExist = func(c *cli.Context) (e error, inPath bool, isRunning bool) {
		return errors.New("docker doesn't exist"), false, false
	}

	pandocImageExists = func(ctx context.Context) bool {
		return true
	}

	dockerPull = func(c *cli.Context) error {
		dockerPullCalled = true
		return nil
	}

	pandocMustExist(&cli.Context{})

	if dockerPullCalled != false {
		t.Fatal("Docker was pulled")
	}

	os.Clearenv()
}

func (tg PandocMustExist) WhenPandocDontExistsAndCannotPull(t *testing.T) {
	os.Setenv("COMPLY_USE_LOCAL_PANDOC", "true")
	dockerPullCalled := false

	pandocBinaryMustExist = func(c *cli.Context) (e error, found, goodVersion, pdfLatex bool) {
		return errors.New("binary doesn't exist"), false, false, false
	}

	dockerMustExist = func(c *cli.Context) (e error, inPath bool, isRunning bool) {
		return nil, true, false
	}

	pandocImageExists = func(ctx context.Context) bool {
		return false
	}

	dockerPull = func(c *cli.Context) error {
		dockerPullCalled = true
		return nil
	}

	pandocMustExist(&cli.Context{})

	if dockerPullCalled != false {
		t.Fatal("Docker was pulled")
	}
	os.Clearenv()
}
