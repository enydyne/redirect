package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var (
	Redirects     = map[string]string{}
	RedirectsHtml []byte
	Addr          = ":4321"
)

func init() {
	file, err := os.Open("urls.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if port, ok := os.LookupEnv("PORT"); ok {
		Addr = fmt.Sprintf(":%s", port)
	}

	urlsHtml := strings.Builder{}
	urlsHtml.WriteString(`<!doctype html><html><head><meta charset="UTF-8"><meta content="width=device-width, initial-scale=1" name="viewport"><title>Redirect</title></head><body><table><tbody><tr><th>ID</th><th>Target URL</th></tr>`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid line: %s", line))
		}
		url := strings.TrimSpace(parts[1])
		id := strings.TrimSpace(parts[0])
		if id == "" || url == "" {
			panic(fmt.Errorf("invalid line format: %s", line))
		}
		Redirects[id] = url
		urlsHtml.WriteString(fmt.Sprintf("<tr><td>%s</td><td><a href=%q><code>%s</code></a></td></tr>", id, url, url))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	urlsHtml.WriteString("</tbody></table></body></html>")

	RedirectsHtml = []byte(urlsHtml.String())
}

func main() {
	if len(os.Args) == 2 {
		url := os.Args[1]
		length := 10
		chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		var id string
		for {
			var randomString []rune
			for i := 0; i < length; i++ {
				randomString = append(randomString, chars[rand.Intn(len(chars))])
			}
			id = string(randomString)
			if _, ok := Redirects[id]; !ok {
				Redirects[id] = url
				break
			}
		}

		file, err := os.OpenFile("urls.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("%s=%s\n", id, url)); err != nil {
			panic(err)
		}

	} else {
		http.HandleFunc("/", handle)
		if err := http.ListenAndServe(Addr, nil); err != nil {
			panic(err)
		}
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	if dest, ok := Redirects[r.URL.Path[1:]]; ok {
		http.Redirect(w, r, dest, http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(RedirectsHtml)
}
