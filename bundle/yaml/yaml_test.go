package yaml

import (
	"bytes"
	"github.com/alimy/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestContextRenderYAML tests that the response is serialized as YAML
// and Content-Type is set to application/x-yaml
func TestContextRenderYAML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.YAML(http.StatusCreated, gin.H{"foo": "bar"})

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "foo: bar\n", w.Body.String())
	assert.Equal(t, "application/x-yaml; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestContextShouldBindWithYAML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("foo: bar\nbar: foo"))
	c.Request.Header.Add("Content-Type", gin.MIMEXML) // set fake content-type

	var obj struct {
		Foo string `yaml:"foo"`
		Bar string `yaml:"bar"`
	}
	assert.NoError(t, c.ShouldBindYAML(&obj))
	assert.Equal(t, "foo", obj.Bar)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, 0, w.Body.Len())
}

func TestContextBindWithYAML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("foo: bar\nbar: foo"))
	c.Request.Header.Add("Content-Type", gin.MIMEXML) // set fake content-type

	var obj struct {
		Foo string `yaml:"foo"`
		Bar string `yaml:"bar"`
	}
	assert.NoError(t, c.BindYAML(&obj))
	assert.Equal(t, "foo", obj.Bar)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, 0, w.Body.Len())
}
