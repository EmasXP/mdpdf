package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const INDEX = `
<!doctype html>
<html lang="en">
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<title>MD</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
	<link href="https://fonts.googleapis.com/css2?family=EB+Garamond:ital,wght@0,400..800;1,400..800&family=Inconsolata:wght@200..900&display=swap" rel="stylesheet">
	<body>

		<form method="post">
			<textarea name="md" style="width:100%; height: 400px">---
papersize: a4
pagestyle: empty
fontsize: 12pt
linestretch: 1.2
geometry: margin=6.25em
header-includes:
  - \pagenumbering{gobble}
  - \usepackage{titling}
  - \setlength{\droptitle}{-3.88em}
  - \pretitle{\fontsize{21}{21}\selectfont\bfseries}
  - \posttitle{\par\vspace{-3.1em}}
---

</textarea>
			<button>Convert Markdown to PDF</button>
		</form>

		<div id="information">
			<h1>What's this?</h1>

			<p>
				This is just a web page to convert <a href="https://en.wikipedia.org/wiki/Markdown">Markdown</a> into PDF. It does that using <a href="https://pandoc.org/">Pandoc</a> and the <a href="https://github.com/octaviopardo/EBGaramond12/">EB Garamond 12</a> font. Pandoc and LaTeX generates super crisp PDF files. Before I used this approach I converted the Markdown to HTML, and used a web browser to "print" the page to a PDF file. For some reason the fonts are rendered more blurry when taking that route.
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

			<h3>Defining the title of the document</h3>

			<p>
				That is done using the <code>title</code> variable:
			</p>

			<pre><code>---
title: My document
---</code></pre>

			<h4>Markdown title</h4>

			<p>
				Markdown files commonly use the H1 (<code># Title</code>) for the document title. If that is the case, my suggestion is to use the <code>title</code> variable in the meta block instead. That way the title will be typeset in a larger font size.
			</p>

			<p>
				A <code>title</code> is rendered in the center of the page, and has large top and bottom margins. Markdown files are commonly <i>not</i> rendered this way. If you want to make the <code>title</code> look more like a regular H1, add this to the meta block:
			</p>

			<pre><code>---
header-includes:
- \usepackage{titling}
- \setlength{\droptitle}{-3.88em}
- \pretitle{\fontsize{21}{21}\selectfont\bfseries}
- \posttitle{\par\vspace{-3.1em}}
---</code></pre>

			<p>
				This example will make the title more distinct than an ordinary H1, which is suitable for a document title.
			</p>

			<h3>Paragraph spacing</h3>

			<p>
				Vertical space is added between paragraphs by default. That can be changed to use indentation instead:
			</p>

			<pre><code>---
indent: true
---</code></pre>

			<h3>Language of the document</h3>

			<p>
				To specify the language of the document, use the <code>lang</code> variable:
			</p>

			<pre><code>---
lang: en-GB
---</code></pre>

			<h3>Removal of page numbers</h3>

			<p>
				When using <code>pagestyle: empty</code>, page numbers are removed. However, the first page will still have a page number. To remove that too, one can add the following:
			</p>

			<pre><code>---
header-includes:
- \pagenumbering{gobble}
---</code></pre>

			<h3>More variables</h3>

			<p>
				More variables are found here: <a href="https://pandoc.org/chunkedhtml-demo/6.2-variables.html">https://pandoc.org/chunkedhtml-demo/6.2-variables.html</a>.
			</p>

			<p>
				The <code>mainfont</code> is hard-coded and cannot be changed.
			</p>
		</div>

		<style>
		<!--
		body {
			font-family: 'EB Garamond', serif;
		}
		#information {
			max-width: 60em;
			margin: auto;
			padding: 1em;
			margin-bottom: 4em;
		}
		code, textarea {
			font-family: 'Inconsolata', monospace;
			font-size: 0.85em;
		}
		pre {
			background: #f8f8f8;
			padding: 0.6em 1em;
			overflow-x: auto;
			border-radius: 0.5em;
			border: 1px solid #eeeeee;
		}
		h2, h3, h4 {
			margin-top: 2em;
		}
		button {
			padding: 0.5em 1em;
		}
		-->
		</style>

	</body>
</html>
`

func main() {
	// Define the optional flag
	templatePath := flag.String("template", "default", "Path or name of the LaTeX template file")

	// Custom usage function to show usage info
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <input> [output]\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the flags
	flag.Parse()

	// Get positional arguments
	args := flag.Args()
	if len(args) < 1 {
		web(*templatePath)
	}

	input := args[0]

	var output string
	if len(args) >= 2 {
		output = args[1]
	} else {
		// Derive output from input
		ext := filepath.Ext(input)
		if ext == ".md" {
			output = strings.TrimSuffix(input, ".md") + ".pdf"
		} else {
			output = input + ".pdf"
		}
	}

	cli(input, output, *templatePath)
}

func cli(input, output, templatePath string) {
	fmt.Println(templatePath)
	command := exec.Command(
		"pandoc",
		input,
		"--sandbox=false",
		"-o",
		output,
		"-f",
		"markdown",
		"--to",
		"pdf",
		"--template",
		templatePath,
		"-V",
		"mainfont=EB Garamond",
		// Trying out sans-serif fonts for a future sans-serif version:
		//"mainfont=Inter",
		//"mainfont=Carlito",
		"--pdf-engine",
		"xelatex",
	)

	var stderr bytes.Buffer
	command.Stderr = &stderr

	_, err := command.Output()
	if err != nil {
		error, _ := io.ReadAll(&stderr)
		log.Println(string(error))
		return
	}
}

func web(templatePath string) {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, INDEX)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		command := exec.Command(
			"pandoc",
			"--sandbox=true",
			"-f",
			"markdown",
			"--to",
			"pdf",
			"--template",
			templatePath,
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
