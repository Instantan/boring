# boring

> CLI for boring webdev with go

## Install
```bash
go install github.com/Instantan/boring
```

## Commands
```bash
boring 
```
Hot code reloading of .scss, .css, .ts, .js, .templ and .go files

```bash
boring generate
```
Generates, bundles, minifies .scss, .css, .ts, .js and .templ files

```bash
boring help
```
Prints the help screen


## More informations

### SCSS
Boring bundles the dart-sass compiler and automatically compiles and minifies all scss and css files

It treats every scss or css file as a entrypoint, wich means if you have the following three scss/css files:

pico.css
lib.scss
main.scss

It will produce:

pico.min.css
lib.min.css
main.min.css

To ignore a scss or css file from the compiler prefix it with an underscore.
So if you only want the main file because you import the other files in it your directory should look like:

_pico.css
_lib.scss
main.scss

Wich will produce the following file:

main.min.css

### ESBuild
It also comes with esbuild to bundle and minify all js/ts files.
It treats every js file as a entrypoint wich means you get a compiled and minified version for every js file.
To exclude a js file from the bundler prefix is with a underscore. For more information see the SCSS description, it works exactly the same way

### Templ
Templ is used to generate efficient HTML Templates.

## Build with
- [https://github.com/sass/dart-sass](https://www.npmjs.com/package/sass)
- [https://github.com/evanw/esbuild](https://github.com/evanw/esbuild)
- [https://github.com/a-h/templ](https://github.com/a-h/templ)
- [https://github.com/dop251/goja](https://github.com/dop251/goja)