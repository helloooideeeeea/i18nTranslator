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
	defaultLocale string 						// ex. "en-US"
	dirPath	string 								// ex. "$GOPATH/src/i18nTranslator/example/config"
	dictionaries map[string]*map[string]string 	// ex. ["en-US"]["Hello"]["hey, %s"]
}

/*
Constructor

The first argument is the directory where the locale files are stored, expecting a relative path from the project directory.
(第一引数は、ロケールファイルが格納されているディレクトリで、プロジェクトディレクトリからの相対パスが渡される事を期待しています。)

(The second argument expects a default locale identifier to be passed.)
第二引数は、デフォルトのロケール識別子が渡される事を期待しています。
*/
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

/*
Loads a list of files with a certain file extension (such as .en) under the specified directory.
(Recursively search the specified directory.)

指定したディレクトリ配下にある、"何かしらのファイル拡張子(.enなど)が付いている"ファイル一覧を読み込みます。
(指定したディレクトリを再帰的に探索します。)
*/
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

		ext := filepath.Ext(f)[1:] // extension's dot erase

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

		dictionaries[strings.ToLower(ext)] = &m
	}

	translator.dictionaries = dictionaries
	return nil
}

/*
Get the value from the key for each locale.
The first argument assumes that "Accept-Language" of Request Header is passed.

https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language

ロケール別のキーからバリューを取得します。
第一引数は、Request Header の "Accept-Language" が渡される事を想定しています。
*/
func (translator *I18nTranslator) Translate(lang string, key string) (string, bool) {

	parsedLang := translator.parse(lang)
	if parsedLang == "" {
		return translator.TranslateByDefaultLocale(key)
	}

	dic := *translator.dictionaries[parsedLang]
	key, ok := dic[key]
	return key, ok
}

/*
Get value from key with default locale identifier.
デフォルトのロケール識別子で、キーからバリューを取得します。
*/
func (translator *I18nTranslator) TranslateByDefaultLocale(key string) (string, bool) {
	dic := *translator.dictionaries[translator.defaultLocale]
	key, ok := dic[key]
	return key, ok
}

/*
Parse "Accept-Language" of RequestHeader of the first argument.

	ja,en-US;q=0.9,en;q=0.8,da;q=0.7

In the above case, the locale identifier of the locale file loaded in the memory is searched from the left.
If there is no locale identifier loaded in memory, it will try to return the value of the default locale identifier.
The locale identifier is unified with the lower case letter when loading it into memory, so you do not have to worry about uppercase and lowercase letters.

第一引数のRequestHeaderの"Accept-Language"をパースします。

	ja,en-US;q=0.9,en;q=0.8,da;q=0.7

上記の様な場合、メモリにロードしているロケールファイルのロケール識別子を、左から検索します。
メモリにロードしているロケール識別子が存在しない場合、デフォルトのロケール識別子のValueを返却しようとします。
ロケール識別子は、メモリに読み込む際に、Lower Case Letterで統一していますので、大文字、小文字で悩む必要はないです。
*/
func (translator *I18nTranslator) parse(httpAcceptLanguage string) string {
	// * or ""
	if httpAcceptLanguage == "*" || httpAcceptLanguage == "" {
		return ""
	}

	// ja,en-US;q=0.9,en;q=0.8,da;q=0.7
	hals := strings.Split(strings.ToLower(httpAcceptLanguage),",")
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

/*
Outputs the key value loaded in memory for each locale. For debugging.
メモリにロードしているキーバリューをロケール別に出力します。デバッグ用です。
*/
func (translator *I18nTranslator) PrintDebugLoadedDictionaries() {
	for lang, dict := range translator.dictionaries {
		fmt.Printf("----- lang [%s] -----\n", lang)
		for key, value := range *dict {
			fmt.Printf("[%s] : [%s]\n", key, value)
		}
	}
}

/*
Recursively obtains the file path from the directory path of the first argument.
第一引数のディレクトリパスから、再帰的にファイルパスを取得します。
*/
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



