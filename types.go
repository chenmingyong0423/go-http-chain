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

import "io"

const (
	HeaderContentType = "Content-Type"

	ContentTypeTextPlain                    = "text/plain"
	ContentTypeTextPlainCharacterUTF8       = "text/plain; charset=utf-8"
	ContentTypeApplicationJSON              = "application/json"
	ContentTypeApplicationJSONCharacterUTF8 = "application/json; charset=utf-8"
	ContentTypeApplicationXML               = "application/xml"
	ContentTypeApplicationXMLCharacterUTF8  = "application/xml; charset=utf-8"
	ContentTypeTextXML                      = "text/xml"
	ContentTypeTextXMLCharacterUTF8         = "text/xml; charset=utf-8"
	ContentTypeApplicationOctetStream       = "application/octet-stream"
	ContentTypeMultipartFormData            = "multipart/form-data"
	ContentTypeFormURLEncoded               = "application/x-www-form-urlencoded"
)

type BodyEncodeFunc func(body any) (io.Reader, error)
