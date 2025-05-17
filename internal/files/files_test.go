package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	type test struct {
		testFile []byte
		dirName  string
		fileName string
	}

	tests := []test{
		{testFile: []byte("test file"), dirName: "/", fileName: "test.txt"},
		{testFile: []byte(""), fileName: "test.txt"},
		{testFile: []byte("test file"), fileName: "test.txt"},
		{testFile: []byte("windows"), fileName: "windows.txt"},
	}

	for _, tc := range tests {
		testFs := afero.NewBasePathFs(afero.NewMemMapFs(), "test")
		err := WriteFile(testFs, tc.fileName, tc.testFile)
		path := tc.dirName + tc.fileName
		assert.NoError(t, err, "Expected no error when writing '%s' to %s", tc.testFile, path)

		readFile, _ := afero.ReadFile(testFs, path)
		assert.Equal(t, tc.testFile, readFile, "expected: '%s', got: '%s'", tc.testFile, readFile)
	}

}
