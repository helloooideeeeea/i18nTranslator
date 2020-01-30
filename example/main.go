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

	acceptLanguage := "en-US;q=0.9,en;q=0.8,da;q=0.7"

	str , _:= i18n.Translate(acceptLanguage,"HELLO")

	fmt.Println(fmt.Sprintf(str, "Jimmy"))

	i18n.PrintDebugLoadedDictionaries()
}
