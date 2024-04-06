# mdpdf

This is just a personal project that I use to convert Markdown to PDF. It does that using Pandoc and a slightly customized LaTeX template. It generates high quality PDF files due to this approach, simply because Pandoc and LaTeX are awesome.

## Starting the service

The only thing to remember is to have the correct working directory whens tarting the service. You need to start the service from the directory where `template.latex` is located. Otherwise `template.latex` will not be found by the service during run time. For example:

```shell
cd mdpdf
./mdpdf
```

## Firejail

The `pandoc` command is called with `--sandbox=true`, but to increase the sandbox even further, you can use `firejail`.

You will need to create a firejail profile, in this example I am going to call it `mdpdf.profile`:

```plain
whitelist /absolute/path/to/mdpdf/template.latex
```

This will work if you have put the `mdpdf` binary on your PATH (which is accessible by the root user.)

If you don't want to add `mdpdf` to your PATH, then you need to add this additional line to `mdpdf.profile`:

```plain
whitelist /absolute/path/to/mdpdf/mdpdf
```

If you want to add `mdpdf` to your PATH, you can do that like this:

```shell
cp mdpdf /usr/bin/
```

(and here `mdpdf` is the binary file, not the folder)

And now you can start the service:

```shell
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
