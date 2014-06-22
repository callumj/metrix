package resource_bundle

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
)

type CachedFile struct {
	Data        []byte
	ContentType string
	Hash        string
}

var (
	ErrNotExist = errors.New("Specified file does not exist")
	ErrLoading  = errors.New("Unable to load file at this time")
	ErrReading  = errors.New("Unable to read file at this time")
)

var CachedResources map[string]CachedFile

func FetchFile(key string) (CachedFile, error) {
	if CachedResources == nil {
		CachedResources = make(map[string]CachedFile)
	}
	if val, found := CachedResources[key]; found {
		return val, nil
	}

	return fallbackFetchFileFromDisk(key)
}

func fallbackFetchFileFromDisk(key string) (CachedFile, error) {
	if productionMode {
		return CachedFile{}, ErrNotExist
	}

	finalPath := getAssetPath(key)
	if len(finalPath) == 0 {
		return CachedFile{}, ErrNotExist
	}

	buf := bytes.NewBuffer(nil)
	f, err := os.Open(finalPath)
	if err != nil {
		return CachedFile{}, ErrLoading
	}
	defer f.Close()

	written, err := io.Copy(buf, f)
	if written == 0 || err != nil {
		return CachedFile{}, ErrReading
	}

	cType := getMimeType(finalPath)

	resource := newCachedFile(cType, buf.Bytes())

	return resource, nil
}

func getAssetPath(path string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	join := fmt.Sprintf("%v/assets/%v", wd, path)

	fullPath, err := filepath.Abs(join)
	if err != nil {
		return ""
	}

	_, err = os.Stat(fullPath)
	if err != nil {
		return ""
	}

	return fullPath
}

func getMimeType(file string) string {
	fileExt := filepath.Ext(file)
	return mime.TypeByExtension(fileExt)
}

func newCachedFile(cType string, data []byte) CachedFile {
	h := md5.New()
	h.Write(data)
	hashResult := hex.EncodeToString(h.Sum(nil))
	return CachedFile{
		Data:        data,
		ContentType: cType,
		Hash:        hashResult,
	}
}
