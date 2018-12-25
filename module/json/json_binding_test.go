// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"encoding/json"
	"github.com/alimy/gin/binding"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required"`
}

type FooStructUseNumber struct {
	Foo interface{} `json:"foo" binding:"required"`
}

type FooStructForMapType struct {
	// Unknown type: not support map
	MapFoo map[string]interface{} `form:"map_foo"`
}

type FooStructForSliceType struct {
	SliceFoo []int `form:"slice_foo"`
}

type FooStructForStructType struct {
	StructFoo struct {
		Idx int `form:"idx"`
	}
}

type FooStructForStructPointerType struct {
	StructPointerFoo *struct {
		Name string `form:"name"`
	}
}

type FooStructForSliceMapType struct {
	// Unknown type: not support map
	SliceMapFoo []map[string]interface{} `form:"slice_map_foo"`
}

type FooBarStructForIntType struct {
	IntFoo int `form:"int_foo"`
	IntBar int `form:"int_bar" binding:"required"`
}

type FooBarStructForInt8Type struct {
	Int8Foo int8 `form:"int8_foo"`
	Int8Bar int8 `form:"int8_bar" binding:"required"`
}

type FooBarStructForInt16Type struct {
	Int16Foo int16 `form:"int16_foo"`
	Int16Bar int16 `form:"int16_bar" binding:"required"`
}

type FooBarStructForInt32Type struct {
	Int32Foo int32 `form:"int32_foo"`
	Int32Bar int32 `form:"int32_bar" binding:"required"`
}

type FooBarStructForInt64Type struct {
	Int64Foo int64 `form:"int64_foo"`
	Int64Bar int64 `form:"int64_bar" binding:"required"`
}

type FooBarStructForUintType struct {
	UintFoo uint `form:"uint_foo"`
	UintBar uint `form:"uint_bar" binding:"required"`
}

type FooBarStructForUint8Type struct {
	Uint8Foo uint8 `form:"uint8_foo"`
	Uint8Bar uint8 `form:"uint8_bar" binding:"required"`
}

type FooBarStructForUint16Type struct {
	Uint16Foo uint16 `form:"uint16_foo"`
	Uint16Bar uint16 `form:"uint16_bar" binding:"required"`
}

type FooBarStructForUint32Type struct {
	Uint32Foo uint32 `form:"uint32_foo"`
	Uint32Bar uint32 `form:"uint32_bar" binding:"required"`
}

type FooBarStructForUint64Type struct {
	Uint64Foo uint64 `form:"uint64_foo"`
	Uint64Bar uint64 `form:"uint64_bar" binding:"required"`
}

type FooBarStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
	BoolBar bool `form:"bool_bar" binding:"required"`
}

type FooBarStructForFloat32Type struct {
	Float32Foo float32 `form:"float32_foo"`
	Float32Bar float32 `form:"float32_bar" binding:"required"`
}

type FooBarStructForFloat64Type struct {
	Float64Foo float64 `form:"float64_foo"`
	Float64Bar float64 `form:"float64_bar" binding:"required"`
}

type FooStructForStringPtrType struct {
	PtrFoo *string `form:"ptr_foo"`
	PtrBar *string `form:"ptr_bar" binding:"required"`
}

type FooStructForMapPtrType struct {
	PtrBar *map[string]interface{} `form:"ptr_bar"`
}

func TestBindingDefault(t *testing.T) {
	assert.Equal(t, &jsonBinding{}, binding.DefaultWith("POST", binding.MIMEJSON))
	assert.Equal(t, &jsonBinding{}, binding.DefaultWith("PUT", binding.MIMEJSON))
}

