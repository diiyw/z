package file

import (
	"sync"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	fileCache *Cache
)

type Cache struct {
	documents map[string]protocol.TextDocumentItem // uri -> document
	mu        sync.RWMutex
}

func Document() *Cache {
	if fileCache != nil {
		return fileCache
	}
	fileCache = &Cache{
		documents: make(map[string]protocol.TextDocumentItem),
	}
	return fileCache
}

func (f *Cache) Set(uri string, doc protocol.TextDocumentItem) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.documents[uri] = doc
}

func (f *Cache) Delete(uri string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.documents, uri)
}

func (f *Cache) Get(uri string) (protocol.TextDocumentItem, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	doc, ok := f.documents[uri]
	return doc, ok
}

func (f *Cache) GetText(uri string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	doc := f.documents[uri]
	return doc.Text
}
