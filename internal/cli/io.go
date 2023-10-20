package cli

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readFile(file string, autoresolver ...string) (string, error) {
	if len(file) == 0 {
		return "", errors.New("filepath cant be empty")
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if file[0] != '/' {
		file = pwd + "/" + file
	} else {
		file = pwd + file
	}

	if len(autoresolver) > 0 && !fileExists(file) {
		for _, resolve := range autoresolver {
			if fileExists(file + resolve) {
				file = file + resolve
				break
			}
		}
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func getRelevantFiles(directory string, fileExtensions []string) ([]string, error) {
	absolutePaths := []string{}
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}

			if !info.IsDir() && isFileRelevant(path, fileExtensions) {
				absolutePaths = append(absolutePaths, path)
			}

			return nil
		})
	return absolutePaths, err
}

func shouldSkipDir(name string) bool {
	return name == "node_modules" || (len(name) > 0 && name[0] == '_')
}

func isFileRelevant(path string, fileExtension []string) bool {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(filepath.Base(path), ext)
	if len(base) == 0 || base[0] == '_' || strings.HasSuffix(base, ".min") {
		return false
	}
	for i := range fileExtension {
		if ext == fileExtension[i] {
			return true
		}
	}
	return false
}

func turnFilepathIntoMinifiedVersion(path string) string {
	suff := filepath.Ext(path)
	raw := strings.TrimSuffix(path, suff)
	return raw + ".min" + suff
}

func changeExtensionOfFilepath(path string, ext string) string {
	raw := strings.TrimSuffix(path, filepath.Ext(path))
	return raw + ext
}

func overwriteFile(path string, contents []byte) error {
	return os.WriteFile(path, contents, 0666)
}

func isDirectoryEmpty(pathToDir string, ignoredFilepaths []string) (bool, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return false, err
	}
	testStr := strings.ToLower(strings.Join(ignoredFilepaths, "; "))
	for i := range entries {
		if !strings.Contains(testStr, strings.ToLower(entries[i].Name())) {
			return false, nil
		}
	}
	return true, nil
}

func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func absPathCwd(path string) (string, error) {
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(wd, path)
	}
	return path, nil
}
