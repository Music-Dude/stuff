package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var outdir string

func saveImage(filename string) {
    path := outdir+"/"+filename

    if _, err := os.Stat(path); os.IsNotExist(err) {
        fmt.Println("Creating "+path)

        resp, err := http.Get("https://i.redd.it/" + filename)
        if err != nil {
            log.Fatalln(err)
        }
        defer resp.Body.Close()

        file, err := os.Create(path)
        if err != nil {
            log.Fatalln(err)
        }
        defer file.Close()
        file.ReadFrom(resp.Body)
    } else {
        fmt.Println(path+" already exists, skipping")
    }
}

func main() {

	flag.StringVar(&outdir, "out", "imgs", "Output directory to save images in")

	flag.Parse()

	// subs arg is required
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

    os.Mkdir(outdir, 0777)

	subs := strings.Join(flag.Args(), "+")
	client := http.Client{}

	var lastpost string
	urlre := regexp.MustCompile(`"url": "https?://i.redd.it/(.+?)"`)
	idre := regexp.MustCompile(`"after": "(\w+)"`)

	for {
		//var data map[string]interface{}

		req, err := http.NewRequest("GET", "https://reddit.com/r/"+subs+".json?limit=100&after="+string(lastpost), nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("User-Agent", "The Scraperer/0.1")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
        defer resp.Body.Close()

		substrs := urlre.FindAllStringSubmatch(string(content), -1)
		for _, str := range substrs {
			saveImage(str[1])
		}

		after := idre.FindStringSubmatch(string(content))
		if len(after) < 2 {
			continue
        }
		lastpost = after[1]

		//json.Unmarshal(content, &data)
		//fmt.Println(string(content))
		//posts := data["data"].(map[string]interface{})["children"].([]interface{})

		//for _, post := range posts {
		//	//fmt.Printf("%v\n\n\n", post["data"].(map[string]string)["url"])
		//	fmt.Println(post.(map[string]interface{})["url"].(string), "\n\n\n\n")
		//}
	}
}
