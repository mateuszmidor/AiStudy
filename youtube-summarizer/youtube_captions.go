package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// YouTubeTranscript represents the root element of the XML document.
type YouTubeTranscript struct {
	XMLName xml.Name      `xml:"transcript"`
	Text    []YouTubeText `xml:"text"`
}

// YouTubeText represents a text element within the transcript.
type YouTubeText struct {
	XMLName xml.Name `xml:"text"`
	Start   string   `xml:"start,attr"`
	Dur     string   `xml:"dur,attr"`
	Content string   `xml:",innerxml"`
}

func getCaptions(videoURL string) ([]string, error) {
	// Fetch the video webpage
	resp, err := http.Get(videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the video page: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve the video page. Status code: %d", resp.StatusCode)
	}

	// Parse the video webpage
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the video page: %v", err)
	}

	// Search for the captions URL
	var captionsURL string
	var findCaptionsURL func(*html.Node)
	findCaptionsURL = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "type" && a.Val == "application/ld+json" {
					continue
				}
			}
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "captionTracks") {
				re := regexp.MustCompile(`"captionTracks":(\[.*?\])`)
				match := re.FindStringSubmatch(n.FirstChild.Data)
				if len(match) > 1 {
					captionTracks := match[1]
					re := regexp.MustCompile(`"baseUrl":"(.*?)"`)
					urlMatch := re.FindStringSubmatch(captionTracks)
					if len(urlMatch) > 1 {
						captionsURL = strings.Replace(urlMatch[1], `\u0026`, "&", -1)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findCaptionsURL(c)
		}
	}
	findCaptionsURL(doc)

	if captionsURL == "" {
		return nil, fmt.Errorf("no captions found for this video")
	}

	// Fetch the captions XML
	// fmt.Println(captionsURL)
	captionsResp, err := http.Get(captionsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve captions: %v", err)
	}
	defer captionsResp.Body.Close()
	if captionsResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve captions. Status code: %d", captionsResp.StatusCode)
	}

	captionsData, err := ioutil.ReadAll(captionsResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read captions data: %v", err)
	}

	return extractCaptions(string(captionsData)), nil
}

func extractCaptions(inputXML string) []string {
	inputXML = strings.ReplaceAll(inputXML, "\n", " ")
	var transcript YouTubeTranscript

	// Unmarshal the inputXML into the transcript variable
	err := xml.Unmarshal([]byte(inputXML), &transcript)
	if err != nil {
		fmt.Printf("Failed to parse XML: %v\n", err)
		return nil
	}

	// Extract and return the Content of each Text element
	var captions []string
	for _, text := range transcript.Text {
		captions = append(captions, text.Content)
	}

	return captions
}
