package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	chapterDirNamePattern = regexp.MustCompile("^Chapter\\d*_")
	htmlFileNamePatten    = regexp.MustCompile(".html$")

	welcomeTemplate = template.Must(template.New("welcome").Parse(welcomeTemplateStr))
)

type HtmlPage struct {
	Name    string
	Path    string
	Content string
}

func loadHtmlPages() (pages []HtmlPage) {
	fileInfoUnderWorkingDir, err := ioutil.ReadDir(".")
	if err == nil {
		for _, fileInfo := range fileInfoUnderWorkingDir {
			if fileInfo.IsDir() && chapterDirNamePattern.MatchString(fileInfo.Name()) {
				log.Println("found chapter dir: " + fileInfo.Name())

				for _, htmlFile := range listHtmlFiles(fileInfo) {
					path := fileInfo.Name() + "/" + htmlFile.Name()
					log.Println("loading " + path)

					content, err := ioutil.ReadFile(path)
					if err != nil {
						log.Printf("error when reading html file %v, caused by\n%v\n", path, err)
						continue
					}
					pages = append(pages, HtmlPage{htmlFile.Name(), path, string(content)})
				}
			}
		}
	}
	return
}

func listHtmlFiles(dir os.FileInfo) (htmlFiles []os.FileInfo) {
	fileInfoUnderDir, err := ioutil.ReadDir(dir.Name())
	if err == nil {
		log.Println("go through chapter dir: " + dir.Name())
		for _, fileInfo := range fileInfoUnderDir {
			if !fileInfo.IsDir() && htmlFileNamePatten.MatchString(fileInfo.Name()) {
				log.Println("found html file: " + fileInfo.Name())
				htmlFiles = append(htmlFiles, fileInfo)
			}
		}
	} else {
		log.Printf("error occured when get file info under %v, caused by\n%v\n", dir.Name(), err)
	}
	return
}

func main() {
	pages := loadHtmlPages()
	pagePathStr := ""
	for _, page := range pages {
		pagePathStr = pagePathStr + page.Path + "\n"
	}
	log.Printf("following html5 pages have been loaded\n%s\n", pagePathStr)

	findPage := func(name string) (HtmlPage, bool) {
		for _, page := range pages {
			if name == page.Name {
				return page, true
			}
		}
		return HtmlPage{}, false
	}

	http.HandleFunc("/html5", func(w http.ResponseWriter, r *http.Request) {
		pageParam := r.FormValue("page")
		if pageParam == "" {
			welcomeTemplate.Execute(w, pages)
		} else {
			page, found := findPage(pageParam)
			if found {
				io.WriteString(w, page.Content)
			} else {
				io.WriteString(w, "Page not found... :(")
			}
		}
	})

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("Failed to startup the server", err)
	}
}

const welcomeTemplateStr = `
<!DOCTYPE HTML>
<html>
	<head>
	    <title>Practicing HTML5</title>
	</head>
	<body>
		<h1>Welcome to the practicing server for HTML5!</h1>
	    
	    <ul>
	    {{range .}}
	    <li><a href="/html5?page={{.Name}}">{{.Name}}</a></li>
	    {{end}}
	    </ul>
	</body>
</html>
`
