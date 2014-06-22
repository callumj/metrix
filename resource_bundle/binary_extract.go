package resource_bundle

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var productionMode bool = false
var assetReplaceReg = regexp.MustCompile(`^assets/`)

func FetchFilesFromSelf() {
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

	zipBufR, archiveSize := buildZipArchive(f, fi)
	if zipBufR == nil || archiveSize == 0 {
		return
	}

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

			addResource(fName, readBuf.Bytes())
		}
	}
	productionMode = true
}

func addResource(fName string, data []byte) {
	if CachedResources == nil {
		CachedResources = make(map[string]CachedFile)
	}

	cType := getMimeType(fName)
	CachedResources[fName] = newCachedFile(cType, data)
	log.Printf("Loaded asset %v (%v)\r\n", fName, cType)
}

func buildZipArchive(fPointer *os.File, fInfo os.FileInfo) (*bytes.Reader, int64) {
	info := readMetaFromFile(fPointer, fInfo)
	if len(info) == 0 {
		return nil, 0
	}

	fullString, sizeStr := info[0], info[1]
	archiveSize, _ := strconv.ParseInt(sizeStr, 10, 64)

	seekLen := fInfo.Size() - (int64(len(fullString)) + archiveSize)
	curOff, _ := fPointer.Seek(seekLen, 0)
	log.Printf("Seeked to: %v for %v\r\n", curOff, archiveSize)

	zipBuf := bytes.NewBuffer(nil)
	_, err := io.CopyN(zipBuf, fPointer, archiveSize)

	if err != nil {
		return nil, 0
	}

	return bytes.NewReader(zipBuf.Bytes()), archiveSize
}

func readMetaFromFile(fPointer *os.File, fInfo os.FileInfo) []string {
	fPointer.Seek(fInfo.Size()-1024, 0)
	buf := make([]byte, 1024)

	_, err := fPointer.Read(buf)
	if err != nil {
		return []string{}
	}

	contents := string(buf)
	re := regexp.MustCompile("ArchiveLength:(\\d+)$")
	info := re.FindStringSubmatch(contents)

	if len(info) != 2 {
		return []string{}
	} else {
		return info
	}
}
