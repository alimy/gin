package render

import "net/http"

type EmptyRenderFactory struct {}

type EmptyRender struct {}

// Instance apply opts to build a new EmptyRender instance
func (_ EmptyRenderFactory) Instance(data interface{}, opts ...interface{}) Render {
	return &EmptyRender{}
}

// Render writes data with custom ContentType.
func (_ EmptyRender)Render(http.ResponseWriter) error {
	return nil
}

// WriteContentType writes custom ContentType.
func (_ EmptyRender)WriteContentType(w http.ResponseWriter) {
	return
}