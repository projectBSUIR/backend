package zipper

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

var TEMPDIRECTORY string = "./temp/"

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
	zips := make([]*os.File, len(paths)-1)
	zipWriters := make([]*zip.Writer, len(paths)-1)

	propertiesTempFile, _ := os.CreateTemp(TEMPDIRECTORY, names[2])
	fileInfo, err := os.Stat(propertiesTempFile.Name())
	if err != nil {
		return nil, err
	}

	properties, err := os.OpenFile(propertiesTempFile.Name(),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		fileInfo.Mode())

	if err != nil {
		return nil, err
	}

	CloseZips := func() {
		for i := 0; i < len(zips); i++ {
			zipWriters[i].Close()
			os.Remove(zips[i].Name())
		}
		os.Remove(propertiesTempFile.Name())
		os.RemoveAll(TEMPDIRECTORY)
	}

	for i := 0; i < len(zips); i++ {
		zips[i], err = os.CreateTemp(TEMPDIRECTORY, names[i])
		if err != nil {
			return nil, err
		}
		zipWriters[i] = zip.NewWriter(zips[i])
	}

	for _, f := range zipFile.File {
		if f.FileInfo().IsDir() {
			continue
		}
		for pathIndex := 0; pathIndex < len(zips); pathIndex++ {
			if strings.HasPrefix(f.Name, paths[pathIndex]) {
				err := CloneFileToZip(f, paths[pathIndex], zipWriters[pathIndex])
				if err != nil {
					CloseZips()
					return nil, err
				}
			}
		}

		if strings.HasPrefix(f.Name, paths[2]) {
			CloneFileToNoneZip(f, properties)
		}
	}
	byteFiles := make([][]byte, len(paths))
	for i := 0; i < len(zips); i++ {
		byteFiles[i], err = os.ReadFile(zips[i].Name())
		if err != nil {
			return nil, err
		}
	}
	byteFiles[2], err = os.ReadFile(propertiesTempFile.Name())
	if err != nil {
		return nil, err
	}
	CloseZips()
	return byteFiles, nil
}

func CloneFileToNoneZip(f *zip.File, targetFile *os.File) error {
	fileReader, err := f.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	return nil
}

func CloneFileToZip(f *zip.File, path string, zipFileWriter *zip.Writer) error {
	path = strings.TrimRightFunc(path, func(r rune) bool {
		return r != '/'
	})

	file, err := f.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	nameFile, _ := strings.CutPrefix(f.Name, path)
	zipFile, err := zipFileWriter.Create(nameFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, file)
	if err != nil {
		return err
	}
	return nil
}

func CreateTempDir() error {
	var err error
	TEMPDIRECTORY, err = os.MkdirTemp(".", "temp")
	return err
}
