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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	resp *http.Response
	err  error
}

func (r *Response) Result() (*http.Response, error) {
	return r.resp, r.err
}

func (r *Response) DecodeRespBody(dst any) error {
	if r.err != nil {
		return r.err
	}
	return DecodeRespBody(r.resp.Header.Get(HeaderContentType), r.resp.Body, dst)
}

func DecodeRespBody(contentType string, body io.ReadCloser, dst any) error {
	switch contentType {
	case ContentTypeApplicationJSON, ContentTypeApplicationJSONCharacterUTF8:
		return json.NewDecoder(body).Decode(dst)
	case ContentTypeApplicationXML, ContentTypeApplicationXMLCharacterUTF8, ContentTypeTextXML, ContentTypeTextXMLCharacterUTF8:
		return xml.NewDecoder(body).Decode(dst)
	case ContentTypeTextPlain, ContentTypeTextPlainCharacterUTF8:
		bytes, err := io.ReadAll(body)
		if err != nil {
			return err
		}

		strPtr, ok := dst.(*string)
		if !ok {
			return fmt.Errorf("expected dst to be *string, but got %T", dst)
		}
		*strPtr = string(bytes)
		return nil
	default:
		return &UnsupportedContentTypeError{ContentType: contentType}
	}
}
