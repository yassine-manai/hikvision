package hikvision

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// saveXMLData saves the XML data to a file
func saveXMLData(xmlData, licensePlate, clientIP string) error {
	// Create directory if it doesn't exist
	dir := filepath.Join("storage", "xml", clientIP)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create filename with timestamp
	filename := filepath.Join(dir, fmt.Sprintf("%s_%s.xml",
		time.Now().Format("20060102_150405"),
		licensePlate))

	// Write XML data to file
	if err := ioutil.WriteFile(filename, []byte(xmlData), 0644); err != nil {
		return fmt.Errorf("failed to write XML file: %w", err)
	}

	return nil
}

// saveImages saves the images from the capture
func saveImages(images []Image, clientIP string) error {
	for _, img := range images {
		// Create directory if it doesn't exist
		dir := filepath.Join("storage", "images", clientIP, img.Type)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Create filename with timestamp
		filename := filepath.Join(dir, fmt.Sprintf("%s_%s.jpg",
			time.Now().Format("20060102_150405"),
			img.FileName))

		// Write image data to file
		if err := ioutil.WriteFile(filename, img.Data, 0644); err != nil {
			return fmt.Errorf("failed to write image file: %w", err)
		}
	}
	return nil
}

// ReadFileContent reads content from a multipart file header
func ReadFileContent(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	return ioutil.ReadAll(src)
}
