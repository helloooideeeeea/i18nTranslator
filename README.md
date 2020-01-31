# i18n Translator
### It helps you translate multiple languages.
very simple rule.

This library supports Accept-Language of Request Header.

And, locale file is key-value file format.

```
# config/message.en
HELLO=Hello %s.  
EMAIL_REQUIRED=Please fill email box.
```
```
# config/message.en-US
HELLO=Hey, %s!
EMAIL_REQUIRED=You need to fill email box.
```

Source Code Example.
```go
// exmaple/main.go

i18n, err := i18nTranslator.New("example/config", "en")
if err != nil {
    log.Fatal(err)
}

// example Request Header's Accept Language. 
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
acceptLanguage := "en-us;q=0.9,en;q=0.8,da;q=0.7"

// Translate
str , _:= i18n.Translate(acceptLanguage,"HELLO")

fmt.Println(fmt.Sprintf(str, "Jimmy")) // -> Hey, Jimmy!

// print debug (loaded key-value dump on memory)
i18n.PrintDebugLoadedDictionaries() // ↓
//
//　----- DEBUG -----
//  ----- lang [en] -----
//  [HELLO] : [Hello, %s.]
//  [EMAIL_REQUIRED] : [Please fill email box.]
//  ----- lang [en-us] -----
//  [HELLO] : [Hey, %s!]
//  [EMAIL_REQUIRED] : [You need to fill email box.]
//  ----- lang [ja] -----
//  [HELLO] : [こんにちわ。%sさん]
//  [EMAIL_REQUIRED] : [Emailは必須です。]

```

### rule 1.
Locale file naming rule is “hogehoge.local_identifier”
Examples are "form_err_message.en-US" or "form_err_message.en-us"

(The file extension can be either lowercase or uppercase, as long as it is a locale identifier. 
When loading files, use the same lower case letter.)

### rule 2.
The directory containing the locale files must be relative to the project directory.
Then, the directories under the directory are read recursively, and all files except for files without extension are recognized as locale files.  
Examples are "config/App1/message.en" or "config/App2/error.ja"

####You can perform test actions
$ git clone https://github.com/helloooideeeeea/i18nTranslator
$ cd i18nTranslator  
$ go run example/main.go  
This library is very small. So,

I wrote the specifications as detailed as possible in the source code.  
I want you to read.  
$ vim i18nTranslator.go

