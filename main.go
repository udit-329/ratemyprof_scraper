package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func getProf(first string, last string) {
	// url := "https://www.ratemyprofessors.com/search/teachers?query=" + first + "%20" + last + "&sid=U2Nob29sLTE0NjY="
	url := "https://www.ratemyprofessors.com/search/teachers?query=christian&sid=U2Nob29sLTE0NjY="
	// url := "https://www.amazon.ca/Mattel-GAMES-W2085-Card-Game/dp/B00CTH0A1Q"
	response, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		fmt.Println(err)
	}

	n, err := io.Copy(os.Stdout, response.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Number of bytes copied to STDOUT:", n)

	document, err := goquery.NewDocumentFromReader(response.Body)

	document.Find("a").Each(func(i int, s *goquery.Selection) {
		// as, _ := s.Find("a").Attr("href")
		as, _ := s.Attr("href")
		fmt.Println(as)
		fmt.Println("f")
	})
	// link := document.Find("#productTitle")
	// fmt.Println(link)
	// fmt.Println(a)
	fmt.Println("fs")
}

func main() {
	fmt.Println("hello world!")
	getProf("yuan", "tian")
}
