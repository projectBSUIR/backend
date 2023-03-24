package routes

import (
	"archive/zip"
	"bytes"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"os"
	"strings"
)

func ContestHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("Contest")
	log.Println(fileHeader.Size)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	file, err := fileHeader.Open()

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
	testsFile, _ := os.Create("tests.zip")
	testsWriter := zip.NewWriter(testsFile)
	defer testsWriter.Close()
	defer testsFile.Close()

	for _, f := range contestZip.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name, "tests/") {
			testsWriter.Copy(f)
		}
	}
	defer os.Remove("tests.zip")
	return c.Status(fiber.StatusOK).SendFile("tests.zip")
}
