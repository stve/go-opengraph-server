package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rojters/opengraph"
)

// Media - core data structure for image/video data
type Media struct {
	Height    string `json:"height,omitempty"`
	SecureURL string `json:"secure_url,omitempty"`
	Type      string `json:"type,omitempty"`
	URL       string `json:"url,omitempty"`
	Width     string `json:"width,omitempty"`
}

func (m *Media) isIncomplete() bool {
	return isEmpty(m.Height) || isEmpty(m.SecureURL) || isEmpty(m.Type) || isEmpty(m.URL) || isEmpty(m.Width)
}

// GraphData - core data structure for top-level opengraph data
type GraphData struct {
	SiteName    string  `json:"site_name,omitempty"`
	Type        string  `json:"type,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	URL         string  `json:"url,omitempty"`
	Videos      []Media `json:"videos,omitempty"`
	Images      []Media `json:"images,omitempty"`
}

// IncompleteVideo - return the first incomplete Video media to assign properties to
func (r *GraphData) IncompleteVideo() *Media {
	foundIndex := -1
	for i := 0; i < len(r.Videos); i++ {
		if r.Videos[i].isIncomplete() {
			foundIndex = i
			break
		}
	}

	if foundIndex >= 0 {
		return &r.Videos[foundIndex]
	}

	if len(r.Videos) == 0 {
		r.Videos = []Media{Media{}}
		return &r.Videos[0]
	}

	r.Videos = append(r.Videos, Media{})
	return &r.Videos[len(r.Videos)-1]
}

// IncompleteImage - return the first incomplete Image media to assign properties to
func (r *GraphData) IncompleteImage() *Media {
	foundIndex := -1
	for i := 0; i < len(r.Images); i++ {
		if r.Images[i].isIncomplete() {
			foundIndex = i
			break
		}
	}

	if foundIndex >= 0 {
		return &r.Images[foundIndex]
	}

	if len(r.Images) == 0 {
		r.Images = []Media{Media{}}
		return &r.Images[0]
	}

	r.Images = append(r.Images, Media{})
	return &r.Images[len(r.Images)-1]
}

// Error - default error message type
type Error struct {
	Message string `json:"message"`
}

const opengraphPrefix string = "og"
const keySeparator string = ":"
const defaultPort string = "8000"

func extract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(w)

	url := r.URL.Query().Get("url")
	if isEmpty(url) {
		w.WriteHeader(http.StatusBadRequest)
		writeError(encoder, "Invalid URL")
		return
	}

	res, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(encoder, "Error accessing URL: "+err.Error())
		return
	}

	if res.StatusCode != 200 {
		w.WriteHeader(http.StatusBadRequest)
		writeError(encoder, "Invalid URL")
		return
	}

	metaData, err := opengraph.ExtractPrefix(res.Body, opengraphPrefix)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(encoder, err.Error())
		return
	}

	structuredData := structureMetaData(metaData)

	response, err := json.Marshal(structuredData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(encoder, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return
}

func writeError(encoder *json.Encoder, message string) {
	if err := encoder.Encode(Error{Message: message}); err != nil {
		panic(err)
	}
}

func structureMetaData(md []opengraph.MetaData) GraphData {
	data := GraphData{}

	for i := range md {
		property := md[i].Property

		switch property {
		case "site_name":
			data.SiteName = md[i].Content
		case "title":
			data.Title = md[i].Content
		case "type":
			data.Type = md[i].Content
		case "description":
			data.Description = md[i].Content
		case "url":
			data.URL = md[i].Content
		case "image":
			m := data.IncompleteImage()
			m.URL = md[i].Content
		default:
			if isNested(property) {
				handleNestedData(&data, md[i])
			}
		}
	}

	return data
}

func handleNestedData(data *GraphData, md opengraph.MetaData) {
	prop := strings.Split(md.Property, keySeparator)
	parent := prop[0]
	key := prop[1]

	switch parent {
	case "video":
		media := data.IncompleteVideo()
		switch key {
		case "secure_url":
			media.SecureURL = md.Content
		case "width":
			media.Width = md.Content
		case "height":
			media.Height = md.Content
		case "type":
			media.Type = md.Content
		case "url":
			media.URL = md.Content
		}
	case "image":
		media := data.IncompleteImage()
		switch key {
		case "secure_url":
			media.SecureURL = md.Content
		case "width":
			media.Width = md.Content
		case "height":
			media.Height = md.Content
		case "type":
			media.Type = md.Content
		case "url":
			media.URL = md.Content
		}
	}
}

func isNested(str string) bool {
	return strings.Contains(str, keySeparator)
}

func isEmpty(str string) bool {
	return (len(str) == 0)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	http.HandleFunc("/", extract)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
