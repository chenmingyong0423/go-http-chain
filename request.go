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
	"context"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	client         *http.Client
	url            string
	method         string
	headers        http.Header
	queryValues    url.Values
	body           any
	bodyBytes      io.Reader
	bodyEncodeFunc BodyEncodeFunc
}

func (r *Request) AddQuery(key, value string) *Request {
	r.queryValues.Add(key, value)
	return r
}

func (r *Request) SetQuery(key, value string) *Request {
	r.queryValues.Set(key, value)
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.headers.Set(key, value)
	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	r.headers.Add(key, value)
	return r
}

func (r *Request) SetBodyEncodeFunc(fn BodyEncodeFunc) *Request {
	r.bodyEncodeFunc = fn
	return r
}

func (r *Request) SetBody(body any) *Request {
	r.body = body
	if bodyReader, ok := body.(io.Reader); ok {
		r.bodyBytes = bodyReader
	}
	return r
}

func (r *Request) encodeBody() (io.Reader, error) {
	if r.bodyBytes != nil {
		return r.bodyBytes, nil
	}
	if r.body != nil {
		body, err := r.bodyEncodeFunc(r.body)
		if err != nil {
			return nil, err
		}
		r.bodyBytes = body
		return body, nil
	}
	return nil, nil
}

func (r *Request) Call(ctx context.Context) *Response {
	response := new(Response)
	body, err := r.encodeBody()
	if err != nil {
		response.err = err
		return response
	}
	req, err := http.NewRequestWithContext(ctx, r.method, r.url, body)
	if err != nil {
		response.err = err
		return response
	}
	req.Header = r.headers
	req.URL.RawQuery = r.queryValues.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		response.err = err
		return response
	}
	response.resp = resp
	return response
}