func TestBindingJSONNilBody(t *testing.T) {
	var obj FooStruct
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := (&jsonBinding{}).Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingJSON(t *testing.T) {
	testBodyBinding(t,
		jsonBinding{}, "json",
		"/", "/",
		`{"foo": "bar"}`, `{"bar": "foo"}`)
}

func TestBindingJSONUseNumber(t *testing.T) {
	testBodyBindingUseNumber(t,
		jsonBinding{}, "json",
		"/", "/",
		`{"foo": 123}`, `{"bar": "foo"}`)
}

func TestBindingJSONUseNumber2(t *testing.T) {
	testBodyBindingUseNumber2(t,
		jsonBinding{}, "json",
		"/", "/",
		`{"foo": 123}`, `{"bar": "foo"}`)
}

func TestValidationFails(t *testing.T) {
	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := (&jsonBinding{}).Bind(req, &obj)
	assert.Error(t, err)
}

func TestValidationDisabled(t *testing.T) {
	backup := binding.Validator
	binding.Validator = nil
	defer func() { binding.Validator = backup }()

	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := (&jsonBinding{}).Bind(req, &obj)
	assert.NoError(t, err)
}

func TestExistsSucceeds(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"hoge" binding:"exists"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"hoge": 0}`)
	err := (&jsonBinding{}).Bind(req, &obj)
	assert.NoError(t, err)
}

func TestExistsFails(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"foo" binding:"exists"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"boen": 0}`)
	err := (&jsonBinding{}).Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingFormForType(t *testing.T) {
	testFormBindingForType(t, "POST",
		"/", "/",
		"map_foo=", "bar2=1", "Map")

	testFormBindingForType(t, "POST",
		"/", "/",
		"slice_foo=1&slice_foo=2", "bar2=1&bar2=2", "Slice")

	testFormBindingForType(t, "GET",
		"/?slice_foo=1&slice_foo=2", "/?bar2=1&bar2=2",
		"", "", "Slice")

	testFormBindingForType(t, "POST",
		"/", "/",
		"slice_map_foo=1&slice_map_foo=2", "bar2=1&bar2=2", "SliceMap")

	testFormBindingForType(t, "GET",
		"/?slice_map_foo=1&slice_map_foo=2", "/?bar2=1&bar2=2",
		"", "", "SliceMap")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int_foo=&int_bar=-12", "bar2=-123", "Int")

	testFormBindingForType(t, "GET",
		"/?int_foo=&int_bar=-12", "/?bar2=-123",
		"", "", "Int")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int8_foo=&int8_bar=-12", "bar2=-123", "Int8")

	testFormBindingForType(t, "GET",
		"/?int8_foo=&int8_bar=-12", "/?bar2=-123",
		"", "", "Int8")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int16_foo=&int16_bar=-12", "bar2=-123", "Int16")

	testFormBindingForType(t, "GET",
		"/?int16_foo=&int16_bar=-12", "/?bar2=-123",
		"", "", "Int16")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int32_foo=&int32_bar=-12", "bar2=-123", "Int32")

	testFormBindingForType(t, "GET",
		"/?int32_foo=&int32_bar=-12", "/?bar2=-123",
		"", "", "Int32")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int64_foo=&int64_bar=-12", "bar2=-123", "Int64")

	testFormBindingForType(t, "GET",
		"/?int64_foo=&int64_bar=-12", "/?bar2=-123",
		"", "", "Int64")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint_foo=&uint_bar=12", "bar2=123", "Uint")

	testFormBindingForType(t, "GET",
		"/?uint_foo=&uint_bar=12", "/?bar2=123",
		"", "", "Uint")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint8_foo=&uint8_bar=12", "bar2=123", "Uint8")

	testFormBindingForType(t, "GET",
		"/?uint8_foo=&uint8_bar=12", "/?bar2=123",
		"", "", "Uint8")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint16_foo=&uint16_bar=12", "bar2=123", "Uint16")

	testFormBindingForType(t, "GET",
		"/?uint16_foo=&uint16_bar=12", "/?bar2=123",
		"", "", "Uint16")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint32_foo=&uint32_bar=12", "bar2=123", "Uint32")

	testFormBindingForType(t, "GET",
		"/?uint32_foo=&uint32_bar=12", "/?bar2=123",
		"", "", "Uint32")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint64_foo=&uint64_bar=12", "bar2=123", "Uint64")

	testFormBindingForType(t, "GET",
		"/?uint64_foo=&uint64_bar=12", "/?bar2=123",
		"", "", "Uint64")

	testFormBindingForType(t, "POST",
		"/", "/",
		"bool_foo=&bool_bar=true", "bar2=true", "Bool")

	testFormBindingForType(t, "GET",
		"/?bool_foo=&bool_bar=true", "/?bar2=true",
		"", "", "Bool")

	testFormBindingForType(t, "POST",
		"/", "/",
		"float32_foo=&float32_bar=-12.34", "bar2=12.3", "Float32")

	testFormBindingForType(t, "GET",
		"/?float32_foo=&float32_bar=-12.34", "/?bar2=12.3",
		"", "", "Float32")

	testFormBindingForType(t, "POST",
		"/", "/",
		"float64_foo=&float64_bar=-12.34", "bar2=12.3", "Float64")

	testFormBindingForType(t, "GET",
		"/?float64_foo=&float64_bar=-12.34", "/?bar2=12.3",
		"", "", "Float64")

	testFormBindingForType(t, "POST",
		"/", "/",
		"ptr_bar=test", "bar2=test", "Ptr")

	testFormBindingForType(t, "GET",
		"/?ptr_bar=test", "/?bar2=test",
		"", "", "Ptr")

	testFormBindingForType(t, "POST",
		"/", "/",
		"idx=123", "id1=1", "Struct")

	testFormBindingForType(t, "GET",
		"/?idx=123", "/?id1=1",
		"", "", "Struct")

	testFormBindingForType(t, "POST",
		"/", "/",
		"name=thinkerou", "name1=ou", "StructPointer")

	testFormBindingForType(t, "GET",
		"/?name=thinkerou", "/?name1=ou",
		"", "", "StructPointer")
}

