package models

import "encoding/xml"

type Urlset struct {
	XMLName xml.Name  `xml:"urlset"`
	Text    string    `xml:",chardata"`
	Xmlns   string    `xml:"xmlns,attr"`
	URL     []URLItem `xml:"url"`
}

type URLItem struct {
	Text       string  `xml:",chardata"`
	Loc        string  `xml:"loc"`
	Lastmod    string  `xml:"lastmod"`
	Changefreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}
