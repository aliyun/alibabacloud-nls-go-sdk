/*
streamtts.go

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
	STREAM_TTS_NAMESPACE = "FlowingSpeechSynthesizer"
	//name field
	STREAM_TTS_START_NAME              = "StartSynthesis"
	STREAM_TTS_RUN_NAME                = "RunSynthesis"
	STREAM_TTS_STOP_NAME               = "StopSynthesis"
	STREAM_TTS_SYNTHESIS_STARTED_NAME  = "SynthesisStarted"
	STREAM_TTS_SENTENCE_BEGIN_NAME     = "SentenceBegin"
	STREAM_TTS_SENTENCE_SYNTHESIS_NAME = "SentenceSynthesis"
	STREAM_TTS_SENTENCE_END_NAME       = "SentenceEnd"
	STREAM_TTS_COMPLETED_NAME          = "SynthesisCompleted"
)

type StreamSpeechSynthesisStartParam struct {
	SessionId              string `json:"session_id"`
	Voice                  string `json:"voice"`
	Format                 string `json:"format,omitempty"`
	SampleRate             int    `json:"sample_rate,omitempty"`
	Volume                 int    `json:"volume,omitempty"`
	SpeechRate             int    `json:"speech_rate,omitempty"`
	PitchRate              int    `json:"pitch_rate,omitempty"`
	EnableSubtitle         bool   `json:"enable_subtitle,omitempty"`
	EnablePhonemeTimestamp bool   `json:"enable_phoneme_timestamp,omitempty"`
}

func DefaultStreamSpeechSynthesisParam() StreamSpeechSynthesisStartParam {
	return StreamSpeechSynthesisStartParam{
		SessionId:  getUuid(),
		Voice:      "longxiaochun",
		Format:     "wav",
		SampleRate: 16000,
		Volume:     50,
		SpeechRate: 0,
		PitchRate:  0,
	}
}

type StreamSpeechSynthesis struct {
	nls    *nlsProto
	taskId string

	startCh chan bool
	stopCh  chan bool

	//completeChan chan bool
	lk sync.Mutex

	onTaskFailed        func(text string, param interface{})
	onStarted           func(text string, param interface{})
	onSentenceBegin     func(text string, param interface{})
	onSentenceEnd       func(text string, param interface{})
	onSentenceSynthesis func(text string, param interface{})
	onSynthesisResult   func(data []byte, param interface{})
	onCompleted         func(text string, param interface{})
	onClose             func(param interface{})

	StartParam map[string]interface{}
	UserParam  interface{}
}

func checkStreamTtsNlsProto(proto *nlsProto) *StreamSpeechSynthesis {
	if proto == nil {
		log.Default().Fatal("empty proto check failed")
		return nil
	}

	tts, ok := proto.param.(*StreamSpeechSynthesis)
	if !ok {
		log.Default().Fatal("proto param not SpeechSynthesis instance")
		return nil
	}

	return tts
}

func onStreamTtsTaskFailedHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkStreamTtsNlsProto(proto)
	if tts.onTaskFailed != nil {
		tts.onTaskFailed(string(text), tts.UserParam)
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	if tts.startCh != nil {
		tts.startCh <- false
		close(tts.startCh)
		tts.startCh = nil
	}

	if tts.stopCh != nil {
		tts.stopCh <- false
		close(tts.stopCh)
		tts.stopCh = nil
	}
}

func onStreamTtsConnectedHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkStreamTtsNlsProto(proto)

	req := CommonRequest{}
	req.Context = DefaultContext
	req.Header.Appkey = tts.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = STREAM_TTS_START_NAME
	req.Header.Namespace = STREAM_TTS_NAMESPACE
	req.Header.TaskId = tts.taskId
	req.Payload = tts.StartParam

	b, _ := json.Marshal(req)
	tts.nls.logger.Println("send:", string(b))
	tts.nls.cmd(string(b))
}

func onStreamTtsCloseHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkStreamTtsNlsProto(proto)
	if tts.onClose != nil {
		tts.onClose(tts.UserParam)
	}

	tts.nls.shutdown()
}

func onStreamTtsRawResultHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkStreamTtsNlsProto(proto)
	if tts.onSynthesisResult != nil {
		tts.onSynthesisResult(text, tts.UserParam)
	}
}

func onStreamTtsStartedHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStreamTtsNlsProto(proto)
	if st.onStarted != nil {
		st.onStarted(string(text), st.UserParam)
	}
	st.lk.Lock()
	defer st.lk.Unlock()
	if st.startCh != nil {
		st.startCh <- true
		close(st.startCh)
		st.startCh = nil
	}
}

func onStreamTtsSentenceBeginHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStreamTtsNlsProto(proto)
	if st.onSentenceBegin != nil {
		st.onSentenceBegin(string(text), st.UserParam)
	}
}

func onStreamTtsSentenceSynthesisHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStreamTtsNlsProto(proto)
	if st.onSentenceSynthesis != nil {
		st.onSentenceSynthesis(string(text), st.UserParam)
	}
}

func onStreamTtsSentenceEndHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStreamTtsNlsProto(proto)
	if st.onSentenceEnd != nil {
		st.onSentenceEnd(string(text), st.UserParam)
	}
}

func onStreamTtsCompletedHandler(isErr bool, text []byte, proto *nlsProto) {
	tts := checkStreamTtsNlsProto(proto)
	if tts.onCompleted != nil {
		tts.onCompleted(string(text), tts.UserParam)
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	if tts.stopCh != nil {
		tts.stopCh <- true
		tts.stopCh = nil
	}
}

var streamTtsProto = commonProto{
	namespace: STREAM_TTS_NAMESPACE,
	handlers: map[string]func(bool, []byte, *nlsProto){
		CLOSE_HANDLER:                      onStreamTtsCloseHandler,
		CONNECTED_HANDLER:                  onStreamTtsConnectedHandler,
		RAW_HANDLER:                        onStreamTtsRawResultHandler,
		STREAM_TTS_SYNTHESIS_STARTED_NAME:  onStreamTtsStartedHandler,
		STREAM_TTS_SENTENCE_BEGIN_NAME:     onStreamTtsSentenceBeginHandler,
		STREAM_TTS_SENTENCE_SYNTHESIS_NAME: onStreamTtsSentenceSynthesisHandler,
		STREAM_TTS_SENTENCE_END_NAME:       onStreamTtsSentenceEndHandler,
		STREAM_TTS_COMPLETED_NAME:          onStreamTtsCompletedHandler,
		TASK_FAILED_NAME:                   onStreamTtsTaskFailedHandler,
	},
}

func newStreamSpeechSynthesisProto() *commonProto {
	return &streamTtsProto
}

func NewStreamSpeechSynthesis(config *ConnectionConfig,
	logger *NlsLogger,
	onTaskFailed func(string, interface{}),
	onStarted func(text string, param interface{}),
	onSentenceBegin func(text string, param interface{}),
	onSentenceEnd func(text string, param interface{}),
	onSentenceSynthesis func(text string, param interface{}),
	onSynthesisResult func([]byte, interface{}),
	completed func(string, interface{}),
	closed func(interface{}),
	param interface{}) (*StreamSpeechSynthesis, error) {
	tts := new(StreamSpeechSynthesis)
	proto := newStreamSpeechSynthesisProto()
	if logger == nil {
		logger = DefaultNlsLog()
	}

	nls, err := newNlsProto(config, proto, logger, tts)
	if err != nil {
		return nil, err
	}

	tts.nls = nls
	tts.UserParam = param
	tts.onTaskFailed = onTaskFailed
	tts.onStarted = onStarted
	tts.onSentenceBegin = onSentenceBegin
	tts.onSentenceEnd = onSentenceEnd
	tts.onSentenceSynthesis = onSentenceSynthesis
	tts.onSynthesisResult = onSynthesisResult
	tts.onCompleted = completed
	tts.onClose = closed
	return tts, nil
}

func (tts *StreamSpeechSynthesis) Start(param StreamSpeechSynthesisStartParam, extra map[string]interface{}) (chan bool, error) {
	if tts.nls == nil {
		return nil, errors.New("empty nls: using NewStreamSpeechSynthesis to create a valid instance")
	}

	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(b, &tts.StartParam)
	if extra != nil {
		if tts.StartParam == nil {
			tts.StartParam = extra
		} else {
			for k, v := range extra {
				tts.StartParam[k] = v
			}
		}
	}
	tts.taskId = getUuid()
	err = tts.nls.Connect()
	if err != nil {
		return nil, err
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()

	tts.startCh = make(chan bool, 1)
	return tts.startCh, nil
}

func (tts *StreamSpeechSynthesis) SendText(text string) error {
	if tts.nls == nil {
		return errors.New("empty nls: using NewStreamSpeechSynthesis to create a valid instance")
	}

	req := CommonRequest{}
	req.Header.Appkey = tts.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = STREAM_TTS_RUN_NAME
	req.Header.Namespace = STREAM_TTS_NAMESPACE
	req.Header.TaskId = tts.taskId
	req.Payload = map[string]interface{}{
		"text": text,
	}
	b, _ := json.Marshal(req)
	return tts.nls.cmd(string(b))
}

func (tts *StreamSpeechSynthesis) Stop() (chan bool, error) {
	if tts.nls == nil {
		return nil, errors.New("empty nls: using NewStreamSpeechSynthesis to create a valid instance")
	}

	req := CommonRequest{}
	req.Header.Appkey = tts.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = STREAM_TTS_STOP_NAME
	req.Header.Namespace = STREAM_TTS_NAMESPACE
	req.Header.TaskId = tts.taskId

	b, _ := json.Marshal(req)
	err := tts.nls.cmd(string(b))
	if err != nil {
		return nil, err
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	tts.stopCh = make(chan bool, 1)
	return tts.stopCh, nil
}

func (tts *StreamSpeechSynthesis) Shutdown() {
	if tts.nls == nil {
		return
	}

	tts.lk.Lock()
	defer tts.lk.Unlock()
	tts.nls.shutdown()
	if tts.startCh != nil {
		tts.startCh <- false
		close(tts.startCh)
		tts.startCh = nil
	}

	if tts.stopCh != nil {
		tts.stopCh <- false
		close(tts.stopCh)
		tts.stopCh = nil
	}
}
