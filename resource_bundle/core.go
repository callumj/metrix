package resource_bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

type CachedFile struct {
	Data        []byte
	ContentType string
}

var CachedResources map[string]CachedFile

var productionMode bool = false
var assetReplaceReg = regexp.MustCompile(`^assets/`)

func FetchFilesFromSelf() {
	if CachedResources == nil {
		CachedResources = make(map[string]CachedFile)
	}

	pathToSelf, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Printf("Not checking self for asset zipfile. %v\r\n", err)
		return
	}

	f, err := os.Open(pathToSelf)
	if err != nil {
		log.Printf("Unable open self to look for bundled zip: %v\r\n", err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Printf("Unable query self to look for bundled zip: %v\r\n", err)
		return
	}

	f.Seek(fi.Size()-1024, 0)
	buf := make([]byte, 1024)

	_, err = f.Read(buf)
	contents := string(buf)
	re := regexp.MustCompile("ArchiveLength:(\\d+)$")
	info := re.FindStringSubmatch(contents)

	if len(info) != 2 {
		return
	}

	fullString, sizeStr := info[0], info[1]
	archiveSize, _ := strconv.ParseInt(sizeStr, 10, 64)

	seekLen := fi.Size() - (int64(len(fullString)) + archiveSize)
	curOff, _ := f.Seek(seekLen, 0)
	log.Printf("Seeked to: %v for %v\r\n", curOff, archiveSize)

	zipBuf := bytes.NewBuffer(nil)
	_, err = io.CopyN(zipBuf, f, archiveSize)

	zipBufR := bytes.NewReader(zipBuf.Bytes())

	zipR, err := zip.NewReader(zipBufR, archiveSize)
	if err != nil {
		log.Printf("Unable to load self for bundled zip: %v\r\n", err)
		return
	}

	for _, f := range zipR.File {
		if f.UncompressedSize64 != 0 {
			fName := assetReplaceReg.ReplaceAllString(f.Name, "")
			log.Printf("Preparing to load asset %s\r\n", fName)
			fPntr, err := f.Open()
			if err != nil {
				log.Printf("Failed to open %s: %v\r\n", fName, err)
			}

			readBuf := bytes.NewBuffer(nil)
			io.Copy(readBuf, fPntr)

			cType := getMimeType(fName)
			CachedResources[fName] = CachedFile{
				Data:        readBuf.Bytes(),
				ContentType: cType,
			}
			log.Printf("Loaded asset %v (%v)\r\n", fName, cType)
		}
	}
	productionMode = true
}

func FetchFile(key string) (CachedFile, error) {
	if CachedResources == nil {
		CachedResources = make(map[string]CachedFile)
	}
	if val, found := CachedResources[key]; found {
		return val, nil
	}

	if productionMode {
		return CachedFile{}, errors.New("Production mode is active")
	}

	finalPath := getAssetPath(key)
	if len(finalPath) == 0 {
		return CachedFile{}, errors.New("Cannot find file")
	}

	buf := bytes.NewBuffer(nil)
	f, err := os.Open(finalPath)
	if err != nil {
		return CachedFile{}, errors.New("Cannot load file")
	}
	defer f.Close()

	written, err := io.Copy(buf, f)
	if written == 0 || err != nil {
		return CachedFile{}, errors.New("Cannot read file")
	}

	cType := getMimeType(finalPath)

	CachedResources[key] = CachedFile{
		Data:        buf.Bytes(),
		ContentType: cType,
	}

	return CachedResources[key], nil
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

	return fullPath
}

func getMimeType(file string) string {
	fileExt := filepath.Ext(file)
	return mime.TypeByExtension(fileExt)
}
