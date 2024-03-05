package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/BNPrashanth/h24-assignment/models"
)

const genericErrorMessage = "Oops, something went wrong. Please try again later."

func HandleAnalyseWebPage(w http.ResponseWriter, r *http.Request) {
	params := models.ParseRequestParams{}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Error occured while decoding parameters:", err.Error())

		parseError := models.ParseResponseError{
			Message:    genericErrorMessage,
			DevMessage: err.Error(),
			Success:    false,
		}

		json.NewEncoder(w).Encode(parseError)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if params.URL == "" {
		fmt.Println("No url provided in the request body.")

		parseError := models.ParseResponseError{
			Message: "Please provide a valid url to analyse.",
			Success: false,
		}

		json.NewEncoder(w).Encode(parseError)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	fmt.Println("URL received:", params.URL)

	resp, err := http.Get(params.URL)
	if err != nil {
		fmt.Println("Error occured while sending request:", err.Error())

		parseError := models.ParseResponseError{
			Message:    genericErrorMessage,
			DevMessage: err.Error(),
			Success:    false,
		}

		json.NewEncoder(w).Encode(parseError)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Got a non 200 response:", resp.StatusCode)

		parseError := models.ParseResponseError{
			Message: "Oops, the provided url is either invalid or failed to load. Please check the url and try again.",
			Status:  resp.StatusCode,
			Success: false,
		}

		json.NewEncoder(w).Encode(parseError)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error occured while reading response body:", err.Error())

		parseError := models.ParseResponseError{
			Message:    genericErrorMessage,
			DevMessage: err.Error(),
			Success:    false,
		}

		json.NewEncoder(w).Encode(parseError)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(data))
	// fmt.Println(string(data))

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error occured while parsing the response body via goquery:", err.Error())

		parseError := models.ParseResponseError{
			Message:    genericErrorMessage,
			DevMessage: err.Error(),
			Success:    false,
		}

		json.NewEncoder(w).Encode(parseError)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	parseResponse := models.ParseResponse{}

	// Get the HTML version
	version := getHTMLVersion(string(data))
	fmt.Println("HTML version:", version)
	parseResponse.HTMLVersion = version

	// Get the title
	titleSelection := doc.Find("title").First()
	if titleSelection != nil {
		fmt.Println("Title:", titleSelection.Text())
		parseResponse.Title = titleSelection.Text()
	}

	// Get the headings
	parseResponse.Headings = getHeadings(*doc)

	// Get the links
	links := make(chan models.ParseResponseLink)
	go getLinks(params.URL, *doc, links)
	parseResponse.Links = <-links

	// Check if there is a login form
	parseResponse.HasLoginForm = checkForLoginForm(*doc)

	parseResponse.Success = true

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parseResponse)
}

func getHTMLVersion(htmlString string) string {
	htmlString = strings.ToLower(htmlString)

	docTypeIndex := strings.Index(htmlString, "<!doctype html")
	htmlTagIndex := strings.Index(htmlString, "<html")

	// Check if the document is HTML5
	if strings.Index(htmlString, "<!doctype html>") < strings.Index(htmlString, "<html") {
		return "HTML5"
	}

	// Check if the document is XHTML
	dtdXHtmlIndex := strings.Index(htmlString, "dtd xhtml 1.1")
	if docTypeIndex < htmlTagIndex && dtdXHtmlIndex > docTypeIndex && dtdXHtmlIndex < htmlTagIndex {
		return "XHTML"
	}

	// Check if the document is HTML4
	dtdHtml4Index := strings.Index(htmlString, "dtd html 4.01")
	if docTypeIndex < htmlTagIndex && dtdHtml4Index > docTypeIndex && dtdHtml4Index < htmlTagIndex {
		return "HTML 4.1"
	}

	// Check if the document is HTML3.2
	dtdHtml32Index := strings.Index(htmlString, "dtd html 3.2")
	if docTypeIndex < htmlTagIndex && dtdHtml32Index > docTypeIndex && dtdHtml32Index < htmlTagIndex {
		return "HTML 3.2"
	}

	// Check if the document is HTML2.0
	dtdHtml2Index := strings.Index(htmlString, "dtd html 2.0")
	if docTypeIndex < htmlTagIndex && dtdHtml2Index > docTypeIndex && dtdHtml2Index < htmlTagIndex {
		return "HTML 2.0"
	}

	// Return HTML version 1.0 for any other
	return "HTML 1.0"
}

func getHeadings(doc goquery.Document) models.ParseResponseHeadings {
	headingsResult := models.ParseResponseHeadings{}

	headingCount := 0
	headings := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	for _, heading := range headings {
		headingTag := models.ParseResponseHeading{
			Tag: heading,
		}

		doc.Find(heading).Each(func(i int, s *goquery.Selection) {
			headingTag.Count++
			headingTag.List = append(headingTag.List, s.Text())
		})

		headingCount += headingTag.Count

		switch heading {
		case "h1":
			headingsResult.H1 = headingTag
		case "h2":
			headingsResult.H2 = headingTag
		case "h3":
			headingsResult.H3 = headingTag
		case "h4":
			headingsResult.H4 = headingTag
		case "h5":
			headingsResult.H5 = headingTag
		case "h6":
			headingsResult.H6 = headingTag
		}
	}

	headingsResult.Count = headingCount

	return headingsResult
}

