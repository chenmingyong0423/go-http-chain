// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpchain

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func mapsEqual(t *testing.T, m1, m2 map[string][]string, ignoreLen bool) {
	if !ignoreLen {
		require.Equal(t, m1, m2)
	}
	for key, aValue := range m1 {
		bValue := m2[key]
		require.ElementsMatch(t, aValue, bValue)
	}
}

func TestNewWithClient(t *testing.T) {
	client := NewWithClient(http.DefaultClient)
	require.NotNil(t, client)
}

func TestNewDefault(t *testing.T) {
	client := NewDefault()
	require.NotNil(t, client)
}

func TestClient(t *testing.T) {
	client := NewDefault()
	client.SetHeader("X-Addr", "海南省").
		AddHeader("X-Name", "陈明勇").
		SetQuery("name", "陈明勇").
		AddQuery("addr", "海南省")
	mapsEqual(t, client.headers, http.Header{
		"X-Addr": {"海南省"},
		"X-Name": {"陈明勇"},
	}, false)
	mapsEqual(t, client.queryValues, url.Values{
		"addr": {"海南省"},
		"name": {"陈明勇"},
	}, false)
}

func TestClient_Request(t *testing.T) {
	client := NewDefault()
	client.SetHeader("X-Addr", "海南省").
		AddHeader("X-Name", "陈明勇").
		SetQuery("name", "陈明勇").
		AddQuery("addr", "海南省")
	r := client.Request("http://localhost:8080", http.MethodGet)
	compareRequest(t, r, http.MethodGet, client)
	g := client.Get("http://localhost:8080")
	compareRequest(t, g, http.MethodGet, client)
	p := client.Post("http://localhost:8080")
	compareRequest(t, p, http.MethodPost, client)
	put := client.Put("http://localhost:8080")
	compareRequest(t, put, http.MethodPut, client)
	d := client.Delete("http://localhost:8080")
	compareRequest(t, d, http.MethodDelete, client)
	patch := client.Patch("http://localhost:8080")
	compareRequest(t, patch, http.MethodPatch, client)
	head := client.Head("http://localhost:8080")
	compareRequest(t, head, http.MethodHead, client)
	opt := client.Options("http://localhost:8080")
	compareRequest(t, opt, http.MethodOptions, client)
	conn := client.Connect("http://localhost:8080")
	compareRequest(t, conn, http.MethodConnect, client)
	trace := client.Trace("http://localhost:8080")
	compareRequest(t, trace, http.MethodTrace, client)
}

func compareRequest(t *testing.T, r *Request, method string, client *Client) {
	require.NotNil(t, r)
	require.Equal(t, "http://localhost:8080", r.url)
	require.Equal(t, method, r.method)
	mapsEqual(t, r.headers, client.headers, false)
	mapsEqual(t, r.queryValues, client.queryValues, false)
	require.Equal(t, client.client, r.Client)
}
