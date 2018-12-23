// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"
	"sync"
)

// RenderFactory Names
const (
	JSONRenderFactory           = iota // JSON Render Factory
	IntendedJSONRenderFactory          // IntendedJSON Render Factory
	PureJSONRenderFactory              // PureJSON Render Factory
	AsciiJSONRenderFactory             // AsciiJSON Render Factory
	JsonpJSONRenderFactory             // JsonpJSON Render Factory
	SecureJSONRenderFactory            // SecureJSON RenderFactory
	XMLRenderFactory                   // XML Render Factory
	StringRenderFactory                // String Render Factory
	RedirectRenderFactory              // Redirect Render Factory
	DataRenderFactory                  // Data Render Factory
	HTMLRenderFactory                  // HTML Render Factory
	HTMLDebugRenderFactory             // HTMLDebug Render Factory
	HTMLProductionRenderFactory        // HTMLProduction Render Factory
	YAMLRenderFactory                  // YAML Render Factory
	MsgPackRenderFactory               // MsgPack Render Factory
	ReaderRenderFactory                // Reader Render Factory
	ProtoBufRenderFactory              // ProtoBuf Render Factory
)

var (
	renderFactoriesMu sync.RWMutex
	renderFactories   = make(map[int]RenderFactory)
)

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Render interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(w http.ResponseWriter)
}

// HTMLRender interface is to be implemented by HTMLProduction and HTMLDebug.
type RenderFactory interface {
	// Instance apply opts to build a new Render instance
	Instance(data interface{}, opts ...interface{}) Render
}

// Register makes a binding available by the provided name.
// If Register is called twice with the same name or if binding is nil,
// it panics.
func Register(name int, factory RenderFactory) {
	renderFactoriesMu.Lock()
	defer renderFactoriesMu.Unlock()

	if factory == nil {
		panic("gin: Register RenderFactory is nil")
	}
	if _, dup := renderFactories[name]; dup {
		panic("gin: Register called twice for RenderFactories")
	}

	renderFactories[name] = factory
}

// Default returns the appropriate RenderFactory instance based on the render type.
func Default(name int) RenderFactory {
	renderFactoriesMu.RLock()
	defer renderFactoriesMu.RUnlock()

	if renderFactory, ok := renderFactories[name]; ok {
		return renderFactory
	} else {
		return EmptyRenderFactory{}
	}
}

func WriteContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
