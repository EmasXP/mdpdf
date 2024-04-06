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
			<textarea name="md" style="width:100%; height: 400px">---
papersize: a4
pagestyle: empty
fontsize: 12pt
linestretch: 1.2
geometry: margin=6.25em
large-h1: true
---

</textarea>
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
			More technical stuff: The LaTeX engine used is <a href="https://tug.org/xetex/">XeTeX</a>, the site itself is built using <a href="https://go.dev/">Go</a>.
		</p>

		<p>
			The source code can be found here: <a href="https://github.com/EmasXP/mdpdf">https://github.com/EmasXP/mdpdf</a>
		</p>

		<h2>The meta block</h2>

		<p>
		One can add a meta data block in the beginning of the Markdown source, like this:
		</p>

		<pre><code>---
Meta data goes here
---
The document starts here.</code></pre>

		<p>
			I've added a predefined meta data block that I believe works well for Markdown files (at least, this is how I personally think the Markdown files should look when printed on paper.)
		</p>

		<p>
			As an example, one can add a title like this:
		</p>

		<pre><code>---
title: My document
---</code></pre>

		<p>
			Vertical space is added between paragraphs by default. That can be changed to use indentation instead:
		</p>

		<pre><code>---
title: My document
indent: true
---</code></pre>

		<p>
			I've kept <code>title</code> in this example just to show you how it looks when several variables are defined.
		</p>

		<p>
			To specify the language of the document, use the <code>lang</code> variable:
		</p>

		<pre><code>---
lang: en-GB
---</code></pre>

		<p>
			More variables are found here: <a href="https://pandoc.org/chunkedhtml-demo/6.2-variables.html">https://pandoc.org/chunkedhtml-demo/6.2-variables.html</a>.
		</p>

		<p>
			The <code>mainfont</code> is hard-coded and cannot be changed.
		</p>

		<h3>Additional variables</h3>

		<p>
			Two additional variables has been added:
		</p>

		<h4>large-h1</h4>

		<p>
			Since Markdown most often does not have a title variable, but uses the H1 (<code>#</code>) as such, the <code>large-h1</code> can be used to simply make the H1 larger:
		</p>

		<pre><code>---
large-h1: true
---</code></pre>

		<h4>additional</h4>

		<p>
			Another way of achieving a larger H1 is to use the <code>additional</code> variable:
		</p>

		<pre><code>---
additional: \usepackage{sectsty}
  \sectionfont{\fontsize{21}{21}\selectfont}
---</code></pre>

	</body>
</html>
`

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, INDEX)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		command := exec.Command(
			"pandoc",
			"--sandbox=true",
			"-f",
			"gfm",
			"--to",
			"pdf",
			"--template",
			"template.latex",
			"-V",
			"mainfont=EB Garamond",
			// Trying out sans-serif fonts for a future sans-serif version:
			//"mainfont=Inter",
			//"mainfont=Carlito",
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
