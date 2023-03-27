package routes

import (
	"archive/zip"
	"bytes"
	"github.com/gofiber/fiber/v2"
	"io"
	"os"
	"strings"
)

func ContestHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("Contest")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	file, err := fileHeader.Open()
	defer file.Close()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	buf := new(bytes.Buffer)

	fileSize, err := io.Copy(buf, file)
	if err != nil {
		panic(err)
	}

	contestZip, err := zip.NewReader(bytes.NewReader(buf.Bytes()), fileSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	zipFile, err := os.Create("D:/Projects/BSUIR/backend/temp/tests.zip")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	testsWriter := zip.NewWriter(zipFile)

	for _, f := range contestZip.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name, "tests/") {
			err := CloneTest(f, testsWriter)
			if err != nil {
				testsWriter.Close()
				zipFile.Close()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			}
		}
	}

	c.SendFile("D:/Projects/BSUIR/backend/temp/tests.zip")

	testsWriter.Close()
	zipFile.Close()
	os.Remove("D:/Projects/BSUIR/backend/temp/")
	return c.SendStatus(fiber.StatusOK)
}

func CloneTest(f *zip.File, testsWriter *zip.Writer) error {
	test, err := f.Open()
	if err != nil {
		return err
	}
	defer test.Close()
	nameTest, _ := strings.CutPrefix(f.Name, "tests/")
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
