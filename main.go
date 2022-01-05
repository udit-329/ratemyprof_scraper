package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseRatings(profId string) {
	url := "https://www.ratemyprofessors.com/ShowRatings.jsp?tid=" + profId
	fmt.Println(url)
}

func getProfId(jsonString string, first string, last string) (legacyIdStr string) {

	var result map[string]interface{}
	json.Unmarshal([]byte(jsonString), &result)

	for _, element := range result {

		// Only check maps
		if _, ok := element.(map[string]interface{}); ok {

			//Check if the map has 'firstName' and 'lastName' keys
			firstName, foundFirst := element.(map[string]interface{})["firstName"]
			lastName, foundLast := element.(map[string]interface{})["lastName"]

			if foundFirst && foundLast {
				if strings.Contains(strings.ToLower(firstName.(string)), first) && strings.Contains(strings.ToLower(lastName.(string)), last) {
					legacyId, _ := element.(map[string]interface{})["legacyId"]
					legacyIdStr := strconv.Itoa(int(legacyId.(float64)))

					return legacyIdStr
				}
			}
		}
	}
	return ""
}

func getProf(first string, last string) {

	url := "https://www.ratemyprofessors.com/search/teachers?query=" + first + "%20" + last + "&sid=U2Nob29sLTE0NjY="

	response, err := http.Get(url)
	// fmt.Println(url)
	if err != nil {
		fmt.Println(err)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)

	document.Find("script").Each(func(i int, s *goquery.Selection) {

		text := s.Text()

		// Get the script tag we want
		if strings.Contains(text, "window.__RELAY_STORE__") {

			//Split based on semicolons
			arrOfStrings := strings.Split(text, ";")

			//Get the string we want
			for _, str := range arrOfStrings {
				if strings.Contains(str, "window.__RELAY_STORE__") {
					JsonPart := strings.SplitN(str, "=", 2)[1]
					JsonPart = strings.TrimSpace(JsonPart)

					profId := getProfId(JsonPart, first, last)

					if profId != "" {
						parseRatings(profId)
					}
				}
			}
		}
	})
}

func main() {
	fmt.Println("hello world!")
	getProf("christian", "muise")
}
