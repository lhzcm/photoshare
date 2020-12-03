package models

import "encoding/xml"

type KeliResult struct {
	XMLName    xml.Name `xml:"mob"`
	Version    string   `xml:"version,attr"`
	Result     int      `xml:"result"`
	ResultInfo string   `xml:"resultInfo"`
}
