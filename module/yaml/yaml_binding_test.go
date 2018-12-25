// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package yaml

import (
	"bytes"
	"github.com/alimy/gin/binding"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required"`
}

func TestBindingDefault(t *testing.T) {
	assert.Equal(t, &yamlBinding{}, binding.DefaultWith("POST", binding.MIMEYAML))
	assert.Equal(t, &yamlBinding{}, binding.DefaultWith("PUT", binding.MIMEYAML))
}

func TestBindingYAML(t *testing.T) {
	testBodyBinding(t,
		yamlBinding{}, "yaml",
		"/", "/",
		`foo: bar`, `bar: foo`)
}

func TestBindingYAMLFail(t *testing.T) {
	testBodyBindingFail(t,
		yamlBinding{}, "yaml",
		"/", "/",
		`foo:\nbar`, `bar: foo`)
}

func testBodyBinding(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
}

func testBodyBindingFail(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
	assert.Equal(t, "", obj.Foo)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
