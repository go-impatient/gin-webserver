package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewServerInformation(t *testing.T) {
	var c JSONObject
	c = make(JSONObject)
	if !c.Bool("a", true) {
		t.Fatalf("TestNewServerInformation failed")
	}
	if 1 != c.Int("a", 1) {
		t.Fatalf("TestNewServerInformation failed")
	}
	if "a" != c.String("a", "a") {
		t.Fatalf("TestNewServerInformation failed")
	}
	if nil == c.Object("a") {
		t.Fatalf("TestNewServerInformation failed")
	}

	c.Set("a", true)
	c.Set("b", 1)
	c.Set("c", "a")
	c.Object("d").Set("a", false)
	if !c.Bool("a", true) {
		t.Fatalf("TestNewServerInformation failed")
	}
	if 1 != c.Int("b", 1) {
		t.Fatalf("TestNewServerInformation failed")
	}
	if "a" != c.String("c", "a") {
		t.Fatalf("TestNewServerInformation failed")
	}
	if c.Object("d").Bool("a", true) {
		t.Fatalf("TestNewServerInformation failed")
	}

	if "a" != c.String("a", "a") {
		t.Fatalf("TestNewServerInformation failed")
	}
	if !c.Bool("b", true) {
		t.Fatalf("TestNewServerInformation failed")
	}
	if 1 != c.Int("c", 1) {
		t.Fatalf("TestNewServerInformation failed")
	}
	if "a" != c.String("d", "a") {
		t.Fatalf("TestNewServerInformation failed")
	}
	if nil == c.Object("a") {
		t.Fatalf("TestNewServerInformation failed")
	}

	c.Set(1, 1)
	c.Set(uint(2), 2)
	c.Set(1.2, 1.2)
	c.Set(1+1i, 0)
	c.Set(nil, nil)
	c.Set(make(map[string]string), "")
	c.Set(make([]string, 1), "")
	type a struct{}
	c.Set(a{}, "")
	c.Set(&a{}, "")

	b, _ := json.Marshal(c)
	fmt.Println(string(b))
}
