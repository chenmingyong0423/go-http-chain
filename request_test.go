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
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	r := NewDefault().Request("http://localhost:8080", "GET")
	r.SetHeader("X-Addr", "海南省").
		AddHeader("X-Name", "陈明勇").
		SetQuery("name", "陈明勇").
		AddQuery("addr", "海南省").
		SetBody(map[string]string{"name": "陈明勇", "addr": "海南省"}).
		SetBodyEncodeFunc(func(body any) (io.Reader, error) {
			return nil, nil
		})
	mapsEqual(t, r.headers, http.Header{
		"X-Addr": {"海南省"},
		"X-Name": {"陈明勇"},
	}, false)
	mapsEqual(t, r.queryValues, url.Values{
		"addr": {"海南省"},
		"name": {"陈明勇"},
	}, false)
	require.Equal(t, r.method, "GET")
	require.Equal(t, r.url, "http://localhost:8080")
	require.Equal(t, r.body, map[string]string{"name": "陈明勇", "addr": "海南省"})
	require.NotNil(t, r.bodyEncodeFunc)

	r.SetBody(&bytes.Buffer{})
	require.Equal(t, r.bodyBytes, &bytes.Buffer{})
}

func TestRequest_Do(t *testing.T) {
	// 启动一个 web 服务器
	go runWebServer(t)

	time.Sleep(500 * time.Millisecond)
	t.Run("test GET", func(t *testing.T) {
		client := NewDefault()
		getReq := client.AddHeader("X-Name", "陈明勇").
			AddQuery("name", "陈明勇").Get("http://localhost:8080/test")
		resp, err := getReq.Do(context.Background())
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, 200)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, string(body), "hello world")
	})
	t.Run("test POST", func(t *testing.T) {
		client := NewDefault()
		postReq := client.Post("http://localhost:8080/test").SetBody(map[string]string{"name": "陈明勇"}).SetBodyEncodeFunc(func(body any) (io.Reader, error) {
			data, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(data), nil
		})
		resp, err := postReq.Do(context.Background())
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, 200)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, string(body), "hello world")
	})

	t.Run("test request error", func(t *testing.T) {
		client := NewDefault()

		deleteReq := client.Delete("http://localhost:8080/test")
		resp, err := deleteReq.Do(nil)
		require.Nil(t, resp)
		require.Equal(t, err, errors.New("net/http: nil Context"))
	})
	t.Run("test Do error", func(t *testing.T) {
		client := NewDefault()

		deleteReq := client.Delete("")
		resp, err := deleteReq.Do(context.Background())
		require.Nil(t, resp)
		var urlErr = &url.Error{}
		require.True(t, errors.As(err, &urlErr))
	})

}

var once sync.Once

func runWebServer(t *testing.T) {
	once.Do(func() {
		http.HandleFunc("GET /test", func(writer http.ResponseWriter, request *http.Request) {
			mapsEqual(t, http.Header{
				"X-Name": {"陈明勇"},
			}, request.Header, true)
			mapsEqual(t, request.URL.Query(), url.Values{
				"name": {"陈明勇"},
			}, true)
			_, _ = writer.Write([]byte("hello world"))
		})
		http.HandleFunc("POST /test", func(writer http.ResponseWriter, request *http.Request) {
			body, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			mp := make(map[string]string)
			err = json.Unmarshal(body, &mp)
			require.NoError(t, err)
			require.Equal(t, mp, map[string]string{"name": "陈明勇"})
			_, _ = writer.Write([]byte("hello world"))
		})
		http.HandleFunc("PUT /test", func(writer http.ResponseWriter, request *http.Request) {
			body, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			mp := make(map[string]string)
			err = json.Unmarshal(body, &mp)
			require.NoError(t, err)
			require.Equal(t, mp, map[string]string{"name": "陈明勇"})
			result := map[string]any{"name": "陈明勇"}
			data, err := json.Marshal(result)
			require.NoError(t, err)
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(data)
		})
		http.HandleFunc("DELETE /test", func(writer http.ResponseWriter, request *http.Request) {
			result := map[string]any{"name": "陈明勇"}
			data, err := json.Marshal(result)
			require.NoError(t, err)
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = writer.Write(data)
		})
		http.HandleFunc("GET /xml", func(writer http.ResponseWriter, request *http.Request) {
			type XmlStruct struct {
				Name string `xml:"name"`
			}
			result := XmlStruct{
				Name: "陈明勇",
			}
			data, err := xml.Marshal(result)
			require.NoError(t, err)
			writer.Header().Set("Content-Type", "application/xml")
			_, _ = writer.Write(data)
		})
		http.HandleFunc("GET /invalid-ct", func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "xxx")
			_, _ = writer.Write(nil)
		})
		err := http.ListenAndServe(":8080", nil)
		require.NoError(t, err)
	})
}

