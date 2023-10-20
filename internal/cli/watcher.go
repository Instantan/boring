package cli

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fileWatcher   *fsnotify.Watcher
	compilerPool  *CompilerPool
	changeTracker *ChangeTracker
	scheduled     *sync.Mutex
}

func NewWatcher() *Watcher {
	watcher := &Watcher{}
	fileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	watcher.fileWatcher = fileWatcher
	watcher.changeTracker = NewChangeTracker()
	watcher.compilerPool = NewCompilerPool(watcher.changeTracker)
	watcher.scheduled = &sync.Mutex{}
	return watcher
}

func (w *Watcher) Run() error {
	defer w.fileWatcher.Close()

	w.scheduleBuildAndRun()
	w.changeTracker.Init()

	err := w.fileWatcher.Add("./")
	if err != nil {
		return err
	}
	w.watchSubfolders()

	relevantFiles, _ := getRelevantFiles("./", relevantForRecompilationFileExtensions)
	for i := range relevantFiles {
		w.changeTracker.DidFileChange(relevantFiles[i])
	}

	for {
		select {
		case event, ok := <-w.fileWatcher.Events:
			if !ok {
				return nil
			}
			if event.Op.Has(fsnotify.Create) {
				if err := w.watchFolder(event.Name); err != nil {
					printError(err)
				}
			} else if event.Op.Has(fsnotify.Remove) {
				if err := w.unwatchFolder(event.Name); err != nil {
					printError(err)
				}
			}

			w.scheduled.Lock()
			if !isFileRelevant(event.Name, relevantForRecompilationFileExtensions) || !w.changeTracker.DidFileChange(event.Name) {
				w.scheduled.Unlock()
				continue
			}
			w.scheduleBuildAndRun()
			w.scheduled.Unlock()
		case err, ok := <-w.fileWatcher.Errors:
			if !ok {
				return nil
			}
			printError(err)
		}
	}
}

func (w *Watcher) watchSubfolders() error {
	return filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if shouldSkipDir(info.Name()) {
					return filepath.SkipDir
				}
				w.fileWatcher.Add(path)
			}
			return nil
		})
}

func (w *Watcher) watchFolder(path string) error {
	path, err := absPathCwd(path)
	if err != nil {
		return err
	}
	if isDir, err := isDirectory(path); !isDir {
		return err
	}
	printInternal("watching: %v", path)
	return w.fileWatcher.Add(path)
}

func (w *Watcher) unwatchFolder(path string) error {
	path, err := absPathCwd(path)
	if err != nil {
		return err
	}
	if isDir, err := isDirectory(path); !isDir {
		return err
	}
	printInternal("watching: %v", path)
	return w.fileWatcher.Remove(path)
}

func (w *Watcher) scheduleBuildAndRun() {
	println()
	w.compilerPool.GenerateAssetsAndTempl()
	w.compilerPool.RunGo()
	time.Sleep(time.Millisecond * 1000)
}
