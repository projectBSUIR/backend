package zipper

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

const TEMPDIRECTORY = "D:/Projects/BSUIR/backend/temp/"

func makeTempDir(name string) string {
	return TEMPDIRECTORY + name
}

func ExtractAllInOrder(file multipart.File, paths, names []string) ([][]byte, error) {
	buf := new(bytes.Buffer)
	fileSize, err := io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	zipFile, err := zip.NewReader(bytes.NewReader(buf.Bytes()), fileSize)
	if err != nil {
		return nil, err
	}
	var pathIndex int = 0
	zips := make([]*os.File, len(paths))
	zipWriters := make([]*zip.Writer, len(paths))

	CloseZips := func() {
		for i := 0; i < len(paths); i++ {
			zipWriters[i].Close()
			zips[i].Close()
		}
	}

	for i := 0; i < len(paths); i++ {
		zips[i], err = os.Create(makeTempDir(names[i]))
		if err != nil {
			return nil, err
		}
		zipWriters[i] = zip.NewWriter(zips[i])
	}

	for _, f := range zipFile.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if pathIndex+1 < len(paths) && strings.HasPrefix(f.Name, paths[pathIndex+1]) {
			pathIndex++
		}
		if strings.HasPrefix(f.Name, paths[pathIndex]) {
			err := CloneFile(f, paths[pathIndex], zipWriters[pathIndex])
			if err != nil {
				CloseZips()
				return nil, err
			}
		}
	}
	CloseZips()
	byteFiles := make([][]byte, len(paths))
	for i := 0; i < len(paths); i++ {
		byteFiles[i], err = os.ReadFile(makeTempDir(names[i]))
	}
	return byteFiles, nil
}

func CloneFile(f *zip.File, path string, testsWriter *zip.Writer) error {
	strings.TrimRightFunc(path, func(r rune) bool {
		return r != '/'
	})

	test, err := f.Open()
	if err != nil {
		return err
	}
	defer test.Close()
	nameTest, _ := strings.CutPrefix(f.Name, path)
	zipTest, err := testsWriter.Create(nameTest)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipTest, test)
	if err != nil {
		return err
	}
	return nil
}
