package cache_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/cache"
)

func TestCache(t *testing.T) {
	c := cache.NewMemoryCache()

	var err error

	_, err = c.Get("a")
	if err != cache.ErrNotFound {
		t.Errorf("want: %s, got: %s", cache.ErrNotFound, err)
	}

	if c.Exists("a") {
		t.Errorf("key 'a' should not exist")
	}

	err = c.Put("a", 1)
	if err != nil {
		t.Errorf("error was not expected at this point. err:%s", err)
	}

	if !c.Exists("a") {
		t.Errorf("key 'a' should exist")
	}

	x, err := c.Get("a")
	if err != nil {
		t.Errorf("error was not expected at this point. err: %s", err)
	}
	if x != 1 {
		t.Errorf("want: %v, got: %v", 1, x)
	}

	err = c.Invalidate("a")
	if err != nil {
		t.Errorf("error was not expected at this point. err: %s", err)
	}
}

func TestCacheWithDifferentValueTypes(t *testing.T) {
	c := cache.NewMemoryCache()

	var testCase = []struct {
		key   string
		value interface{}
	}{
		{key: "integer", value: 1},
		{key: "float", value: 3.5},
		{key: "string", value: "test"},
	}

	for _, tc := range testCase {
		err := c.Put(tc.key, tc.value)
		if err != nil {
			t.Errorf("error was not expected at this point. key: %s, err: %s", tc.key, err)
		}
	}

	for _, tc := range testCase {
		got, err := c.Get(tc.key)
		if err != nil {
			t.Errorf("error was not expected at this point. key: %s, err: %s", tc.key, err)
		}
		if got != tc.value {
			t.Errorf("want: %v, got: %v", tc.value, got)
		}
	}
}
