/*
tts.go

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
	"encoding/json"
	"errors"
	"log"
	"sync"
)

type SpeechSynthesisParam struct {
	TextType               string  `json:"text_type"`
	Format                 string  `json:"format,omitempty"`
	SampleRate             int     `json:"sample_rate,omitempty"`
	Volume                 int     `json:"volume"`
	Rate                   float32 `json:"rate"`
	Pitch                  float32 `json:"pitch"`
	EnablePhonemeTimestamp bool    `json:"phoneme_timestamp_enabled"`
	EnableWordTimestamp    bool    `json:"word_timestamp_enabled"`
}

func DefaultSpeechSynthesisParam() SpeechSynthesisParam {
	return SpeechSynthesisParam{
		TextType:               PLAIN_TEXT_TYPE_NAME,
		Format:                 "wav",
		SampleRate:             16000,
		Volume:                 50,
		Rate:                   1.0,
		Pitch:                  1.0,
		EnablePhonemeTimestamp: false,
		EnableWordTimestamp:    false,
	}
}

type SpeechSynthesisInput struct {
	Text string `json:"text"`
}

type SpeechSynthesis struct {
	dash   *dashProto
	taskId string

	completeChan chan bool
	lk           sync.Mutex

	onTaskFailed      func(text string, param interface{})
	onSynthesisResult func(data []byte, param interface{})
	onCompleted       func(text string, param interface{})
	onMetaInfo        func(text string, param interface{})
	onClose           func(param interface{})
	onStarted         func(taskid string, param interface{})

	StartParam map[string]interface{}
	InputParam map[string]interface{}
	UserParam  interface{}
}

func checkTtsNlsProto(proto *dashProto) *SpeechSynthesis {
	if proto == nil {
		log.Default().Fatal("empty proto check failed")
		return nil
	}

	tts, ok := proto.param.(*SpeechSynthesis)
	if !ok {
		log.Default().Fatal("proto param not SpeechSynthesis instance")
		return nil
	}

	return tts
}

func onTtsTaskFailedHandler(isErr bool, text []byte, proto *dashProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onTaskFailed != nil {
		tts.onTaskFailed(string(text), tts.UserParam)
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	if tts.completeChan != nil {
		tts.completeChan <- false
		close(tts.completeChan)
		tts.completeChan = nil
	}
}

func onTtsConnectedHandler(isErr bool, text []byte, proto *dashProto) {
	tts := checkTtsNlsProto(proto)

	req := CommonRequest{}
	req.Header.TaskId = getUuid()
	req.Header.Action = TASK_RUN_NAME
	req.Header.Streaming = STREAMING_MODE_OUT_NAME
	req.Payload.Model = tts.dash.proto.Model
	req.Payload.TaskGroup = tts.dash.proto.TaskGroup
	req.Payload.Task = tts.dash.proto.Task
	req.Payload.Function = tts.dash.proto.Function
	req.Payload.Input = tts.InputParam
	req.Payload.Parameters = tts.StartParam

	tts.taskId = req.Header.TaskId

	b, _ := json.Marshal(req)
	tts.dash.logger.Println("send:", string(b))
	tts.dash.cmd(string(b))
}

func onTtsTaskStartedHandler(isErr bool, text []byte, proto *dashProto) {
	tts := checkTtsNlsProto(proto)

	if tts.onStarted != nil {
		tts.onStarted(tts.taskId, tts.UserParam)
	}
}

func onTtsCloseHandler(isErr bool, text []byte, proto *dashProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onClose != nil {
		tts.onClose(tts.UserParam)
	}

	tts.dash.shutdown()
}

func onTtsMetaInfoHandler(isErr bool, text []byte, proto *dashProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onMetaInfo != nil {
		tts.onMetaInfo(string(text), tts.UserParam)
	}
}

func onTtsRawResultHandler(isErr bool, text []byte, proto *dashProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onSynthesisResult != nil {
		tts.onSynthesisResult(text, tts.UserParam)
	}
}

func onTtsCompletedHandler(isErr bool, text []byte, proto *dashProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onCompleted != nil {
		tts.onCompleted(string(text), tts.UserParam)
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	if tts.completeChan != nil {
		tts.completeChan <- true
		close(tts.completeChan)
		tts.completeChan = nil
	}
}

var ttsProto = commonProto{
	Model:     "",
	TaskGroup: TASK_GROUP_AUDIO_NAME,
	Task:      TASK_TTS_NAME,
	Function:  TASK_TTS_SYNTHESIZE_FUNCTION_NAME,
	handlers: map[string]func(bool, []byte, *dashProto){
		CLOSE_HANDLER:             onTtsCloseHandler,
		CONNECTED_HANDLER:         onTtsConnectedHandler,
		RAW_HANDLER:               onTtsRawResultHandler,
		TASK_STARTED_EVENT_NAME:   onTtsTaskStartedHandler,
		TASK_COMPLETED_EVENT_NAME: onTtsCompletedHandler,
		TASK_FAILED_EVENT_NAME:    onTtsTaskFailedHandler,
		TASK_RESULT_EVENT_NAME:    onTtsMetaInfoHandler,
	},
}

func NewSpeechSynthesis(config *ConnectionConfig,
	logger *NlsLogger,
	started func(string, interface{}),
	taskfailed func(string, interface{}),
	synthesisresult func([]byte, interface{}),
	metainfo func(string, interface{}),
	completed func(string, interface{}),
	closed func(interface{}),
	param interface{}) (*SpeechSynthesis, error) {
	tts := new(SpeechSynthesis)
	proto := &ttsProto
	if logger == nil {
		logger = DefaultNlsLog()
	}

	dash, err := newDashProto(config, proto, logger, tts)
	if err != nil {
		return nil, err
	}

	tts.dash = dash
	tts.UserParam = param
	tts.onTaskFailed = taskfailed
	tts.onSynthesisResult = synthesisresult
	tts.onMetaInfo = metainfo
	tts.onCompleted = completed
	tts.onClose = closed

	return tts, nil
}

func (tts *SpeechSynthesis) Start(model string,
	text string,
	param SpeechSynthesisParam,
	extra map[string]interface{}) (chan bool, error) {
	if tts.dash == nil {
		return nil, errors.New("empty dash obj: using NewSpeechSynthesis to create a valid instance")
	}

	if model == "" {
		return nil, errors.New("empty model")
	}
	tts.dash.proto.Model = model

	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(b, &tts.StartParam)
	if extra != nil {
		if tts.StartParam == nil {
			tts.StartParam = extra
		} else {
			for k, v := range extra {
				tts.StartParam[k] = v
			}
		}
	}

	input := SpeechSynthesisInput{
		Text: text,
	}

	b, err = json.Marshal(input)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(b, &tts.InputParam)
	err = tts.dash.Connect()
	if err != nil {
		return nil, err
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	tts.completeChan = make(chan bool, 1)
	return tts.completeChan, nil
}

func (tts *SpeechSynthesis) Shutdown() {
	if tts.dash == nil {
		return
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	tts.dash.shutdown()
	if tts.completeChan != nil {
		tts.completeChan <- false
		close(tts.completeChan)
		tts.completeChan = nil
	}
}
