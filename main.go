package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

const INDEX = `
<!doctype html>
<html lang="en">
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<title>MD</title>
	<body>

		<form method="post">
			<textarea name="md" style="width:100%; height: 400px"></textarea>
			<button>To PDF</button>
		</form>

		<h2>What's this?</h2>

		<p>
			This is just a web page to convert <a href="https://en.wikipedia.org/wiki/Markdown">Markdown</a> into PDF. It does that using <a href="https://pandoc.org/">Pandoc</a> and a customized <a href="https://en.wikipedia.org/wiki/LaTeX">LaTeX</a> template. The most notable change is the use of the <a href="https://github.com/octaviopardo/EBGaramond12/">EB Garamond 12</a> font. Pandoc and LaTeX generates super crisp PDF files. Before I used this approach I converted the Markdown to HTML, and used a web browser to "print" the page to a PDF file. For some reason the fonts are rendered more blurry when taking that route.
		</p>

		<p>
			This web page is hosted on an old Raspberry Pi, so it is rather slow. Even though this is intended for my personal use, you are free to use it too.
		</p>

		<p>
			More technical stuff: The LaTeX engine used is <a href="https://tug.org/xetex/">XeTeX</a>, the site itself is built using <a href="https://go.dev/">Go</a>. It's amazing how much tech is needed for such a small web site.
		</p>

		<p>
			The source code can be found here: <a href="https://github.com/EmasXP/mdpdf">https://github.com/EmasXP/mdpdf</a>
		</p>
	</body>
</html>
`

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
		io.WriteString(w, INDEX)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request){
		command := exec.Command(
			"pandoc",
            "-f",
            "gfm",
            "--to",
            "pdf",
            "-V",
            "papersize:a4",
            "-V",
            "pagestyle=empty",
            "-V",
            "linestretch=1.2",
            "--template",
            "template.latex",
            "-V",
            "fontsize:12pt",
            "-V",
            "documentclass=scrartcl",
            "-V",
            "mainfont=EB Garamond",
            "--pdf-engine",
            "xelatex",
		)

		md := r.PostFormValue("md")
		command.Stdin = strings.NewReader(md)

		var stderr bytes.Buffer
		command.Stderr = &stderr

		stdout, err := command.Output()
		if err != nil {
			error, _ := io.ReadAll(&stderr)
			w.Header().Set("Content-Type", "text/plain")
			w.Write(error)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Write(stdout)
	})

	fmt.Println(":9472")
	http.ListenAndServe(":9472", nil)
}
