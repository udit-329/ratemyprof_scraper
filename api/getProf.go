package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func parseRatings(profId string) map[string]interface{} {
	url := "https://www.ratemyprofessors.com/ShowRatings.jsp?tid=" + profId

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)

	profRatings := make(map[string]interface{})
	allRatings := make([]map[string]interface{}, 0)

	document.Find("script").Each(func(i int, s *goquery.Selection) {

		text := s.Text()

		// Get the script tag we want
		if strings.Contains(text, "window.__RELAY_STORE__") {

			//Split based on semicolons AND } because we can have semicolons in comments
			arrOfStrings := strings.Split(text, "};")

			//Get the string we want
			for _, str := range arrOfStrings {
				if strings.Contains(str, "window.__RELAY_STORE__") {
					JsonPart := strings.SplitN(str, "=", 2)[1]
					JsonPart = strings.TrimSpace(JsonPart)
					//Add } since it got remove when we split the string
					JsonPart = JsonPart + "}"

					var result map[string]interface{}
					json.Unmarshal([]byte(JsonPart), &result)

					for _, element := range result {
						// Only check maps
						if _, ok := element.(map[string]interface{}); ok {

							//Parse Teacher
							if element.(map[string]interface{})["__typename"] == "Teacher" {

								profRatings["wouldTakeAgainPercent"] = element.(map[string]interface{})["wouldTakeAgainPercent"].(float64)
								profRatings["avgDifficulty"] = element.(map[string]interface{})["avgDifficulty"].(float64)
								profRatings["numRatings"] = int(element.(map[string]interface{})["numRatings"].(float64))
								profRatings["avgRating"] = element.(map[string]interface{})["avgRating"].(float64)

								//Parse Ratings
							} else if element.(map[string]interface{})["__typename"] == "Rating" {
								rating := make(map[string]interface{})

								rating["class"] = element.(map[string]interface{})["class"].(string)
								rating["grade"] = element.(map[string]interface{})["grade"].(string)
								rating["comment"] = element.(map[string]interface{})["comment"].(string)
								rating["difficultyRating"] = element.(map[string]interface{})["difficultyRating"].(float64)
								//From what I understand, the quality is the average of clarity and helpful
								rating["clarityRating"] = element.(map[string]interface{})["clarityRating"].(float64)
								rating["helpfulRating"] = element.(map[string]interface{})["helpfulRating"].(float64)

								allRatings = append(allRatings, rating)
							}
						}
					}
				}
			}
		}
	})
	profRatings["ratings"] = allRatings

	return profRatings
}

func getProfId(jsonString string, first string, last string) string {

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

func getProf(first string, last string, schoolCode string) (profRatings map[string]interface{}) {

	url := "https://www.ratemyprofessors.com/search/teachers?query=" + first + "%20" + last + "&sid=" + schoolCode + "="

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{"error": err}
	}

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{"error": err}
	}

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
						profRatings = parseRatings(profId)
					}
				}
			}
		}
	})
	return profRatings
}

func start(c *gin.Context) {

	first := c.Query("first")
	last := c.Query("last")
	schoolCode := c.Query("school")

	c.JSON(200, getProf(first, last, schoolCode))
}

func Main(w http.ResponseWriter, r *http.Request) {

	// rg := gin.Default()
	// rg.GET("/getProf", start)
	// rg.Run()
	profDetails := getProf("yuan", "tian", "U2Nob29sLTE0NjY")
	fmt.Println(profDetails)
	fmt.Fprintf(w, "hi")
}
