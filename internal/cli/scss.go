package cli

import (
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/Instantan/boring/internal/bundles"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/dop251/goja_nodejs/url"
)

type ScssFile struct {
	Contents string `json:"contents"`
	Syntax   string `json:"syntax"`
}

type ScssCompilerPool struct {
	pool *sync.Pool
}

type ScssCompiler func(string) (string, error)

func NewScssCompilerPool() ScssCompilerPool {
	pool := ScssCompilerPool{}
	pool.pool = &sync.Pool{
		New: func() any {
			return pool.newCompiler()
		},
	}
	return pool
}

func (p ScssCompilerPool) newCompiler() ScssCompiler {
	scssRuntime := p.newJsRuntime()

	prog, err := goja.Compile("scss", string(bundles.SCSSSource()), false)
	if err != nil {
		panic(err)
	}
	_, err = scssRuntime.RunProgram(prog)
	if err != nil {
		panic(err)
	}

	scssRuntime.Set("canonicalize", p.canonicalizeScssUrl)
	scssRuntime.Set("loadFile", p.loadScssImport)

	fn, ok := goja.AssertFunction(scssRuntime.Get("compileScssToCss"))
	if !ok {
		panic("Not a function")
	}
	return ScssCompiler(func(scss string) (string, error) {
		res, err := fn(goja.Undefined(), scssRuntime.ToValue(scss))
		if err != nil {
			return "", err
		}
		return res.String(), err
	})
}

func (p ScssCompilerPool) newJsRuntime() *goja.Runtime {
	scssRuntime := goja.New()
	scssRuntime.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	new(require.Registry).Enable(scssRuntime)
	url.Enable(scssRuntime)
	console.Enable(scssRuntime)
	return scssRuntime
}

func (p ScssCompilerPool) GetCompiler() ScssCompiler {
	return p.pool.Get().(ScssCompiler)
}

func (p ScssCompilerPool) PutCompiler(compiler ScssCompiler) {
	p.pool.Put(compiler)
}

func (p ScssCompilerPool) CompileString(scss string) (string, error) {
	c := p.GetCompiler()
	res, err := c(scss)
	p.PutCompiler(c)
	return res, err
}

func (p ScssCompilerPool) CompileFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	result, err := p.CompileString(string(data))
	return []byte(result), err
}

func (p ScssCompilerPool) canonicalizeScssUrl(url string) string {
	switch true {
	case p.isFileImport(url):
		return url
	case p.isWebImport(url):
		return url
	default: // defaults to file imports
		return "file:" + url
	}
}

func (p ScssCompilerPool) isFileImport(url string) bool {
	return strings.HasPrefix(url, "file:")
}

func (p ScssCompilerPool) isWebImport(url string) bool {
	return strings.HasPrefix(url, "http:") || strings.HasPrefix(url, "https:")
}

func (p ScssCompilerPool) loadScssImport(url string) *ScssFile {
	var contents string
	var err error
	switch true {
	case p.isFileImport(url):
		contents, err = readFile(url[7:])
	case p.isWebImport(url):
		contents, err = getRequest(url)
	default:
		contents, err = "", errors.New("unkown type of import")
	}
	if err != nil {
		return nil
	}
	return &ScssFile{
		Contents: contents,
		Syntax:   "scss",
	}
}
