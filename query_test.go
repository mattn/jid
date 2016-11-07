package jig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	var assert = assert.New(t)

	assert.True(validate([]rune(".test.name")))
	assert.True(validate([]rune(".test.name.")))
	assert.True(validate([]rune(".test[0].name.")))
	assert.True(validate([]rune(".[0].name.")))
	assert.True(validate([]rune(".name[9][1]")))
	assert.True(validate([]rune(".[0][1].name.")))

	assert.False(validate([]rune("[0].name.")))
	assert.False(validate([]rune(".test[0]].name.")))
	assert.False(validate([]rune(".test..name")))
	assert.False(validate([]rune(".test.name..")))
	assert.False(validate([]rune(".test[[0]].name.")))
	assert.False(validate([]rune(".test[0]name.")))
	assert.False(validate([]rune(".test.[0].name.")))
}

func TestNewQuery(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".name")
	q := NewQuery(v)

	assert.Equal(*q.query, []rune(".name"))
	assert.Equal(*q.complete, []rune(""))
}

func TestNewQueryWithInvalidQuery(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("name")
	q := NewQuery(v)

	assert.Equal(*q.query, []rune(""))
	assert.Equal(*q.complete, []rune(""))
}

func TestNewQueryWithString(t *testing.T) {
	var assert = assert.New(t)

	q := NewQueryWithString(".name")

	assert.Equal(*q.query, []rune(".name"))
	assert.Equal(*q.complete, []rune(""))
}

func TestNewQueryWithStringWithInvalidQuery(t *testing.T) {
	var assert = assert.New(t)

	q := NewQueryWithString("name")

	assert.Equal(*q.query, []rune(""))
	assert.Equal(*q.complete, []rune(""))
}

func TestQueryGet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test")
	q := NewQuery(v)

	assert.Equal(q.Get(), []rune(".test"))
}

func TestQuerySet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".hello")
	q := NewQuery(v)

	assert.Equal([]rune(".world"), q.Set([]rune(".world")))
	assert.Equal("", string(q.Set([]rune(""))))
}

func TestQuerySetWithInvalidQuery(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".hello")
	q := NewQuery(v)

	assert.Equal(q.Set([]rune("world")), []rune(".hello"))
}

func TestQueryAdd(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".hello")
	q := NewQuery(v)

	assert.Equal(q.Add([]rune("world")), []rune(".helloworld"))
}
func TestQueryInsert(t *testing.T) {
	var assert = assert.New(t)
	v := []rune(".hello.world")
	q := NewQuery(v)

	assert.Equal([]rune(".hello.world"), q.Insert([]rune("w"), 0))
	assert.Equal([]rune(".whello.world"), q.Insert([]rune("w"), 1))
	assert.Equal([]rune(".wwhello.world"), q.Insert([]rune("w"), 1))
	assert.Equal([]rune(".wwhello.world"), q.Insert([]rune("."), 1))
	assert.Equal([]rune(".wwh.ello.world"), q.Insert([]rune("."), 4))
	assert.Equal([]rune(".wwh.ello.world"), q.Insert([]rune("a"), 20))
}
func TestQueryStringInsert(t *testing.T) {
	var assert = assert.New(t)
	q := NewQueryWithString(".hello.world")

	assert.Equal(".hello.world", q.StringInsert("w", 0))
	assert.Equal(".whello.world", q.StringInsert("w", 1))
	assert.Equal(".wwhello.world", q.StringInsert("w", 1))
	assert.Equal(".wwhello.world", q.StringInsert(".", 1))
	assert.Equal(".wwh.ello.world", q.StringInsert(".", 4))
	assert.Equal(".wwh.ello.worlda", q.StringInsert("a", 15))
	assert.Equal(".wwh.ello.worlda", q.StringInsert("a", 20))
}

func TestQueryClear(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test")
	q := NewQuery(v)

	assert.Equal(q.Clear(), []rune(""))
}

func TestQueryDelete(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".helloworld")
	q := NewQuery(v)

	assert.Equal([]rune(".helloworl"), q.Delete(-1))
	assert.Equal([]rune(".hellowor"), q.Delete(-1))
	assert.Equal([]rune(".hellow"), q.Delete(-2))
	assert.Equal([]rune(""), q.Delete(-8))

	q = NewQuery([]rune(".hello.world"))
	assert.Equal([]rune(".hello.world"), q.Delete(0))
	assert.Equal([]rune(".ello.world"), q.Delete(1))
	assert.Equal([]rune(".llo.world"), q.Delete(1))
	assert.Equal([]rune(".llo.world"), q.Delete(0))
	assert.Equal([]rune(".ll.world"), q.Delete(3))
	assert.Equal([]rune(".llworld"), q.Delete(3))
	assert.Equal([]rune(".llorld"), q.Delete(3))
}

