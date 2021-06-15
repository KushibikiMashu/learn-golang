package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// GET /index
func TestIndexHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	IndexHandler(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got = %d, want = 200", resp.StatusCode)
	}

	actual, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("cannot read test response: %v", err)
	}

	if string(actual) != "hello world" {
		t.Errorf("got = %s, want = hello world", actual)
	}
}

// GET /echo
func TestEchoHandlerWithoutQuery(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/echo", nil)
	w := httptest.NewRecorder()

	EchoHandler(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got = %d, want = 200", resp.StatusCode)
	}

	actual, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("cannot read test response: %v", err)
	}

	if string(actual) != "no query" {
		t.Errorf("got = %s, want = no query", actual)
	}
}

// GET echo?q=hello
func TestEchoHandlerWithQuery(t *testing.T) {
	keyword := "hello"

	r := httptest.NewRequest(http.MethodGet, "/echo", nil)

	// クエリを組み立てる
	q := r.URL.Query()
	q.Set("q", keyword) // 追加するなら、 q.Add(key, value)
	r.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	EchoHandler(w, r)

	resp := w.Result()
	actual, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("cannot read test response: %v", err)
	}

	if string(actual) != keyword {
		t.Errorf("got = %s, want = %s", actual, keyword)
	}
}

func TestJsonGetHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/json", nil)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	JsonGetHandler(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got = %d, want = 200", resp.StatusCode)
	}

	// "application/json" という文字列が含まれていなければエラー
	if contentType := resp.Header.Get("Content-Type"); !strings.Contains(contentType, "application/json") {
		t.Error("Content Type: ", contentType)
	}

	actual, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("cannot read test response: %v", err)
	}

	type ResBody struct {
		// key 名は大文字にする
		Message string `json:"message"`
	}
	var resJson ResBody

	// actual（[]byte）を JSON の構造体に変換し、resJson に代入する
	if parseErr := json.Unmarshal(actual, &resJson); parseErr != nil {
		t.Errorf("parseErr: %v", parseErr)
	}

	if resJson.Message != "hello" {
		t.Errorf("got = %s, want = %s", resJson.Message, "hello")
	}
}

func TestJsonPostHandlerWithJsonBody(t *testing.T) {
	postBody := map[string]string{
		"name": "self",
	}
	postJson, err := json.Marshal(postBody)

	if err != nil {
		t.Errorf("Cannot make body json: %v", err)
	}

	r := httptest.NewRequest("POST", "/json", bytes.NewBuffer(postJson))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	JsonPostHandler(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode: got = %d, want %d", resp.StatusCode, 200)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Cannot read response body: %v", err)
	}

	type Actual struct {
		Result string `json:"result"`
	}
	var actual Actual

	if parseErr := json.Unmarshal(body, &actual); parseErr != nil {
		t.Errorf("Cannot parse json: %v", parseErr)
	}

	if result := actual.Result; result != "self" {
		t.Errorf("result: got = %v, want %v", result, "self")
	}
}
