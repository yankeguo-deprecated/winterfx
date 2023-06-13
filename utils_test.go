package winterfx

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func TestFlattenSimpleSlice(t *testing.T) {
	require.Equal(t, "a", flattenSingleSlice([]string{"a"}))
	require.Equal(t, []int{1, 2}, flattenSingleSlice([]int{1, 2}))
}

func TestExtractRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/get?aaa=bbb", nil)

	fm := map[string][]*multipart.FileHeader{}

	m := map[string]any{}
	err := extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "query_aaa": "bbb"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`{"hello":"world"}`)))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "header_content_type": "application/json;charset=utf-8", "hello": "world", "query_aaa": "bbb"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "header_content_type": "application/x-www-form-urlencoded;charset=utf-8", "hello": "world", "query_aaa": "bbb"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "text/plain;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "header_content_type": "text/plain;charset=utf-8", "query_aaa": "bbb", "body": "hello=world"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "application/x-custom")

	m = map[string]any{}
	err = extractRequest(m, fm, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "header_content_type": "application/x-custom", "query_aaa": "bbb", "body": []byte("hello=world")}, m)
}