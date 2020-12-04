package models

import "encoding/xml"

const (
	KeliXmlHead = `<?xml version="1.0" encoding="GB2312"?>` + "\n"
)

type KeliResult struct {
	XMLName    xml.Name `xml:"mob"`
	Version    string   `xml:"version,attr"`
	Result     int      `xml:"result"`
	ResultInfo string   `xml:"resultInfo"`
}