func TestGetKeywords(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test.name")
	q := NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("test"),
		[]rune("name"),
	})

	v = []rune("")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{})

	v = []rune(".test.name.")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("test"),
		[]rune("name"),
		[]rune(""),
	})

	v = []rune(".hello")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("hello"),
	})

	v = []rune(".hello.")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("hello"),
		[]rune(""),
	})

	v = []rune(".hello[")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("hello"),
		[]rune("["),
	})

	v = []rune(".hello[12")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("hello"),
		[]rune("[12"),
	})

	v = []rune(".hello[0]")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("hello"),
		[]rune("[0]"),
	})

	v = []rune(".hello[13][0]")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("hello"),
		[]rune("[13]"),
		[]rune("[0]"),
	})

	v = []rune(".[3][23].hello[13][0]")
	q = NewQuery(v)
	assert.Equal(q.GetKeywords(), [][]rune{
		[]rune("[3]"),
		[]rune("[23]"),
		[]rune("hello"),
		[]rune("[13]"),
		[]rune("[0]"),
	})

}

func TestGetLastKeyword(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test.name")
	q := NewQuery(v)
	assert.Equal(q.GetLastKeyword(), []rune("name"))

	v = []rune(".test.")
	q = NewQuery(v)
	assert.Equal(q.GetLastKeyword(), []rune(""))

	v = []rune(".test")
	q = NewQuery(v)
	assert.Equal(q.GetLastKeyword(), []rune("test"))
}

func TestStringGetLastKeyword(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test.name")
	q := NewQuery(v)
	assert.Equal(q.StringGetLastKeyword(), "name")

	v = []rune(".test.")
	q = NewQuery(v)
	assert.Equal(q.StringGetLastKeyword(), "")

	v = []rune(".test")
	q = NewQuery(v)
	assert.Equal(q.StringGetLastKeyword(), "test")

	v = []rune("")
	q = NewQuery(v)
	assert.Equal(q.StringGetLastKeyword(), "")
}

func TestPopKeyword(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test.name")
	q := NewQuery(v)
	k, query := q.PopKeyword()
	assert.Equal(k, []rune("name"))
	assert.Equal(query, []rune(".test"))
	assert.Equal(q.Get(), []rune(".test"))

	v = []rune(".test.name.")
	q = NewQuery(v)
	k, query = q.PopKeyword()
	assert.Equal(k, []rune(""))
	assert.Equal(query, []rune(".test.name"))
	assert.Equal(q.Get(), []rune(".test.name"))
}

func TestQueryStringGet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test")
	q := NewQuery(v)

	assert.Equal(q.StringGet(), ".test")
}

func TestQueryStringSet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".hello")
	q := NewQuery(v)

	assert.Equal(q.StringSet(".world"), ".world")
}

func TestQueryStringAdd(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".hello")
	q := NewQuery(v)

	assert.Equal(q.StringAdd("world"), ".helloworld")
}

func TestStringGetKeywords(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test.name")
	q := NewQuery(v)
	assert.Equal(q.StringGetKeywords(), []string{
		"test",
		"name",
	})

	v = []rune(".test.name")
	q = NewQuery(v)
	assert.Equal(q.StringGetKeywords(), []string{
		"test",
		"name",
	})

	v = []rune("")
	q = NewQuery(v)
	kws := q.StringGetKeywords()
	assert.Equal([]string(nil), kws)
	assert.Equal(0, len(kws))
}

func TestStringPopKeyword(t *testing.T) {
	var assert = assert.New(t)

	v := []rune(".test.name")
	q := NewQuery(v)
	k, query := q.StringPopKeyword()
	assert.Equal(k, "name")
	assert.Equal(query, []rune(".test"))
	assert.Equal(q.Get(), []rune(".test"))

	v = []rune(".test.name.")
	q = NewQuery(v)
	k, query = q.StringPopKeyword()
	assert.Equal(k, "")
	assert.Equal(query, []rune(".test.name"))
	assert.Equal(q.Get(), []rune(".test.name"))

	v = []rune(".test.name[23]")
	q = NewQuery(v)
	k, query = q.StringPopKeyword()
	assert.Equal(k, "[23]")
	assert.Equal(query, []rune(".test.name"))
	assert.Equal(q.Get(), []rune(".test.name"))
}
