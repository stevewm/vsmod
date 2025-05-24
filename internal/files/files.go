package files

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func DownloadFile(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL %s: %v", URL, err)
	}

	defer func(b io.ReadCloser) {
		if b == nil {
			return
		}
		err = errors.Join(err, b.Close())
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var data bytes.Buffer
	_, err = io.Copy(&data, resp.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func WriteFile(fs afero.Fs, fileName string, data []byte) error {
	dir := filepath.Dir(fileName)
	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %s: %w", dir, err)
	}

	file, err := fs.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fileName, err)
	}

	defer func(b io.Closer) {
		if b == nil {
			return
		}
		err = errors.Join(err, b.Close())
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", fileName, err)
	}

	return nil
}
