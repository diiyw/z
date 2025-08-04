package main

import (
	"sync"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	cache *DocmentCache
)

type DocmentCache struct {
	documents map[string]protocol.TextDocumentItem // uri -> document
	mu        sync.RWMutex
}

func Document() *DocmentCache {
	if cache != nil {
		return cache
	}
	cache = &DocmentCache{
		documents: make(map[string]protocol.TextDocumentItem),
	}
	return cache
}

func (c *DocmentCache) Set(uri string, doc protocol.TextDocumentItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.documents[uri] = doc
}

func (c *DocmentCache) Get(uri string) (protocol.TextDocumentItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc, ok := c.documents[uri]
	return doc, ok
}

func (c *DocmentCache) GetText(uri string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc := c.documents[uri]
	return doc.Text
}
