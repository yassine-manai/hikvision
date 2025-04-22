package hikvision

import (
	"encoding/xml"
)

// EventNotificationAlert represents the XML structure from Hikvision camera
type EventNotificationAlert struct {
	XMLName              xml.Name `xml:"EventNotificationAlert"`
	Version              string   `xml:"version,attr"`
	Xmlns                string   `xml:"xmlns,attr"`
	IPAddress            string   `xml:"ipAddress"`
	PortNo               string   `xml:"portNo"`
	Protocol             string   `xml:"protocol"`
	MacAddress           string   `xml:"macAddress"`
	ChannelID            string   `xml:"channelID"`
	DateTime             string   `xml:"dateTime"`
	ActivePostCount      int      `xml:"activePostCount"`
	EventType            string   `xml:"eventType"`
	EventState           string   `xml:"eventState"`
	EventDescription     string   `xml:"eventDescription"`
	ChannelName          string   `xml:"channelName"`
	ANPR                 ANPR     `xml:"ANPR"`
	UUID                 string   `xml:"UUID"`
	PicNum               int      `xml:"picNum"`
	IsDataRetransmission bool     `xml:"isDataRetransmission"`
}

type ANPR struct {
	Country              string        `xml:"country"`
	Province             string        `xml:"province"`
	LicensePlate         string        `xml:"licensePlate"`
	Line                 int           `xml:"line"`
	Direction            string        `xml:"direction"`
	ConfidenceLevel      int           `xml:"confidenceLevel"`
	PlateType            string        `xml:"plateType"`
	PlateColor           string        `xml:"plateColor"`
	LicenseBright        int           `xml:"licenseBright"`
	VehicleType          string        `xml:"vehicleType"`
	DetectDir            int           `xml:"detectDir"`
	VehicleInfo          VehicleInfo   `xml:"vehicleInfo"`
	PictureInfoList      []PictureInfo `xml:"pictureInfoList>pictureInfo"`
	OriginalLicensePlate string        `xml:"originalLicensePlate"`
}

type VehicleInfo struct {
	Index               int    `xml:"index"`
	ColorDepth          int    `xml:"colorDepth"`
	Color               string `xml:"color"`
	Length              int    `xml:"length"`
	VehicleLogoRecog    int    `xml:"vehicleLogoRecog"`
	VehicleSubLogoRecog int    `xml:"vehileSubLogoRecog"`
	VehicleModel        int    `xml:"vehileModel"`
}

type PictureInfo struct {
	FileName  string    `xml:"fileName"`
	Type      string    `xml:"type"`
	DataType  int       `xml:"dataType"`
	AbsTime   string    `xml:"absTime"`
	PId       string    `xml:"pId"`
	PlateRect PlateRect `xml:"plateRect,omitempty"`
}

type PlateRect struct {
	X      int `xml:"X"`
	Y      int `xml:"Y"`
	Width  int `xml:"width"`
	Height int `xml:"height"`
}

// Capture represents the processed LPR data
type Capture struct {
	State        string  `json:"country"`
	LicensePlate string  `json:"licensePlate"`
	Direction    string  `json:"direction"`
	Confidence   int     `json:"confidenceLevel"`
	CamIP        string  `json:"ipAddress"`
	CaptureTime  string  `json:"captureTime"`
	VehicleType  string  `json:"vehicleType"`
	Images       []Image `json:"images,omitempty"`
	XMLData      string  `json:"xmlData,omitempty"`
}

type Image struct {
	Type        string `json:"type"` // "licensePlate" or "detection"
	Data        []byte `json:"-"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}
