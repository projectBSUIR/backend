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
	var testsZip *os.File
	var zipWriter *zip.Writer

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

	checkerTempFile, _ := os.CreateTemp(TEMPDIRECTORY, names[1])
	fileInfo, err = os.Stat(propertiesTempFile.Name())
	if err != nil {
		return nil, err
	}

	checker, err := os.OpenFile(checkerTempFile.Name(),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		fileInfo.Mode())

	if err != nil {
		return nil, err
	}

	CloseZips := func() {
		properties.Close()
		propertiesTempFile.Close()
		checker.Close()
		checkerTempFile.Close()
		nameFile := testsZip.Name()
		testsZip.Close()
		os.Remove(nameFile)
		os.Remove(checkerTempFile.Name())
		os.Remove(propertiesTempFile.Name())
	}

	defer CloseZips()

	testsZip, err = os.CreateTemp(TEMPDIRECTORY, names[0])
	if err != nil {
		return nil, err
	}
	zipWriter = zip.NewWriter(testsZip)

	for _, f := range zipFile.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name, paths[0]) {
			err := CloneFileToZip(f, paths[0], zipWriter)
			if err != nil {
				return nil, err
			}
		}

		if strings.HasPrefix(f.Name, paths[1]) {
			err := CloneFileToNoneZip(f, checker)
			if err != nil {
				return nil, err
			}
		}

		if strings.HasPrefix(f.Name, paths[2]) {
			err = CloneFileToNoneZip(f, properties)
			if err != nil {
				return nil, err
			}
		}
	}
	byteFiles := make([][]byte, len(paths))
	zipWriter.Close()
	byteFiles[0], err = os.ReadFile(testsZip.Name())
	if err != nil {
		return nil, err
	}
	byteFiles[1], err = os.ReadFile(checkerTempFile.Name())
	if err != nil {
		return nil, err
	}
	byteFiles[2], err = os.ReadFile(propertiesTempFile.Name())
	if err != nil {
		return nil, err
	}
	return byteFiles, nil
}

func CloneFileToNoneZip(f *zip.File, targetFile *os.File) error {
	fileReader, err := f.Open()
	if err != nil {
		return err
	}

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	fileReader.Close()
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
	nameFile, _ := strings.CutPrefix(f.Name, path)
	if err != nil {
		return err
	}

	zipFile, err := zipFileWriter.Create(nameFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, file)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
