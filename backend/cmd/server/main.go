package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Отдаём простую HTML-страницу с русским текстом
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Coobe — приветствие</title>
    <style>
      body { margin:0; height:100vh; display:flex; align-items:center; justify-content:center;
             font-family: system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
             background: #000; color: #fff; }
      h1 { font-size: 26px; font-weight: 500; }
    </style>
  </head>
  <body>
    <h1>Привет, этот сайт написан на Go.</h1>
  </body>
</html>`)
	})

	log.Println("➡️  Server listening on :8080")
	// Если порт уже занят, тут будет ошибка — будет напечатана в логе
	log.Fatal(http.ListenAndServe(":8080", nil))
}