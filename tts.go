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

package nls

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
)

const (
	//namespace field
	TTS_NAMESPACE = "SpeechSynthesizer"

	//name field
	TTS_START_NAME     = "StartSynthesis"
	TTS_COMPLETED_NAME = "SynthesisCompleted"
	TTS_METAINFO_NAME  = "MetaInfo"
)

type SpeechSynthesisStartParam struct {
	Voice          string `json:"voice"`
	Format         string `json:"format,omitempty"`
	SampleRate     int    `json:"sample_rate,omitempty"`
	Volume         int    `json:"volume"`
	SpeechRate     int    `json:"speech_rate"`
	PitchRate      int    `json:"pitch_rate"`
	EnableSubtitle bool   `json:"enable_subtitle"`
}

func DefaultSpeechSynthesisParam() SpeechSynthesisStartParam {
	return SpeechSynthesisStartParam{
		Voice:          "xiaoyun",
		Format:         "wav",
		SampleRate:     16000,
		Volume:         50,
		SpeechRate:     0,
		PitchRate:      0,
		EnableSubtitle: false,
	}
}

type SpeechSynthesis struct {
	nls    *nlsProto
	taskId string

	completeChan chan bool
	lk           sync.Mutex

	onTaskFailed      func(text string, param interface{})
	onSynthesisResult func(data []byte, param interface{})
	onCompleted       func(text string, param interface{})
	onMetaInfo        func(text string, param interface{})
	onClose           func(param interface{})

	StartParam map[string]interface{}
	UserParam  interface{}
}

func checkTtsNlsProto(proto *nlsProto) *SpeechSynthesis {
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

func onTtsTaskFailedHandler(isErr bool, text []byte, proto *nlsProto) {
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

func onTtsConnectedHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkTtsNlsProto(proto)

	req := CommonRequest{}
	req.Context = DefaultContext
	req.Header.Appkey = tts.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = TTS_START_NAME
	req.Header.Namespace = TTS_NAMESPACE
	req.Header.TaskId = tts.taskId
	req.Payload = tts.StartParam

	b, _ := json.Marshal(req)
	tts.nls.logger.Println("send:", string(b))
	tts.nls.cmd(string(b))
}

func onTtsCloseHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onClose != nil {
		tts.onClose(tts.UserParam)
	}

	tts.nls.shutdown()
}

func onTtsMetaInfoHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onMetaInfo != nil {
		tts.onMetaInfo(string(text), tts.UserParam)
	}
}

func onTtsRawResultHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkTtsNlsProto(proto)
	if tts.onSynthesisResult != nil {
		tts.onSynthesisResult(text, tts.UserParam)
	}
}

func onTtsCompletedHandler(isErr bool, text []byte, proto *nlsProto) {
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
	namespace: TTS_NAMESPACE,
	handlers: map[string]func(bool, []byte, *nlsProto){
		CLOSE_HANDLER:      onTtsCloseHandler,
		CONNECTED_HANDLER:  onTtsConnectedHandler,
		RAW_HANDLER:        onTtsRawResultHandler,
		TTS_COMPLETED_NAME: onTtsCompletedHandler,
		TASK_FAILED_NAME:   onTtsTaskFailedHandler,
		TTS_METAINFO_NAME:  onTtsMetaInfoHandler,
	},
}

func newSpeechSynthesisProto() *commonProto {
	return &ttsProto
}

func NewSpeechSynthesis(config *ConnectionConfig,
	logger *NlsLogger,
	taskfailed func(string, interface{}),
	synthesisresult func([]byte, interface{}),
	metainfo func(string, interface{}),
	completed func(string, interface{}),
	closed func(interface{}),
	param interface{}) (*SpeechSynthesis, error) {
	tts := new(SpeechSynthesis)
	proto := newSpeechSynthesisProto()
	if logger == nil {
		logger = DefaultNlsLog()
	}

	nls, err := newNlsProto(config, proto, logger, tts)
	if err != nil {
		return nil, err
	}

	tts.nls = nls
	tts.UserParam = param
	tts.onTaskFailed = taskfailed
	tts.onSynthesisResult = synthesisresult
	tts.onMetaInfo = metainfo
	tts.onCompleted = completed
	tts.onClose = closed
	return tts, nil
}

func (tts *SpeechSynthesis) Start(text string,
	param SpeechSynthesisStartParam,
	extra map[string]interface{}) (chan bool, error) {
	if tts.nls == nil {
		return nil, errors.New("empty nls: using NewSpeechSynthesis to create a valid instance")
	}

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
	tts.StartParam["text"] = text
	tts.taskId = getUuid()
	err = tts.nls.Connect()
	if err != nil {
		return nil, err
	}

  tts.lk.Lock()
  defer tts.lk.Unlock()
	tts.completeChan = make(chan bool, 1)
	return tts.completeChan, nil
}

func (tts *SpeechSynthesis) Shutdown() {
	if tts.nls == nil {
		return
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	tts.nls.shutdown()
	if tts.completeChan != nil {
		tts.completeChan <- false
		close(tts.completeChan)
		tts.completeChan = nil
	}
}
