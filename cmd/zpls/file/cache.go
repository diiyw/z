package main

import (
	"sync"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	fileCache *FileCache
)

type FileCache struct {
	documents map[string]protocol.TextDocumentItem // uri -> document
	mu        sync.RWMutex
}

func File() *FileCache {
	if fileCache != nil {
		return fileCache
	}
	fileCache = &FileCache{
		documents: make(map[string]protocol.TextDocumentItem),
	}
	return fileCache
}

func (f *FileCache) Set(uri string, doc protocol.TextDocumentItem) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.documents[uri] = doc
}

func (f *FileCache) Delete(uri string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.documents, uri)
}

func (f *FileCache) Get(uri string) (protocol.TextDocumentItem, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	doc, ok := f.documents[uri]
	return doc, ok
}

func (f *FileCache) GetText(uri string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	doc := f.documents[uri]
	return doc.Text
}
