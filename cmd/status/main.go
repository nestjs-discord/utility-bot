package main

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/cmd/status/oja"
	"log"
	"os"
	"sort"
	"text/template"
)

func main() {
	tmpl, err := template.ParseFiles("init.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	dataProviders := map[string]func() []string{
		"../../bot/status/texts_oja.go": oja.Fetch,
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
