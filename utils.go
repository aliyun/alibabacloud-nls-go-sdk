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

package dash

import (
	"io"

	uuid "github.com/satori/go.uuid"
)

const (
	SDK_VERSION  = "0.0.1fix"
	SDK_NAME     = "nls-go-sdk"
	SDK_LANGUAGE = "go"

	//AFORMAT
	PCM = "pcm"
	WAV = "wav"
	MP3 = "mp3"

	DEFAULT_URL             = "wss://dashscope.aliyuncs.com/api-ws/v1/inference"
	APIKEY_ENV_KEY_NAME     = "DASHSCOPE_API_KEY"
	DEFAULT_WS_RBUFFER_SIZE = 4096
	DEFAULT_WS_WBUFFER_SIZE = 4096

	AUTHORIZATION_KEY = "Authorization"

	TASK_GROUP_AUDIO_NAME             = "audio"
	TASK_TTS_NAME                     = "tts"
	TASK_TTS_SYNTHESIZE_FUNCTION_NAME = "SpeechSynthesizer"

	TASK_RUN_NAME           = "run-task"
	STREAMING_NAME          = "streaming"
	STREAMING_MODE_OUT_NAME = "out"

	TEXT_TPYE_NAME       = "text_type"
	PLAIN_TEXT_TYPE_NAME = "PlainText"

	//for audio tts
	AUDIO_FORMAT_KEY            = "format"
	VOLUME_KEY                  = "volume"
	PITCH_KEY                   = "pitch"
	SAMPLE_RATE_KEY             = "sample_rate"
	ENABLE_PHONME_TIMESTAMP_KEY = "phoneme_timestamp_enabled"
	ENABLE_WORD_TIMESTAMP_KEY   = "word_timestamp_enabled"

	//event
	TASK_STARTED_EVENT_NAME   = "task-started"
	TASK_RESULT_EVENT_NAME    = "result-generated"
	TASK_COMPLETED_EVENT_NAME = "task-finished"
	TASK_FAILED_EVENT_NAME    = "task-failed"
)

type Chunk struct {
	Data []byte
}

type ChunkBuffer struct {
	Data []*Chunk
}

type RequestHeader struct {
	Action    string `json:"action"`
	TaskId    string `json:"task_id"`
	Streaming string `json:"streaming"`
}

type ResponseHeader struct {
	TaskId string `json:"task_id"`
	Event  string `json:"event"`
}

type RequestPayload struct {
	Model      string                 `json:"model"`
	TaskGroup  string                 `json:"task_group"`
	Task       string                 `json:"task"`
	Function   string                 `json:"function"`
	Input      map[string]interface{} `json:"input"`
	Parameters map[string]interface{} `json:"parameters"`
}

type ResponsePayload struct {
	Oupput map[string]interface{} `json:"output,omitempty"`
	Usage  map[string]interface{} `json:"usage,omitempty"`
}

type CommonResponse struct {
	Header  ResponseHeader  `json:"header"`
	Payload ResponsePayload `json:"payload,omitempty"`
}

type CommonRequest struct {
	Header  RequestHeader  `json:"header"`
	Payload RequestPayload `json:"payload"`
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
	return uuid.NewV4().String()
}