func testFormBindingForType(t *testing.T, method, path, badPath, body, badBody string, typ string) {
	b := binding.Form
	assert.Equal(t, "form", b.Name())

	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", binding.MIMEPOSTForm)
	}
	switch typ {
	case "Int":
		obj := FooBarStructForIntType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int(0), obj.IntFoo)
		assert.Equal(t, int(-12), obj.IntBar)

		obj = FooBarStructForIntType{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Int8":
		obj := FooBarStructForInt8Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int8(0), obj.Int8Foo)
		assert.Equal(t, int8(-12), obj.Int8Bar)

		obj = FooBarStructForInt8Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Int16":
		obj := FooBarStructForInt16Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int16(0), obj.Int16Foo)
		assert.Equal(t, int16(-12), obj.Int16Bar)

		obj = FooBarStructForInt16Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Int32":
		obj := FooBarStructForInt32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int32(0), obj.Int32Foo)
		assert.Equal(t, int32(-12), obj.Int32Bar)

		obj = FooBarStructForInt32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Int64":
		obj := FooBarStructForInt64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), obj.Int64Foo)
		assert.Equal(t, int64(-12), obj.Int64Bar)

		obj = FooBarStructForInt64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Uint":
		obj := FooBarStructForUintType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint(0x0), obj.UintFoo)
		assert.Equal(t, uint(0xc), obj.UintBar)

		obj = FooBarStructForUintType{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Uint8":
		obj := FooBarStructForUint8Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x0), obj.Uint8Foo)
		assert.Equal(t, uint8(0xc), obj.Uint8Bar)

		obj = FooBarStructForUint8Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Uint16":
		obj := FooBarStructForUint16Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0x0), obj.Uint16Foo)
		assert.Equal(t, uint16(0xc), obj.Uint16Bar)

		obj = FooBarStructForUint16Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Uint32":
		obj := FooBarStructForUint32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x0), obj.Uint32Foo)
		assert.Equal(t, uint32(0xc), obj.Uint32Bar)

		obj = FooBarStructForUint32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Uint64":
		obj := FooBarStructForUint64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0x0), obj.Uint64Foo)
		assert.Equal(t, uint64(0xc), obj.Uint64Bar)

		obj = FooBarStructForUint64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Float32":
		obj := FooBarStructForFloat32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, float32(0.0), obj.Float32Foo)
		assert.Equal(t, float32(-12.34), obj.Float32Bar)

		obj = FooBarStructForFloat32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Float64":
		obj := FooBarStructForFloat64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, float64(0.0), obj.Float64Foo)
		assert.Equal(t, float64(-12.34), obj.Float64Bar)

		obj = FooBarStructForFloat64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Bool":
		obj := FooBarStructForBoolType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.False(t, obj.BoolFoo)
		assert.True(t, obj.BoolBar)

		obj = FooBarStructForBoolType{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Slice":
		obj := FooStructForSliceType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2}, obj.SliceFoo)

		obj = FooStructForSliceType{}
		req = requestWithBody(method, badPath, badBody)
		err = (&jsonBinding{}).Bind(req, &obj)
		assert.Error(t, err)
	case "Struct":
		obj := FooStructForStructType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t,
			struct {
				Idx int "form:\"idx\""
			}(struct {
				Idx int "form:\"idx\""
			}{Idx: 123}),
			obj.StructFoo)
	case "StructPointer":
		obj := FooStructForStructPointerType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t,
			struct {
				Name string "form:\"name\""
			}(struct {
				Name string "form:\"name\""
			}{Name: "thinkerou"}),
			*obj.StructPointerFoo)
	case "Map":
		obj := FooStructForMapType{}
		err := b.Bind(req, &obj)
		assert.Error(t, err)
	case "SliceMap":
		obj := FooStructForSliceMapType{}
		err := b.Bind(req, &obj)
		assert.Error(t, err)
	case "Ptr":
		obj := FooStructForStringPtrType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Nil(t, obj.PtrFoo)
		assert.Equal(t, "test", *obj.PtrBar)

		obj = FooStructForStringPtrType{}
		obj.PtrBar = new(string)
		err = b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, "test", *obj.PtrBar)

		objErr := FooStructForMapPtrType{}
		err = b.Bind(req, &objErr)
		assert.Error(t, err)

		obj = FooStructForStringPtrType{}
		req = requestWithBody(method, badPath, badBody)
		err = b.Bind(req, &obj)
		assert.Error(t, err)
	}
}

func testBodyBinding(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = (&jsonBinding{}).Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingUseNumber(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = true
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// we hope it is int64(123)
	v, e := obj.Foo.(json.Number).Int64()
	assert.NoError(t, e)
	assert.Equal(t, int64(123), v)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = (&jsonBinding{}).Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingUseNumber2(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = false
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// it will return float64(123) if not use EnableDecoderUseNumber
	// maybe it is not hoped
	assert.Equal(t, float64(123), obj.Foo)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = (&jsonBinding{}).Bind(req, &obj)
	assert.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
