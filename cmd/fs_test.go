package cmd

import (
	"os"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

type mockFileSystem struct {
	ReturnExists    func(string) bool
	ReturnMkdirP    func(string) error
	ReturnCreate    func(string) (File, error)
	ReturnWriteFile func(string, []byte) error
}

func (m *mockFileSystem) Exists(path string) bool {
	return m.ReturnExists(path)
}

func (m *mockFileSystem) MkdirP(path string) error {
	return m.ReturnMkdirP(path)
}

func (m *mockFileSystem) Create(name string) (File, error) {
	return m.ReturnCreate(name)
}

func (m *mockFileSystem) WriteFile(path string, content []byte) error {
	return m.ReturnWriteFile(path, content)
}

type mockFile struct {
	File
	ReturnWrite func([]byte) (int, error)
}

func (m *mockFile) Write(b []byte) (int, error) {
	return m.ReturnWrite(b)
}

func TestFs(t *testing.T) {
	mockey.PatchConvey("fs", t, func() {
		mockey.Mock(os.Stat).Return(nil, nil).Build()
		mockey.Mock(os.Create).Return(nil, nil).Build()
		mockey.Mock(os.WriteFile).Return(nil).Build()
		sys := &osFileSystem{}
		assert.Equal(t, true, sys.Exists(""))
		_, err := sys.Create("")
		assert.Nil(t, err)
		assert.Nil(t, sys.WriteFile("", nil))
	})
}
