package commands

import (
	"errors"
	"html/template"
	"io/ioutil"
	"os"
)

var (
	TestTemplate         = template.Must(template.New("txt").Parse(`{{.Field}}bar`))
	TestTemplateData     = struct{ Field string }{"foo"}
	TestTemplateExpected = "foobar"
	TestTemplateDataKey  = "tplData"
)

type BrokenReader struct{}

func (r *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("reader is broken")
}

func prepTempFile() string {
	var f *os.File
	f, _ = ioutil.TempFile("", "TestFileWriter-")
	fn := f.Name()
	_ = f.Close()
	_ = os.Remove(fn)
	return fn
}

func cleanFile(fn string) {
	_ = os.Remove(fn)
}
