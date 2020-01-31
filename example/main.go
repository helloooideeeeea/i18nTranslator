package main

import (
	"fmt"
	"i18nTranslator"
	"log"
)

func main() {

	i18n, err := i18nTranslator.New("example/config", "en")
	if err != nil {
		log.Fatal(err)
	}

	// example Request Header's Accept Language.
	acceptLanguage := "en-us;q=0.9,en;q=0.8,da;q=0.7"

	// Translate
	str , _:= i18n.Translate(acceptLanguage,"HELLO")

	fmt.Println(fmt.Sprintf(str, "Jimmy"))

	// print debug (loaded key-value dump on memory)
	i18n.PrintDebugLoadedDictionaries()
}
