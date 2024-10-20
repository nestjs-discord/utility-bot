package main

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/cmd/status/ja"
	"github.com/nestjs-discord/utility-bot/cmd/status/oja"
	"github.com/nestjs-discord/utility-bot/cmd/status/pqa"
	"log"
	"os"
	"sort"
	"text/template"
)

func main() {
	for i := 0; i < 100; i++ { // rate limit is 100 requests per hour.
		err := pqa.Scrape()
		if err != nil {
			fmt.Printf("pqa scrape error: %s\n", err)
			break
		}
	}

	tmpl, err := template.ParseFiles("init.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	dataProviders := map[string]func() []string{
		"../../bot/status/texts_ja.go":  ja.Fetch,
		"../../bot/status/texts_oja.go": oja.Fetch,
		"../../bot/status/texts_pqa.go": pqa.Fetch,
	}

	for fileName, dataProvider := range dataProviders {
		texts := dataProvider()
		sort.Strings(texts)                 // ascending order
		outFile, err := os.Create(fileName) // Open or create the output .go file
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(outFile, texts)
		_ = outFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("wrote", fileName)
	}

	fmt.Println("done")
}
