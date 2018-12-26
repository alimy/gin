// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"
	"sync"
)

// RenderType Names
const (
	JSONRenderType           = iota // JSON Render Type
	IntendedJSONRenderType          // IntendedJSON Render Type
	PureJSONRenderType              // PureJSON Render Type
	AsciiJSONRenderType             // AsciiJSON Render Type
	JsonpJSONRenderType             // JsonpJSON Render Type
	SecureJSONRenderType            // SecureJSON Type
	XMLRenderType                   // XML Render Type
	StringRenderType                // String Render Type
	RedirectRenderType              // Redirect Render Type
	DataRenderType                  // Data Render Type
	HTMLRenderType                  // HTML Render Type
	HTMLDebugRenderType             // HTMLDebug Render Type
	HTMLProductionRenderType        // HTMLProduction Render Type
	YAMLRenderType                  // YAML Render Type
	MsgPackRenderType               // MsgPack Render Type
	ReaderRenderType                // Reader Render Type
	ProtoBufRenderType              // ProtoBuf Render Type
	EmptyRenderType                 // Empty Render Type
)

var (
	renderPool = &RenderPool{renderPools: make(map[int]sync.Pool)}
)

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Render interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(w http.ResponseWriter)
}

type RenderRecycler interface {
	Render
	// Setup set data and opts
	Setup(data interface{}, opts ...interface{})
	// Reset clean data and opts
	Reset()
}

// HTMLRender interface is to be implemented by HTMLProduction and HTMLDebug.
type RenderFactory interface {
	// Instance apply opts to build a new Render instance
	Instance() Render
}

// RenderPool contains Render instance
type RenderPool struct {
	mu          sync.RWMutex
	renderPools map[int]sync.Pool
}

func (p *RenderPool) get(name int) RenderRecycler {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var render RenderRecycler
	if pool, ok := p.renderPools[name]; ok {
		render, _ = pool.Get().(RenderRecycler)
	} else {
		pool, _ = p.renderPools[EmptyRenderType]
		render, _ = pool.Get().(RenderRecycler)
	}
	if render == nil {
		render = &EmptyRender{}
	}
	return render
}

func (p *RenderPool) put(name int, render RenderRecycler) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if pool, ok := p.renderPools[name]; ok {
		render.Reset()
		pool.Put(render)
	}
}

func (p *RenderPool) register(name int, factory RenderFactory) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if factory == nil {
		panic("gin: Register RenderFactory is nil")
	}
	if _, dup := p.renderPools[name]; dup {
		panic("gin: Register called twice for RenderFactory")
	}

	p.renderPools[name] = sync.Pool{
		New: func() interface{} {
			return factory.Instance()
		},
	}
}

// Register makes a binding available by the provided name.
// If Register is called twice with the same name or if binding is nil,
// it panics.
func Register(name int, factory RenderFactory) {
	renderPool.register(name, factory)
}

// Default returns the appropriate Render instance based on the render type.
func Default(name int) RenderRecycler {
	return renderPool.get(name)
}

// Recycle put render to sync.Pool
func Recycle(name int, render RenderRecycler) {
	renderPool.put(name, render)
}

func WriteContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
