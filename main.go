package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var Redirects = map[string]string{}

func init() {
	file, err := os.Open("urls.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid line: %s", line))
		}
		Redirects[parts[0]] = parts[1]
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
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
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	if dest, ok := Redirects[r.URL.Path[1:]]; ok {
		http.Redirect(w, r, dest, http.StatusFound)
		return
	}

	http.NotFound(w, r)
}
