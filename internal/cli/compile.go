package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type CompilerPool struct {
	compileLock   sync.Mutex
	scssPool      ScssCompilerPool
	goProcess     *exec.Cmd
	changeTracker *ChangeTracker
}

func NewCompilerPool(changeTracker *ChangeTracker) *CompilerPool {
	return &CompilerPool{
		compileLock:   sync.Mutex{},
		scssPool:      NewScssCompilerPool(),
		changeTracker: changeTracker,
	}
}

func (p *CompilerPool) CompileDirectory(directory string) error {
	p.compileLock.Lock()
	defer p.compileLock.Unlock()

	relevantFiles, err := getRelevantFiles(directory, relevantForJsAndScssBuildingFileExtensions)
	fmt.Println(relevantFiles)
	if err != nil {
		return err
	}

	readWg := &sync.WaitGroup{}
	writeWg := &sync.WaitGroup{}
	errLock := &sync.Mutex{}
	errWg := &sync.WaitGroup{}

	errs := []error{}

	for i := range relevantFiles {
		readWg.Add(1)
		errWg.Add(1)
		go func(i int, l *sync.Mutex, errWg *sync.WaitGroup) {
			err := p.compileFile(relevantFiles[i], readWg, writeWg)
			if err != nil {
				errLock.Lock()
				errs = append(errs, err)
				errLock.Unlock()
			}
			if p.changeTracker != nil {
				p.changeTracker.InitFile(relevantFiles[i])
			}
			errWg.Done()
		}(i, errLock, errWg)
	}

	readWg.Wait()
	writeWg.Wait()

	errWg.Wait()
	errLock.Lock()
	err = errors.Join(errs...)
	errLock.Unlock()

	return err
}

func (p *CompilerPool) stopRunningGoProcess() {
	if p.goProcess != nil && p.goProcess.Process != nil {
		printInternal("\nstopping go process")
		p.goProcess.Process.Kill()
	}
}

func (p *CompilerPool) TestGo() {
	p.stopRunningGoProcess()
	fmt.Println("go test")
	p.goProcess = exec.Command("go", "test")
	p.goProcess.Stdout = os.Stdout
	p.goProcess.Stderr = os.Stderr
	p.goProcess.Stdin = os.Stdin
	if err := p.goProcess.Start(); err != nil {
		printError(err)
	}
}

func (p *CompilerPool) RunGo() {
	p.stopRunningGoProcess()
	fmt.Println("go run .")
	p.goProcess = exec.Command("go", "run", ".")
	p.goProcess.Stdout = os.Stdout
	p.goProcess.Stderr = os.Stderr
	p.goProcess.Stdin = os.Stdin
	if err := p.goProcess.Start(); err != nil {
		printError(err)
	}
}

func (p *CompilerPool) BuildGo() {
	p.stopRunningGoProcess()
	fmt.Println(`go build`)
	p.goProcess = exec.Command("go", "build")
	p.goProcess.Stdout = os.Stdout
	p.goProcess.Stderr = os.Stderr
	p.goProcess.Stdin = os.Stdin
	if err := p.goProcess.Start(); err != nil {
		printError(err)
	}
}

func (p *CompilerPool) GenerateAssetsAndTempl() {
	parallel(func() {
		printMeasuredAction("building assets files", func() error {
			return p.CompileDirectory(".")
		}, "finished building asset files")
	}, func() {
		printMeasuredAction("generating templ files", func() error {
			return p.GenerateTempl()
		}, "finished generating templ files")
	})
}

func (p *CompilerPool) GenerateTempl() error {
	return generateTempl()
}

func (p *CompilerPool) compileFile(path string, readWg, writeWg *sync.WaitGroup) error {
	ext := filepath.Ext(path)
	switch ext {
	case ".js":
		return p.compileJsFile(path, readWg, writeWg)
	case ".ts":
		return p.compileJsFile(path, readWg, writeWg)
	case ".css":
		return p.compileScssFile(path, readWg, writeWg)
	case ".scss":
		return p.compileScssFile(path, readWg, writeWg)
	default:
		return nil
	}
}

func (p *CompilerPool) compileJsFile(path string, readWg, writeWg *sync.WaitGroup) error {
	res, err := buildJsFile(path)
	readWg.Done()
	if err != nil {
		return err
	}
	errs := []error{}
	for i := range res {
		errs = append(errs, overwriteFile(turnFilepathIntoMinifiedVersion(res[i].Path), res[i].Contents))
	}
	return errors.Join(errs...)
}

func (p *CompilerPool) compileScssFile(path string, readWg, writeWg *sync.WaitGroup) error {
	contents, err := p.scssPool.CompileFile(path)
	readWg.Done()
	if err != nil {
		return err
	}
	return overwriteFile(turnFilepathIntoMinifiedVersion(changeExtensionOfFilepath(path, ".css")), contents)
}