func getLinks(url string, doc goquery.Document, linksResp chan<- models.ParseResponseLink) {
	linksResponse := models.ParseResponseLink{}

	loadURL := url
	protocol := "https"
	domainComponent := ""
	baseDomain := ""
	if strings.Contains(url, "://") {
		urlParts := strings.Split(url, "://")
		if len(urlParts) >= 2 {
			protocol = urlParts[0]

			if !strings.Contains(urlParts[1], "/") {
				loadURL = url + "/"
			}

			domainComponent = strings.Split(urlParts[1], "/")[0]
			domainParts := strings.Split(domainComponent, ".")
			if len(domainParts) >= 2 {
				baseDomain = domainParts[len(domainParts)-2] + "." + domainParts[len(domainParts)-1]
			}
		}
	}

	hrefValuesToCheck := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// fmt.Println("a", s.Text())
		linksResponse.Count++
		isInternal := false

		hrefValue, exists := s.Attr("href")
		if exists && hrefValue != "" {
			// fmt.Println("HREF", hrefValue)

			// Check if link is internal or external - 1
			// If there is no domain and is calling a file directly, it can be categorised as an internal link
			// because it'll try to load with the domain/sub-domain of the current page
			if !strings.Contains(hrefValue, "://") {
				linksResponse.InternalCount++
				linksResponse.InternalLinks = append(linksResponse.InternalLinks, hrefValue)
				isInternal = true
			} else {
				// Check if link is internal or external - 2
				// Check the base domain of the href value

				// Check if the base domain if the same in the href value
				// Let's ignore the protocol
				hrefParts := strings.Split(hrefValue, "://")
				if len(hrefParts) >= 2 {
					hrefDomainComponent := strings.Split(hrefParts[1], "/")[0]
					if strings.Contains(hrefDomainComponent, baseDomain) {
						// Internal link
						linksResponse.InternalCount++
						linksResponse.InternalLinks = append(linksResponse.InternalLinks, hrefValue)
						isInternal = true
					} else {
						// External link
						linksResponse.ExternalCount++
						linksResponse.ExternalLinks = append(linksResponse.ExternalLinks, hrefValue)
					}
				}

				// No need to handle the else case here and update the inactive links count
				// The inactive links count will reflect this case when the href value fails to load below
			}

			// Check if the link is active/inactive
			// The tag has a href link with a not null value, but it doesn't load
			// The tag does not have a href value at all
			hrefUrlToLoad := hrefValue
			if isInternal && !strings.Contains(hrefUrlToLoad, "://") {
				if strings.Index(hrefValue, "/") != 0 && !strings.Contains(hrefValue, domainComponent) {
					hrefUrlToLoad = loadURL[0:strings.LastIndex(loadURL, "/")] + "/" + hrefValue
				}

				if strings.Index(hrefValue, "/") == 0 && !strings.Contains(hrefValue, domainComponent) {
					hrefUrlToLoad = protocol + "://" + domainComponent + hrefValue
				}
			}
			// fmt.Println("hrefUrlToLoad", hrefUrlToLoad)
			hrefValuesToCheck = append(hrefValuesToCheck, hrefUrlToLoad)
		} else {
			linksResponse.InactiveCount++
			linksResponse.InactiveLinks = append(linksResponse.InactiveLinks, models.ParseResponseLinkInactive{
				Link: s.Text(),
			})
		}
	})

	// Identify inactive links from the exisitng href values
	if len(hrefValuesToCheck) > 0 {
		responseCount := 0
		ch := make(chan *models.LinksResponse)
		go checkAllLinks(hrefValuesToCheck, ch)
		for responseCount != len(hrefValuesToCheck) {
			select {
			case r := <-ch:
				if r.Err != nil {
					linksResponse.InactiveCount++
					linksResponse.InactiveLinks = append(linksResponse.InactiveLinks, models.ParseResponseLinkInactive{
						Link: r.Link,
					})
				} else {
					if r.Resp.StatusCode != 200 {
						linksResponse.InactiveCount++
						linksResponse.InactiveLinks = append(linksResponse.InactiveLinks, models.ParseResponseLinkInactive{
							Link:   r.Link,
							Status: r.Resp.StatusCode,
						})
					}
				}
				responseCount++
			case <-time.After(5 * time.Second):
				os.Exit(1)
			}
		}
	}

	linksResp <- linksResponse
}

func checkForLoginForm(doc goquery.Document) bool {
	// Check for input fields and at least one should be of type password
	// Check for a button of type submit
	// If both are there then we can conclude that there is a login form

	inputFieldsCount := 0
	hasPasswordField := false
	hasSubmitButton := false
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		inputType, exists := s.Attr("type")
		if exists {
			if strings.ToLower(inputType) != "hidden" {
				inputFieldsCount++
			}

			if strings.ToLower(inputType) == "password" {
				hasPasswordField = true
			}
		}
	})

	doc.Find("button").Each(func(i int, s *goquery.Selection) {
		buttonType, exists := s.Attr("type")
		if exists && strings.ToLower(buttonType) == "submit" {
			hasSubmitButton = true
		}
	})

	return inputFieldsCount >= 2 && hasPasswordField && hasSubmitButton
}

func checkAllLinks(links []string, ch chan *models.LinksResponse) {
	for _, link := range links {
		go func(link string) {
			resp, err := http.Get(link)
			ch <- &models.LinksResponse{
				Link: link,
				Resp: resp,
				Err:  err,
			}
		}(link)
	}
}
