/*
utils.go

Copyright 1999-present Alibaba Group Holding Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/


package nls

import (
	"github.com/satori/go.uuid"
	"io"
	"strings"
)

const (
	SDK_VERSION  = "0.0.1fix"
	SDK_NAME     = "nls-go-sdk"
	SDK_LANGUAGE = "go"

	//AFORMAT
	PCM  = "pcm"
	WAV  = "wav"
	OPUS = "opus"
	OPU  = "opu"

	//token
	DEFAULT_DISTRIBUTE = "cn-shanghai"
	DEFAULT_DOMAIN     = "nls-meta.cn-shanghai.aliyuncs.com"
	DEFAULT_VERSION    = "2019-02-28"

	DEFAULT_SEC_WEBSOCKET_KEY = "x3JJHMbDL1EzLkh9GBhXDw=="
	DEFAULT_SEC_WEBSOCKET_VER = "13"

	DEFAULT_X_NLS_TOKEN_KEY = "X-NLS-Token"

	DEFAULT_URL = "wss://nls-gateway.cn-shanghai.aliyuncs.com/ws/v1"

	TASK_FAILED_NAME = "TaskFailed"
	
  AUDIO_FORMAT_KEY        = "format"
	SAMPLE_RATE_KEY         = "sample_rate"
	ENABLE_INTERMEDIATE_KEY = "enable_intermediate_result"
	ENABLE_PP_KEY           = "enable_punctuation_predition"
	ENABLE_ITN_KEY          = "enable_inverse_text_normalization"
)

type Chunk struct {
	Data []byte
}

type ChunkBuffer struct {
	Data []*Chunk
}

type TokenResult struct {
	UserId     string `json:"UserId"`
	Id         string `json:"Id"`
	ExpireTime int64  `json:"ExpireTime"`
}

type TokenResultMessage struct {
	ErrMsg      string      `json:"ErrMsg"`
	TokenResult TokenResult `json:"Token"`
}

type Header struct {
	MessageId string `json:"message_id"`
	TaskId    string `json:"task_id"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Appkey    string `json:"appkey"`
}

type SDK struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"`
}

type Context struct {
	Sdk       SDK                    `json:"sdk"`
	App       map[string]interface{} `json:"app,omitempty"`
	System    map[string]interface{} `json:"system,omitempty"`
	Device    map[string]interface{} `json:"device,omitempty"`
	Network   map[string]interface{} `json:"network,omitempty"`
	Geography map[string]interface{} `json:"geography,omitempty"`
	Bridge    map[string]interface{} `json:"bridge,omitempty"`
	Custom    map[string]interface{} `json:"custom,omitempty"`
}

var DefaultContext = Context{
	Sdk: SDK{
		Name:     SDK_NAME,
		Version:  SDK_VERSION,
		Language: SDK_LANGUAGE,
	},
}

type CommonResponse struct {
	Header  Header                 `json:"header"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type CommonRequest struct {
	Header  Header                 `json:"header"`
	Payload map[string]interface{} `json:"payload,omitempty"`
	Context Context                `json:"context"`
}

func LoadPcmInChunk(r io.Reader, chunkSize int) *ChunkBuffer {
	buffer := new(ChunkBuffer)
	buffer.Data = make([]*Chunk, 0)
	for {
		chunk := new(Chunk)
		chunk.Data = make([]byte, chunkSize)
		i, err := r.Read(chunk.Data)
		if err == io.EOF {
			break
		} else {
			if i != chunkSize {
				chunk2 := new(Chunk)
				chunk2.Data = make([]byte, i)
				copy(chunk2.Data, chunk.Data)
				chunk = chunk2
			}
			buffer.Data = append(buffer.Data, chunk)
		}
	}

	return buffer
}

func getUuid() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
