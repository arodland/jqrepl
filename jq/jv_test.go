package jq

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestJvKind(t *testing.T) {
	is := is.New(t)

	cases := []struct {
		*Jv
		JvKind
		string
	}{
		{JvNull(), JV_KIND_NULL, "null"},
		{JvFromString("a"), JV_KIND_STRING, "string"},
	}

	for _, c := range cases {
		defer c.Free()
		is.Equal(c.Kind(), c.JvKind)
		is.Equal(c.Kind().String(), c.string)
	}
}

func TestJvString(t *testing.T) {
	is := is.New(t)

	jv := JvFromString("test")
	defer jv.Free()

	str, err := jv.String()

	is.Equal(str, "test")
	is.NoErr(err)

	i := jv.ToGoVal()

	is.Equal(i, "test")
}

func TestJvStringOnNonStringType(t *testing.T) {
	is := is.New(t)

	// Test that on a non-string value we get a go error, not a C assert
	jv := JvNull()
	defer jv.Free()

	_, err := jv.String()
	is.Err(err)
}

func TestJvFromJSONString(t *testing.T) {
	is := is.New(t)

	jv, err := JvFromJSONString("[]")
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_ARRAY)

	jv, err = JvFromJSONString("not valid")
	is.Err(err)
	is.Nil(jv)
}

func TestJvFromFloat(t *testing.T) {
	is := is.New(t)

	jv := JvFromFloat(1.23)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_NUMBER)
	gv := jv.ToGoVal()
	n, ok := gv.(float64)
	is.True(ok)
	is.Equal(n, float64(1.23))
}

func TestJvFromInterface(t *testing.T) {
	is := is.New(t)

	// Null
	jv, err := JvFromInterface(nil)
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_NULL)

	// Boolean
	jv, err = JvFromInterface(true)
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_TRUE)

	jv, err = JvFromInterface(false)
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_FALSE)

	// Float
	jv, err = JvFromInterface(1.23)
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_NUMBER)
	gv := jv.ToGoVal()
	n, ok := gv.(float64)
	is.True(ok)
	is.Equal(n, float64(1.23))

	// Integer
	jv, err = JvFromInterface(456)
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_NUMBER)
	gv = jv.ToGoVal()
	n2, ok := gv.(int)
	is.True(ok)
	is.Equal(n2, 456)

	// String
	jv, err = JvFromInterface("test")
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_STRING)
	gv = jv.ToGoVal()
	s, ok := gv.(string)
	is.True(ok)
	is.Equal(s, "test")

	jv, err = JvFromInterface([]string{"test", "one", "two"})
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_ARRAY)
	gv = jv.ToGoVal()
	is.Equal(gv.([]interface{})[2], "two")

	jv, err = JvFromInterface(map[string]int{"one": 1, "two": 2})
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), JV_KIND_OBJECT)
	gv = jv.ToGoVal()
	is.Equal(gv.(map[string]interface{})["two"], 2)
}

func TestJvDump(t *testing.T) {
	is := is.New(t)

	jv := JvFromString("test")
	defer jv.Free()

	dump := jv.Copy().Dump(JvPrintNone)

	is.Equal(`"test"`, dump)
	dump = jv.Copy().Dump(JvPrintColour)

	is.Equal([]byte("\x1b[0;32m"+`"test"`+"\x1b[0m"), []byte(dump))
}

func TestJvInvalid(t *testing.T) {
	is := is.New(t)

	jv := JvInvalid()

	is.False(jv.IsValid())

	_, ok := jv.Copy().GetInvalidMessageAsString()
	is.False(ok) // "Expected no Invalid message"

	jv = jv.GetInvalidMessage()
	is.Equal(jv.Kind(), JV_KIND_NULL)
}

func TestJvInvalidWithMessage_string(t *testing.T) {
	is := is.New(t)

	jv := JvInvalidWithMessage(JvFromString("Error message 1"))

	is.False(jv.IsValid())

	msg := jv.Copy().GetInvalidMessage()
	is.Equal(msg.Kind(), JV_KIND_STRING)
	msg.Free()

	str, ok := jv.GetInvalidMessageAsString()
	is.True(ok)
	is.Equal("Error message 1", str)
}

func TestJvInvalidWithMessage_object(t *testing.T) {
	is := is.New(t)

	jv := JvInvalidWithMessage(JvObject())

	is.False(jv.IsValid())

	msg := jv.Copy().GetInvalidMessage()
	is.Equal(msg.Kind(), JV_KIND_OBJECT)
	msg.Free()

	str, ok := jv.GetInvalidMessageAsString()
	is.True(ok)
	is.Equal("{}", str)

}
