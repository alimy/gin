// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required"`
}

type FooBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar" xml:"bar" binding:"required"`
}

type FooDefaultBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar,default=hello" xml:"bar" binding:"required"`
}

type FooBarStructForTimeType struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2006-01-02" time_utc:"1" time_location:"Asia/Chongqing"`
	TimeBar time.Time `form:"time_bar" time_format:"2006-01-02" time_utc:"1"`
}

type FooStructForTimeTypeNotFormat struct {
	TimeFoo time.Time `form:"time_foo"`
}

type FooStructForTimeTypeFailFormat struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2017-11-15"`
}

type FooStructForTimeTypeFailLocation struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2006-01-02" time_location:"/asia/chongqing"`
}

type FooStructForMapType struct {
	// Unknown type: not support map
	MapFoo map[string]interface{} `form:"map_foo"`
}

type InvalidNameType struct {
	TestName string `invalid_name:"test_name"`
}

type InvalidNameMapType struct {
	TestName struct {
		MapFoo map[string]interface{} `form:"map_foo"`
	}
}

type FooStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
}

func TestBindingForm(t *testing.T) {
	testFormBinding(t, "POST",
		"/", "/",
		"foo=bar&bar=foo", "bar2=foo")
}

func TestBindingForm2(t *testing.T) {
	testFormBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"", "")
}

func TestBindingFormDefaultValue(t *testing.T) {
	testFormBindingDefaultValue(t, "POST",
		"/", "/",
		"foo=bar", "bar2=foo")
}

func TestBindingFormDefaultValue2(t *testing.T) {
	testFormBindingDefaultValue(t, "GET",
		"/?foo=bar", "/?bar2=foo",
		"", "")
}

func TestBindingFormForTime(t *testing.T) {
	testFormBindingForTime(t, "POST",
		"/", "/",
		"time_foo=2017-11-15&time_bar=", "bar2=foo")
	testFormBindingForTimeNotFormat(t, "POST",
		"/", "/",
		"time_foo=2017-11-15", "bar2=foo")
	testFormBindingForTimeFailFormat(t, "POST",
		"/", "/",
		"time_foo=2017-11-15", "bar2=foo")
	testFormBindingForTimeFailLocation(t, "POST",
		"/", "/",
		"time_foo=2017-11-15", "bar2=foo")
}

func TestBindingFormForTime2(t *testing.T) {
	testFormBindingForTime(t, "GET",
		"/?time_foo=2017-11-15&time_bar=", "/?bar2=foo",
		"", "")
	testFormBindingForTimeNotFormat(t, "GET",
		"/?time_foo=2017-11-15", "/?bar2=foo",
		"", "")
	testFormBindingForTimeFailFormat(t, "GET",
		"/?time_foo=2017-11-15", "/?bar2=foo",
		"", "")
	testFormBindingForTimeFailLocation(t, "GET",
		"/?time_foo=2017-11-15", "/?bar2=foo",
		"", "")
}

func TestBindingFormInvalidName(t *testing.T) {
	testFormBindingInvalidName(t, "POST",
		"/", "/",
		"test_name=bar", "bar2=foo")
}

func TestBindingFormInvalidName2(t *testing.T) {
	testFormBindingInvalidName2(t, "POST",
		"/", "/",
		"map_foo=bar", "bar2=foo")
}

func TestBindingQuery(t *testing.T) {
	testQueryBinding(t, "POST",
		"/?foo=bar&bar=foo", "/",
		"foo=unused", "bar2=foo")
}

func TestBindingQuery2(t *testing.T) {
	testQueryBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"foo=unused", "")
}

func TestBindingQueryFail(t *testing.T) {
	testQueryBindingFail(t, "POST",
		"/?map_foo=", "/",
		"map_foo=unused", "bar2=foo")
}

func TestBindingQueryFail2(t *testing.T) {
	testQueryBindingFail(t, "GET",
		"/?map_foo=", "/?bar2=foo",
		"map_foo=unused", "")
}

func TestBindingQueryBoolFail(t *testing.T) {
	testQueryBindingBoolFail(t, "GET",
		"/?bool_foo=fasl", "/?bar2=foo",
		"bool_foo=unused", "")
}

func createFormPostRequest() *http.Request {
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createDefaultFormPostRequest() *http.Request {
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar"))
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestFail() *http.Request {
	req, _ := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo=bar"))
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormMultipartRequest() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	mw.SetBoundary(boundary)
	mw.WriteField("foo", "bar")
	mw.WriteField("bar", "foo")
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func createFormMultipartRequestFail() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	mw.SetBoundary(boundary)
	mw.WriteField("map_foo", "bar")
	req, _ := http.NewRequest("POST", "/?map_foo=getfoo", body)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func TestBindingFormPost(t *testing.T) {
	req := createFormPostRequest()
	var obj FooBarStruct
	FormPost.Bind(req, &obj)

	assert.Equal(t, "form-urlencoded", FormPost.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingDefaultValueFormPost(t *testing.T) {
	req := createDefaultFormPostRequest()
	var obj FooDefaultBarStruct
	FormPost.Bind(req, &obj)

	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)
}

func TestBindingFormPostFail(t *testing.T) {
	req := createFormPostRequestFail()
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingFormMultipart(t *testing.T) {
	req := createFormMultipartRequest()
	var obj FooBarStruct
	FormMultipart.Bind(req, &obj)

	assert.Equal(t, "multipart/form-data", FormMultipart.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingFormMultipartFail(t *testing.T) {
	req := createFormMultipartRequestFail()
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	assert.Error(t, err)
}

func TestUriBinding(t *testing.T) {
	b := Uri
	assert.Equal(t, "uri", b.Name())

	type Tag struct {
		Name string `uri:"name"`
	}
	var tag Tag
	m := make(map[string][]string)
	m["name"] = []string{"thinkerou"}
	assert.NoError(t, b.BindUri(m, &tag))
	assert.Equal(t, "thinkerou", tag.Name)

	type NotSupportStruct struct {
		Name map[string]interface{} `uri:"name"`
	}
	var not NotSupportStruct
	assert.Error(t, b.BindUri(m, &not))
	assert.Equal(t, map[string]interface{}(nil), not.Name)
}

func testFormBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func testFormBindingDefaultValue(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooDefaultBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)
}

func TestFormBindingFail(t *testing.T) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormPostBindingFail(t *testing.T) {
	b := FormPost
	assert.Equal(t, "form-urlencoded", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormMultipartBindingFail(t *testing.T) {
	b := FormMultipart
	assert.Equal(t, "multipart/form-data", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTime(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStructForTimeType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)

	assert.NoError(t, err)
	assert.Equal(t, int64(1510675200), obj.TimeFoo.Unix())
	assert.Equal(t, "Asia/Chongqing", obj.TimeFoo.Location().String())
	assert.Equal(t, int64(-62135596800), obj.TimeBar.Unix())
	assert.Equal(t, "UTC", obj.TimeBar.Location().String())
}

func testFormBindingForTimeNotFormat(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeNotFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeFailFormat(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeFailFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeFailLocation(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeFailLocation{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingInvalidName(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := InvalidNameType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "", obj.TestName)
}

func testFormBindingInvalidName2(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := InvalidNameMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testQueryBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func testQueryBindingFail(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooStructForMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testQueryBindingBoolFail(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooStructForBoolType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

func TestCanSet(t *testing.T) {
	type CanSetStruct struct {
		lowerStart string `form:"lower"`
	}

	var c CanSetStruct
	assert.Nil(t, mapForm(&c, nil))
}
