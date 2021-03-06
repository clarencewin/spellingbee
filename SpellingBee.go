package main

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type WordInfo struct {
	textOnly  string
	Word      string `json:"word"`
	IsPangram bool   `json:"pangram"`
}

func main() {
	words := dictionaryList()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.StaticFS("ui", http.Dir("static"))

	router.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})

	router.GET("/rest/words", func(c *gin.Context) {
		mustContain := c.Query("mustContain")
		canContain := c.Query("canContain")
		//
		c.JSON(http.StatusOK, wordsMatching(
			words,
			mustContain,
			canContain,
		))
	})

	router.GET("/words", func(c *gin.Context) {
		mustContain := c.Query("mustContain")
		canContain := c.Query("canContain")

		wordsToDisplay := wordsMatching(
			words,
			mustContain,
			canContain)

		c.HTML(http.StatusOK, "wordsList.tmpl", gin.H{"values": wordsToDisplay})
	})

	router.Run(":8080")
}

//Gets dictionary list
func dictionaryList() []string {

	words := make([]string, 0)

	f, err := os.Open("words.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		nextWord := scanner.Text()
		words = append(words, nextWord)
	}

	return words
}

//Gets words containing
func wordsMatching(wordsList []string,
	required string,
	canContain string) []WordInfo {

	canContainMap := map[string]bool{required: true}
	canContainList := strings.Split(canContain, "")

	for allowedLetter := range canContainList {
		canContainMap[canContainList[allowedLetter]] = true
	}

	wordsCont := make([]WordInfo, 0, len(wordsList))

	for _, word := range wordsList {

		reg, err := regexp.Compile("[^a-z]+")

		if err != nil {
			log.Fatal(err)
		}

		wordFormatted := reg.ReplaceAllString(strings.ToLower(word), "")

		if len(wordFormatted) < 4 {
			continue
		}

		if strings.Contains(wordFormatted, required) {

			letters := strings.Split(strings.ToLower(wordFormatted), "")

			addToList := true

			lettersUsed := make(map[string]bool)

			for _, letter := range letters {

				if canContainMap[letter] == true {
					lettersUsed[letter] = true
				} else {
					addToList = false
					break
				}

			}

			isPangram := len(canContainMap) == len(lettersUsed)

			if addToList {
				wordsCont = append(wordsCont, WordInfo{wordFormatted, word, isPangram})
			}
		}
	}

	return wordsCont
}
