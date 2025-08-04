package main

import (
	"sync"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	cache *DocumentCache
)

type DocumentCache struct {
	documents map[string]protocol.TextDocumentItem // uri -> document
	mu        sync.RWMutex
}

func Document() *DocumentCache {
	if cache != nil {
		return cache
	}
	cache = &DocumentCache{
		documents: make(map[string]protocol.TextDocumentItem),
	}
	return cache
}

func (c *DocumentCache) Set(uri string, doc protocol.TextDocumentItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.documents[uri] = doc
}

func (c *DocumentCache) Delete(uri string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.documents, uri)
}

func (c *DocumentCache) Get(uri string) (protocol.TextDocumentItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc, ok := c.documents[uri]
	return doc, ok
}

func (c *DocumentCache) GetText(uri string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc := c.documents[uri]
	return doc.Text
}
