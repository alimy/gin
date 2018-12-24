package xml

import (
	"github.com/alimy/gin/binding"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestBindingBody(t *testing.T) {
	for _, tt := range []struct {
		name    string
		binding binding.BindingBody
		body    string
		want    string
	}{
		{
			name:    "XML binding",
			binding: xmlBinding{},
			body: `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <foo>FOO</foo>
</root>`,
		},
	} {
		t.Logf("testing: %s", tt.name)
		req := requestWithBody("POST", "/", tt.body)
		form := FooStruct{}
		body, _ := ioutil.ReadAll(req.Body)
		assert.NoError(t, tt.binding.BindBody(body, &form))
		assert.Equal(t, FooStruct{"FOO"}, form)
	}
}
