package netgo

import (
	//"os"
	//"path/filepath"
	"encoding/json"
	"strconv"
	"strings"
	//"github.com/op/go-logging"
)

//var log = logging.MustGetLogger("netgo")

func TpbGetResultsJson(in string) []SearchResult {

	type JsonResultHeader struct {
		Class   string            `json:"class"`
		Headers []json.RawMessage `json:"th"`
	}

	type JsonImage struct {
		Alt    string `json:"alt"`
		Src    string `json:"src"`
		Border string `json:"border"`
		Style  string `json:"style"`
		Title  string `json:"title"`
	}

	type JsonLinkInfo struct {
		Class string    `json:"class"`
		Link  string    `json:"href"`
		Title string    `json:"title"`
		Name  string    `json:"content"`
		Image JsonImage `json:"img"`
	}

	type JsonCategoriesAlt struct {
		LinkInfo JsonLinkInfo `json:"a"`
	}
	type JsonCategories struct {
		LinkInfo     JsonLinkInfo      `json:"a"`
		CategoryNull int               `json:"br"`
		CategoryAlt  JsonCategoriesAlt `json:"p"`
	}

	type JsonCategoriesList struct {
		Class      string         `json:"class"`
		Categories JsonCategories `json:"center"`
	}

	type JsonResultTitle struct {
		Class    string       `json:"class"`
		LinkInfo JsonLinkInfo `json:"a"`
	}

	type JsonResultDetail struct {
		Class   string       `json:"class"`
		Link    JsonLinkInfo `json:"a"`
		Content string       `json:"content"`
	}

	type JsonResultInfo struct {
		Info   JsonResultTitle  `json:"div"`
		Links  json.RawMessage  `json:"a"`
		Images json.RawMessage  `json:"img"`
		Detail JsonResultDetail `json:"font"`
	}

	type JsonCount struct {
		Align string `json:"align"`
		Count string `json:"p"`
	}

	type JsonResultEntry struct {
		ResultData []json.RawMessage `json:"td"`
	}

	type JsonResults struct {
		Record []json.RawMessage `json:"tr"`
	}

	type JsonQueryEntry struct {
		Count        int         `json:"count"`
		Created      string      `json:"created"`
		Lang         string      `json:"lang"`
		ResultsEntry JsonResults `json:"results"`
	}

	type JsonQueryResult struct {
		QueryEntry JsonQueryEntry `json:"query"`
	}

	var q JsonQueryResult
	err := json.Unmarshal([]byte(in), &q)

	if err != nil {
		log.Fatal(err)
	}

	if len(q.QueryEntry.ResultsEntry.Record) == 0 {
		log.Info("No records matched the search", err)
		return nil
	}

	var header JsonResultHeader
	var entries []SearchResult

	err = json.Unmarshal(q.QueryEntry.ResultsEntry.Record[0], &header)
	if err != nil {
		log.Error("error:", err)
		return entries
	}

	for t, record := range q.QueryEntry.ResultsEntry.Record[1:] {
		var e JsonResultEntry
		var cat JsonCategoriesList
		var info JsonResultInfo
		var link JsonLinkInfo
		var links []JsonLinkInfo
		var seed JsonCount
		var leech JsonCount
		var image JsonImage
		var images []JsonImage
		var entry SearchResult

		err := json.Unmarshal(record, &e)
		if err != nil {
			log.Error("[%d] [error: %+v]\n%v\n\n", t, err, e)
			continue
		}

		// Get Category
		err = json.Unmarshal(e.ResultData[0], &cat)
		if err != nil {
			log.Error("error:", err)
		}

		// Get File Info
		err = json.Unmarshal(e.ResultData[1], &info)
		if err != nil {
			log.Error("error:", err)
		}

		// Get Links
		err = json.Unmarshal(info.Links, &links)
		if err != nil {
			err = json.Unmarshal(info.Links, &link)
			if err != nil {
				log.Error("error:", err)
			}
			links = append(links, link)
		}
		log.Debug("[%d-Links] %v\n\n", t, links)

		// Get Images
		err = json.Unmarshal(info.Images, &images)
		if err != nil {
			err = json.Unmarshal(info.Images, &image)
			if err != nil {
				log.Error("error:", err)
			}
			images = append(images, image)
		}
		log.Debug("[%d-Images] %v\n\n", t, images)

		// Get Seeds
		err = json.Unmarshal(e.ResultData[2], &seed)
		if err != nil {
			log.Error("error:", err)
		}
		// Get Leeches
		err = json.Unmarshal(e.ResultData[3], &leech)
		if err != nil {
			log.Error("error:", err)
		}

		entry.Name = info.Info.LinkInfo.Name
		entry.Link = info.Info.LinkInfo.Link
		entry.Category = cat.Categories.LinkInfo.Name + "_" + cat.Categories.CategoryAlt.LinkInfo.Name
		entry.Seed, _ = strconv.Atoi(seed.Count)
		entry.Leech, _ = strconv.Atoi(leech.Count)

		for link_index, link := range links {
			log.Debug("[link %d] %v", link_index, link)

			if strings.Contains(link.Image.Alt, "Magnet") {
				entry.Magnet = link.Link
			}
			if strings.Contains(link.Image.Alt, "VIP") {
				entry.VIP = true
			}
			if strings.Contains(link.Image.Alt, "Trusted") {
				entry.Trusted = true
			}
			if strings.Contains(link.Image.Alt, "Download") {
				entry.Torrent = link.Link
			}
		}

		for img_index, img := range images {
			log.Debug("[img %d] %v", img_index, img)

			if strings.Contains(img.Title, "comments") {
				entry.Comments, _ = strconv.Atoi(string(img.Title[17]))
			}
			if strings.Contains(img.Title, "cover image") {
				entry.CoverImage = true
			}
		}

		// Get User
		entry.User = "Anonymous"
		if "" != info.Detail.Link.Name {
			entry.User = info.Detail.Link.Name
		}

		// Get Date

		// Get Size

		log.Debug("[%d] %+v\n\n", t, entry)

		entries = append(entries, entry)
	}

	for i, record := range entries {
		log.Debug("[%d] %v\n", i, record)
	}

	log.Debug("[output]\n%v\n", header.Class)

	return entries
}
