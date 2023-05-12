package zipper

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

const TEMPDIRECTORY = "./temp/"

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
	var tests *os.File
	var testsWriter *zip.Writer

	checker, _ := os.CreateTemp(TEMPDIRECTORY, names[1])
	fileInfo, err := os.Stat(checker.Name())
	if err != nil {
		return nil, err
	}

	propertiesTempFile, _ := os.CreateTemp(TEMPDIRECTORY, names[2])
	fileInfo, err = os.Stat(propertiesTempFile.Name())
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
		testsWriter.Close()
		os.Remove(tests.Name())
		os.Remove(checker.Name())
		os.Remove(propertiesTempFile.Name())
	}

	tests, err = os.CreateTemp(TEMPDIRECTORY, names[0])
	if err != nil {
		return nil, err
	}
	testsWriter = zip.NewWriter(tests)

	for _, f := range zipFile.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name, paths[0]) {
			err := CloneFileToZip(f, paths[0], testsWriter)
			if err != nil {
				CloseZips()
				return nil, err
			}
		}

		if strings.HasPrefix(f.Name, paths[1]) {
			err := CloneFileToNoneZip(f, checker)
			if err != nil {
				CloseZips()
				return nil, err
			}
		}

		if strings.HasPrefix(f.Name, paths[2]) {
			err := CloneFileToNoneZip(f, properties)
			if err != nil {
				CloseZips()
				return nil, err
			}
		}
	}
	byteFiles := make([][]byte, len(paths))
	byteFiles[1], err = os.ReadFile(checker.Name())
	if err != nil {
		return nil, err
	}

	byteFiles[0], err = os.ReadFile(tests.Name())
	if err != nil {
		return nil, err
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
