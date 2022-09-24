package logger_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var fakeErrorViper *viper.Viper
var fakeOKViper *viper.Viper

func TestMain(m *testing.M) {
	setUp()
	out := m.Run()
	cleanUp()
	os.Exit(out)
}

func setUp() {
	fakeErrorViper = viper.New()
	fakeErrorViper.SetConfigType("yaml")
	var yamlExample = []byte(`
log: /fake/path/
`)
	fakeErrorViper.ReadConfig(bytes.NewBuffer(yamlExample))

	fakeOKViper = viper.New()
	fakeOKViper.SetConfigType("yaml")
	var yamlOKExample = []byte(`
log: /tmp/
`)
	fakeOKViper.ReadConfig(bytes.NewBuffer(yamlOKExample))
}

func cleanUp() {
	os.Remove("/tmp/app-log")
}
