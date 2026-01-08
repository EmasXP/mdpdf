# mdpdf

This is just a personal project that I use to convert Markdown to PDF. It does that using Pandoc and the EB Garamond 12, JetBrains Mono NL and Inter fonts. It generates high quality PDF files due to this approach, simply because Pandoc and LaTeX are awesome.

## Installing the fonts

If you don't already have the fonts installed, download them from these addresses:

* EB Garamond 12: https://github.com/octaviopardo/EBGaramond12/ (preferably the OTF versions) 
* JetBrains Mono NL: https://www.jetbrains.com/lp/mono/
* Inter: https://rsms.me/inter/

You need to make sure that the fonts are accessible by the user that is running the service, and that's easiest by putting them in the `/usr/share/fonts` folder, for example:

```sh
sudo mkdir /usr/share/fonts/EBGaramond12
sudo mv EBGaramond-*.otf /usr/share/fonts/EBGaramond12/
sudo mkdir /usr/share/fonts/JetBrainsMonoNL
sudo mv JetBrainsMonoNL-*.ttf /usr/share/fonts/JetBrainsMonoNL/
sudo mkdir /usr/share/fonts/Inter
sudo mv Inter-*.otf /usr/share/fonts/Inter/
```

## Starting the service

It's rather easy:

```sh
cd mdpdf
./mdpdf
```

## Firejail

The `pandoc` command is called with `--sandbox=true`, but to increase the sandbox even further, you can use `firejail`.

```sh
sudo firejail mdpdf
```

This is all you need do if you have put the `mdpdf` binary on your PATH (which is accessible by the root user.)

If you don't want to add `mdpdf` to your PATH, then you need to add this additional line to `mdpdf.profile`:

```plain
whitelist /absolute/path/to/mdpdf/mdpdf
```

You can do like this if you want to add `mdpdf` to your PATH:

```sh
cp mdpdf /usr/bin/
```

(and here `mdpdf` is the binary file, not the folder)

And now you can start the service:

```sh
sudo firejail --profile=mdpdf.profile /absolute/path/to/mdpdf/mdpdf
```

## Supervisor configuration example

```ini
[program:mdpdf]
command=firejail --profile=mdpdf.profile mdpdf
user=root
```

## Using as a CLI tool

```
Usage: ./mdpdf [options] <input> [output]
  -template string
        Path or name of the LaTeX template file (default "default")
```

Example:

```sh
mdpdf my-document.md
```

This will render a `my-document.pdf`.

To specify the output file:

```sh
mdpdf my-document.md output-file.pdf
```

## About a future sans-serif version

It's not very likely that this is ever going to happen. The typeface picks are _very_ opinionated, but it does not hurt to at least think about it.

The currently used sans-serif typeface is Inter, which render itself beautifully on headings and titles. In my personal opinion, it does not render itself as pretty on long text. I think we're better off with Carlito (or possibly Lato) in that case. Combining Inter (headings, titles) with Carlito or Lato (paragraphs) does not sound exciting to me, though.