func TestRequest_DoAndParse(t *testing.T) {
	// 启动一个 web 服务器
	go runWebServer(t)

	t.Run("test PUT", func(t *testing.T) {
		client := NewDefault()
		// 创建一个 map
		data := map[string]string{
			"name": "陈明勇",
		}
		// 序列化 map 为 JSON
		jsonData, err := json.Marshal(data)
		require.NoError(t, err)
		// 创建一个 bytes.Buffer 并写入 JSON 数据
		buffer := new(bytes.Buffer)
		_, err = buffer.Write(jsonData)
		require.NoError(t, err)

		putReq := client.Put("http://localhost:8080/test").SetBody(buffer)
		var result map[string]any
		err = putReq.DoAndParse(context.Background(), &result)
		require.NoError(t, err)
		require.Equal(t, result, map[string]any{"name": "陈明勇"})
	})

	t.Run("test encodeBody error", func(t *testing.T) {
		client := NewDefault()

		putReq := client.Put("http://localhost:8080/test").SetBody(map[string]any{
			"name": "陈明勇",
		}).SetBodyEncodeFunc(func(body any) (io.Reader, error) {
			return nil, errors.New("encode error")
		})
		var result map[string]any
		err := putReq.DoAndParse(context.Background(), &result)
		require.Equal(t, err, errors.New("encode error"))
	})

	t.Run("test decodeBody error", func(t *testing.T) {
		client := NewDefault()

		deleteReq := client.Delete("http://localhost:8080/test")
		err := deleteReq.DoAndParse(context.Background(), nil)
		var invalidUnmarshalError = &json.InvalidUnmarshalError{}
		require.True(t, errors.As(err, &invalidUnmarshalError))
	})

	t.Run("test text/plain", func(t *testing.T) {
		client := NewDefault()
		getReq := client.AddHeader("X-Name", "陈明勇").
			AddQuery("name", "陈明勇").Get("http://localhost:8080/test")
		var result string
		err := getReq.DoAndParse(context.Background(), &result)
		require.NoError(t, err)
		require.Equal(t, result, "hello world")
	})

	t.Run("test text/plain error", func(t *testing.T) {
		client := NewDefault()
		getReq := client.AddHeader("X-Name", "陈明勇").
			AddQuery("name", "陈明勇").Get("http://localhost:8080/test")
		var result string
		err := getReq.DoAndParse(context.Background(), result)
		require.Equal(t, err, fmt.Errorf("expected dst to be *string, but got string"))
	})

	t.Run("test application/xml", func(t *testing.T) {
		client := NewDefault()
		getReq := client.Get("http://localhost:8080/xml")
		var result = struct {
			Name string `xml:"name"`
		}{}
		err := getReq.DoAndParse(context.Background(), &result)
		require.NoError(t, err)
		require.Equal(t, result, struct {
			Name string `xml:"name"`
		}{Name: "陈明勇"})
	})
	t.Run("test invalid content-type", func(t *testing.T) {
		client := NewDefault()
		getReq := client.Get("http://localhost:8080/invalid-ct?name=陈明勇")
		err := getReq.SetQuery("age", "18").DoAndParse(context.Background(), nil)
		require.Equal(t, err, &UnsupportedContentTypeError{ContentType: "xxx"})
	})
}
