package storage

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testObject = struct {
	Message string `json:"message"`
}{"the power you have is to be the best version of yourself you can be, so you can create a better world"}
var key = "ashley rickards"

func TestMain(m *testing.M) {

	db, closer, err := openBitcaskDb()
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(testObject)
	if err != nil {
		panic(err)
	}
	err = db.Put([]byte(key), bytes)
	if err != nil {
		panic(err)
	}
	closer()

	//Run the tests
	m.Run()

	//clean resources.
	db, closer, err = openBitcaskDb()
	if err != nil {
		panic(err)
	}
	db.Delete([]byte(key))
	closer()
}

//Tests if we could get the object by key and deserialize it properly
func Test_Get(t *testing.T) {
	var b = Bitcask{}
	var result struct {
		Message string `json:"message"`
	}

	err := b.Get(key, &result)
	if err != nil {
		t.Error(err)
	}

	if result != testObject {
		t.Errorf("expected %s, got %s", testObject, result)
	}
}

//Tests the function storage.Set()
func Test_Set(t *testing.T) {
	var b = Bitcask{}
	var obj = struct {
		Message string `json:"message"`
	}{"let men decide firmly what they will not do, and they will be free to do vigorously what they ought to do"}

	err := b.Set("mencius", obj)
	if err != nil {
		t.Error(err)
	}

	var result struct {
		Message string `json:"message"`
	}

	err = b.Get("mencius", &result)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(result, obj); diff != "" {
		t.Error(diff)
	}

	b.Remove("mencius")
}

//Tests the function Remove()
func Test_Remove(t *testing.T) {

	var b = Bitcask{}

	//1. Set the key
	var obj = struct {
		Message string `json:"message"`
	}{"to be successful you must accept all challenges that come your way. you can't just accept the ones you like"}

	err := b.Set("mike kafka", obj)
	if err != nil {
		t.Error(err)
	}

	//2. Fail if we can't find the key
	var res struct {
		Message string `json:"message"`
	}
	err = b.Get("mike kafka", &res)
	if err != nil {
		t.Error(err)
	}

	//2. Remove the key - fail if err is not nil
	err = b.Remove("mike kafka")
	if err != nil {
		t.Error(err)
	}

	//3. Check if we can get the key - fail if we can
	err = b.Get("mike kafka", &res)
	if err == nil {
		t.Error("able to get the key after removing it.")
	}
}

//Tests if storage.GetKeys() works
func Test_GetKeys(t *testing.T) {

	var keys = []string{"mike kafka", "mencius", "ashley rickards"}

	var obj = struct {
		Message string `json:"message"`
	}{"a key to success is self-confidence. a key to self-confidence is preparation"}

	var b = Bitcask{}

	//1. Set a couple of keys
	for _, k := range keys {
		err := b.Set(k, obj)
		if err != nil {
			t.Error("couldn't set the key", err)
		}
	}

	//2. Tests if we could get the keys
	result, err := b.GetKeys()
	if err != nil {
		t.Error(err)
	}

	//3. Iterate over the keys that we've set in the previous step.
	for _, k := range keys {
		if !contains(*result, k) {
			t.Errorf("the key %s does not appear in the Getkeys result", k)
		}
	}

	//4. Clean up keys
	for _, k := range keys {
		err := b.Remove(k)
		if err != nil {
			t.Errorf("couldn't remove the key %s", k)
		}
	}
}

//returns true if given string appears in the given array at least once
func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
