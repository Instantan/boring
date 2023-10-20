package cli

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ChangeTracker struct {
	l      sync.Mutex
	hashes map[string][]byte
	hasher hash.Hash
}

func NewChangeTracker() *ChangeTracker {
	return &ChangeTracker{
		l:      sync.Mutex{},
		hashes: map[string][]byte{},
		hasher: sha256.New(),
	}
}

func (ct *ChangeTracker) Init() {
	ct.l.Lock()
	relevantFiles, _ := getRelevantFiles(".", relevantForRecompilationFileExtensions)
	for i := range relevantFiles {
		path := relevantFiles[i]
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		ct.hasher.Reset()
		if _, err := io.Copy(ct.hasher, f); err != nil {
			continue
		}
		hash := ct.hasher.Sum(nil)
		ct.hashes[path] = hash
		f.Close()
	}
	ct.l.Unlock()
}

func (ct *ChangeTracker) InitFile(path string) {
	ct.l.Lock()
	f, err := os.Open(path)
	if err != nil {
		ct.l.Unlock()
		return
	}
	ct.hasher.Reset()
	if _, err := io.Copy(ct.hasher, f); err != nil {
		ct.l.Unlock()
		return
	}
	hash := ct.hasher.Sum(nil)
	ct.hashes[path] = hash
	f.Close()
	ct.l.Unlock()
}

func (ct *ChangeTracker) DidFileChange(path string) bool {
	ct.l.Lock()
	f, err := os.Open(path)
	if err != nil {
		ct.l.Unlock()
		return true
	}

	if ct.isGeneratedTemplFile(path, f) {
		ct.l.Unlock()
		f.Close()
		return false
	}

	ct.hasher.Reset()
	if _, err := io.Copy(ct.hasher, f); err != nil {
		ct.l.Unlock()
		f.Close()
		return true
	}
	hash := ct.hasher.Sum(nil)
	eq := bytes.Equal(ct.hashes[path], hash)
	ct.hashes[path] = hash
	f.Close()
	ct.l.Unlock()
	return !eq
}

func (ct *ChangeTracker) isGeneratedTemplFile(path string, f *os.File) bool {
	if filepath.Ext(path) != ".go" {
		return false
	}
	return ct.hasGeneratedFileHeader(f)
}

func (ct *ChangeTracker) hasGeneratedFileHeader(f *os.File) bool {
	prefix := false
	s := bufio.NewScanner(f)
	if s.Scan() {
		prefix = strings.HasPrefix(s.Text(), "// Code generated by")
	}
	f.Seek(0, 0)
	return prefix
}
