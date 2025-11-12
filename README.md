# mdpdf

This is just a personal project that I use to convert Markdown to PDF. It does that using Pandoc and the EB Garamond 12 font. It generates high quality PDF files due to this approach, simply because Pandoc and LaTeX are awesome.

## Installing the fonts

If you don't already have EB Garamond 12 installed, go to this address https://github.com/octaviopardo/EBGaramond12/ and download the font files, preferably the OTF versions. The same thing goes for the JetBrains Mono NL, which can be downloaded here: https://www.jetbrains.com/lp/mono/

You need to make sure that the fonts are accessible by the user that is running the service, and that's easiest by putting them in the `/usr/share/fonts` folder, for example:

```sh
sudo mkdir /usr/share/fonts/EBGaramond12
sudo mv EBGaramond-*.otf /usr/share/fonts/EBGaramond12/
sudo mkdir /usr/share/fonts/JetBrainsMonoNL
sudo mv JetBrainsMonoNL-*.ttf /usr/share/fonts/JetBrainsMonoNL/
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

