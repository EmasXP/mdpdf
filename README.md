# mdpdf

This is just a personal project that I use to convert Markdown to PDF. It does that using Pandoc and a slightly customized LaTeX template. It generates high quality PDF files due to this approach, simply because Pandoc and LaTeX are awesome.

## Installing the fonts

If you don't already have EB Garamond 12 installed, go to this address https://github.com/octaviopardo/EBGaramond12/ and download the font files, preferably the OTF versions.

You need make sure that the fonts are accessible by the user that is running the service, and that's easiest by putting them in the `/usr/share/fonts` folder, for example:

```sh
sudo mkdir /usr/share/fonts/EBGaramond12
sudo mv EBGaramond-*.otf /usr/share/fonts/EBGaramond12/
```

## Starting the service

The only thing to remember is to have the correct working directory when starting the service. You need to start the service from the directory where `template.latex` is located. Otherwise `template.latex` will not be found by the service during run time. For example:

```sh
cd mdpdf
./mdpdf
```

## Firejail

The `pandoc` command is called with `--sandbox=true`, but to increase the sandbox even further, you can use `firejail`.

You will need to create a firejail profile, in this example I am going to name it `mdpdf.profile`:

```plain
whitelist /absolute/path/to/mdpdf/template.latex
```

This will work if you have put the `mdpdf` binary on your PATH (which is accessible by the root user.)

If you don't want to add `mdpdf` to your PATH, then you need to add this additional line to `mdpdf.profile`:

```plain
whitelist /absolute/path/to/mdpdf/mdpdf
```

You can do that like this if you want to add `mdpdf` to your PATH:

```sh
cp mdpdf /usr/bin/
```

(and here `mdpdf` is the binary file, not the folder)

And now you can start the service:

```sh
sudo firejail --profile=mdpdf.profile /absolute/path/to/mdpdf/mdpdf
```

If you have `mdpdf` on your PATH, you don't need to specify the absolute path, simply saying `mdpdf` is enough.

## Supervisor configuration example

```ini
[program:mdpdf]
command=firejail --profile=mdpdf.profile mdpdf
directory=/absolute/path/to/mdpdf/
user=root
```
