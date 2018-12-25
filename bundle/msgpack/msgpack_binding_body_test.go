package msgpack

import (
	"bytes"
	"github.com/alimy/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
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
			name:    "MsgPack binding",
			binding: msgpackBinding{},
			body:    msgPackBody(t),
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

func msgPackBody(t *testing.T) string {
	test := FooStruct{"FOO"}
	h := new(codec.MsgpackHandle)
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, codec.NewEncoder(buf, h).Encode(test))
	return buf.String()
}
