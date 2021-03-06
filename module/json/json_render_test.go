// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo":  "bar",
		"html": "<b>",
	}

	(&JSONRender{data}).WriteContentType(w)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	err := (&JSONRender{data}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"\\u003cb\\u003e\"}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderJSONPanics(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	assert.Panics(t, func() { (&JSONRender{data}).Render(w) })
}

func TestRenderIndentedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
		"bar": "foo",
	}

	err := (&IndentedJSONRender{data}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "{\n    \"bar\": \"foo\",\n    \"foo\": \"bar\"\n}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderIndentedJSONPanics(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	err := (&IndentedJSONRender{data}).Render(w)
	assert.Error(t, err)
}

func TestRenderSecureJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}

	(&SecureJSONRender{"while(1);", data}).WriteContentType(w1)
	assert.Equal(t, "application/json; charset=utf-8", w1.Header().Get("Content-Type"))

	err1 := (&SecureJSONRender{"while(1);", data}).Render(w1)

	assert.NoError(t, err1)
	assert.Equal(t, "{\"foo\":\"bar\"}", w1.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	datas := []map[string]interface{}{{
		"foo": "bar",
	}, {
		"bar": "foo",
	}}

	err2 := (&SecureJSONRender{"while(1);", datas}).Render(w2)
	assert.NoError(t, err2)
	assert.Equal(t, "while(1);[{\"foo\":\"bar\"},{\"bar\":\"foo\"}]", w2.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w2.Header().Get("Content-Type"))
}

func TestRenderSecureJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	err := (&SecureJSONRender{"while(1);", data}).Render(w)
	assert.Error(t, err)
}

func TestRenderJsonpJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}

	(&JsonpJSONRender{"x", data}).WriteContentType(w1)
	assert.Equal(t, "application/javascript; charset=utf-8", w1.Header().Get("Content-Type"))

	err1 := (&JsonpJSONRender{"x", data}).Render(w1)

	assert.NoError(t, err1)
	assert.Equal(t, "x({\"foo\":\"bar\"})", w1.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	datas := []map[string]interface{}{{
		"foo": "bar",
	}, {
		"bar": "foo",
	}}

	err2 := (&JsonpJSONRender{"x", datas}).Render(w2)
	assert.NoError(t, err2)
	assert.Equal(t, "x([{\"foo\":\"bar\"},{\"bar\":\"foo\"}])", w2.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w2.Header().Get("Content-Type"))
}

func TestRenderJsonpJSONError2(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}
	(&JsonpJSONRender{"", data}).WriteContentType(w)
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))

	e := (&JsonpJSONRender{"", data}).Render(w)
	assert.NoError(t, e)

	assert.Equal(t, "{\"foo\":\"bar\"}", w.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderJsonpJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	err := (&JsonpJSONRender{"x", data}).Render(w)
	assert.Error(t, err)
}

func TestRenderAsciiJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data1 := map[string]interface{}{
		"lang": "GO语言",
		"tag":  "<br>",
	}

	err := (&AsciiJSONRender{data1}).Render(w1)

	assert.NoError(t, err)
	assert.Equal(t, "{\"lang\":\"GO\\u8bed\\u8a00\",\"tag\":\"\\u003cbr\\u003e\"}", w1.Body.String())
	assert.Equal(t, "application/json", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	data2 := float64(3.1415926)

	err = (&AsciiJSONRender{data2}).Render(w2)
	assert.NoError(t, err)
	assert.Equal(t, "3.1415926", w2.Body.String())
}

func TestRenderAsciiJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	assert.Error(t, (&AsciiJSONRender{data}).Render(w))
}
