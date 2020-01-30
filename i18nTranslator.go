package i18nTranslator

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type I18nTranslator struct {
	defaultLocale string
	dirPath	string
	dictionaries map[string]*map[string]string
}

// Constructor
func New(dirPath string, defaultLocale string) (*I18nTranslator, error) {

	t := &I18nTranslator{}
	t.defaultLocale = defaultLocale

	dirPath, _ = filepath.Abs(dirPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, err
	}
	t.dirPath = dirPath

	if err := t.loadFiles(); err != nil {
		return nil, err
	}

	return t, nil
}

func (translator *I18nTranslator) loadFiles() error {

	files, err := dirWalk(translator.dirPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Printf("[%s]: translate language files are not found.\nex. messages.en", translator.dirPath)
	}

	var dictionaries = make(map[string]*map[string]string)

	regex := regexp.MustCompile(`(.+?)=(.+)`)

	for _, f := range files {

		ext := filepath.Ext(f)[1:] // dot erase

		if ext == "" { // TODO I need to think a little more
			log.Printf("[%s]: Undefined File Extension. skip this file\n", f)
			continue
		}

		fp, err := os.Open(f)
		if err != nil {
			continue
		}
		defer fp.Close()

		m := map[string]string{}

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			text := scanner.Text()
			pair := regex.FindStringSubmatch(text)
			if len(pair) != 3 {
				log.Printf("[%s]:[%s] skip this line.\n", f, text)
				continue
			}
			m[pair[1]] = pair[2]
		}

		dictionaries[ext] = &m
	}

	translator.dictionaries = dictionaries
	return nil
}

func (translator *I18nTranslator) Translate(lang string, key string) (string, bool) {

	parsedLang := translator.parse(lang)
	if parsedLang == "" {
		return translator.TranslateByDefaultLocale(key)
	}

	dic := *translator.dictionaries[parsedLang]
	key, ok := dic[key]
	return key, ok
}

func (translator *I18nTranslator) TranslateByDefaultLocale(key string) (string, bool) {
	dic := *translator.dictionaries[translator.defaultLocale]
	key, ok := dic[key]
	return key, ok
}

func (translator *I18nTranslator) parse(httpAcceptLanguage string) string {
	// * or ""
	if httpAcceptLanguage == "*" || httpAcceptLanguage == "" {
		return ""
	}

	// ja,en-US;q=0.9,en;q=0.8,da;q=0.7
	hals := strings.Split(httpAcceptLanguage,",")
	for _, alang := range hals {
		if strings.Contains(alang, ";") {
			alang = strings.Split(alang,";")[0]
		}

		for llang ,_ := range translator.dictionaries {
			if alang == llang {
				return alang
			}
		}
	}
	return ""
}

func (translator *I18nTranslator) PrintDebugLoadedDictionaries() {
	for lang, dict := range translator.dictionaries {
		fmt.Println("----- DEBUG -----")
		fmt.Printf("----- lang [%s] -----\n", lang)
		for key, value := range *dict {
			fmt.Printf("[%s] : [%s]\n", key, value)
		}
	}
}

func dirWalk(dir string) ([]string, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			a , err := dirWalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, a...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths, nil
}



