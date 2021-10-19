package util

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/strongdm/comply/internal/config"
	"gopkg.in/yaml.v2"
)

type TestFixture func()

func ExecuteTests(t *testing.T, testGroupType reflect.Type, beforeEach TestFixture, afterEach TestFixture) {
	testGroup := reflect.New(testGroupType).Elem().Interface()
	for i := 0; i < testGroupType.NumMethod(); i++ {
		m := testGroupType.Method(i)
		t.Run(m.Name, func(t *testing.T) {
			if beforeEach != nil {
				beforeEach()
			}

			in := []reflect.Value{reflect.ValueOf(testGroup), reflect.ValueOf(t)}
			m.Func.Call(in)

			if afterEach != nil {
				afterEach()
			}
		})
	}
}

func MockConfig() {
	config.Config = func() *config.Project {
		p := config.Project{}
		cfgBytes, _ := ioutil.ReadFile(filepath.Join(GetRootPath(), "comply.yml.example"))
		err := yaml.Unmarshal(cfgBytes, &p)
		if err != nil {
			return nil
		}
		return &p
	}
}

func GetRootPath() string {
	_, fileName, _, _ := runtime.Caller(0)
	fileDir := filepath.Dir(fileName)
	return fmt.Sprintf("%s/../../example", fileDir)
}
