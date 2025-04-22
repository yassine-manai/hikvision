package hikvision

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

// ParseXMLData parses the XML data from Hikvision camera
func ParseXMLData(xmlData string) (*EventNotificationAlert, error) {
	var alert EventNotificationAlert
	err := xml.Unmarshal([]byte(xmlData), &alert)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}
	return &alert, nil
}

// ExtractCaptureFromAlert converts EventNotificationAlert to Capture struct
func ExtractCaptureFromAlert(alert *EventNotificationAlert) *Capture {
	return &Capture{
		State:        alert.ANPR.Country,
		LicensePlate: alert.ANPR.LicensePlate,
		Direction:    alert.ANPR.Direction,
		Confidence:   alert.ANPR.ConfidenceLevel,
		CamIP:        alert.IPAddress,
		CaptureTime:  time.Now().Format(time.RFC3339),
		VehicleType:  alert.ANPR.VehicleType,
		XMLData:      "", // Will be set separately if needed
	}
}

// ExtractImagesFromForm processes multipart form to extract images
func ExtractImagesFromForm(form *multipart.Form) ([]Image, error) {
	var images []Image

	// Process license plate image
	if files, ok := form.File["licensePlatePicture.jpg"]; ok && len(files) > 0 {
		img, err := processImageFile(files[0], "licensePlate")
		if err != nil {
			return nil, fmt.Errorf("failed to process license plate image: %w", err)
		}
		images = append(images, img)
	}

	// Process detection image
	if files, ok := form.File["detectionPicture.jpg"]; ok && len(files) > 0 {
		img, err := processImageFile(files[0], "detection")
		if err != nil {
			return nil, fmt.Errorf("failed to process detection image: %w", err)
		}
		images = append(images, img)
	}

	if len(images) == 0 {
		return nil, errors.New("no images found in the request")
	}

	return images, nil
}

func processImageFile(fileHeader *multipart.FileHeader, imageType string) (Image, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return Image{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Image{}, err
	}

	return Image{
		Type:        imageType,
		Data:        data,
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}, nil
}
