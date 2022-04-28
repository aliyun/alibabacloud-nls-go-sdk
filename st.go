/*
st.go

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
	ST_NAMESPACE = "SpeechTranscriber"

	//name field
	ST_START_NAME = "StartTranscription"
	ST_STOP_NAME  = "StopTranscription"
	ST_CTRL_NAME  = "ControlTranscriber"

	ST_STARTED_NAME        = "TranscriptionStarted"
	ST_SENTENCE_BEGIN_NAME = "SentenceBegin"
	ST_SENTENCE_END_NAME   = "SentenceEnd"
	ST_RESULT_CHG_NAME     = "TranscriptionResultChanged"
	ST_COMPLETED_NAME      = "TranscriptionCompleted"
)

type SpeechTranscriptionStartParam struct {
	Format                         string `json:"format,omitempty"`
	SampleRate                     int    `json:"sample_rate,omitempty"`
	EnableIntermediateResult       bool   `json:"enable_intermediate_result"`
	EnablePunctuationPredition     bool   `json:"enable_punctuation_prediction"`
	EnableInverseTextNormalization bool   `json:"enable_inverse_text_normalization"`
	MaxSentenceSilence             int    `json:"max_sentence_silence,omitempty"`
	EnableWords                    bool   `json:"enable_words"`
}

func DefaultSpeechTranscriptionParam() SpeechTranscriptionStartParam {
	return SpeechTranscriptionStartParam{
		Format:                         "pcm",
		SampleRate:                     16000,
		EnableIntermediateResult:       true,
		EnablePunctuationPredition:     true,
		EnableInverseTextNormalization: true,
		MaxSentenceSilence:             800,
		EnableWords:                    false,
	}
}

type SpeechTranscription struct {
	nls    *nlsProto
	taskId string

	startCh chan bool
	stopCh  chan bool

	lk sync.Mutex

	onTaskFailed    func(text string, param interface{})
	onStarted       func(text string, param interface{})
	onSentenceBegin func(text string, param interface{})
	onSentenceEnd   func(text string, param interface{})
	onResultChanged func(text string, param interface{})
	onCompleted     func(text string, param interface{})
	onClose         func(param interface{})

	StartParam map[string]interface{}
	UserParam  interface{}
}

func checkStNlsProto(proto *nlsProto) *SpeechTranscription {
	if proto == nil {
		log.Default().Fatal("empty proto check failed")
		return nil
	}

	st, ok := proto.param.(*SpeechTranscription)
	if !ok {
		log.Default().Fatal("proto param not SpeechTranscription instance")
		return nil
	}

	return st
}

func onStTaskFailedHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)
	if st.onTaskFailed != nil {
		st.onTaskFailed(string(text), st.UserParam)
	}

	st.lk.Lock()
  defer st.lk.Unlock()

  if st.startCh != nil {
		st.startCh <- false
		close(st.startCh)
		st.startCh = nil
	}

	if st.stopCh != nil {
		st.stopCh <- false
		close(st.stopCh)
		st.stopCh = nil
	}
}

func onStConnectedHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)

	req := CommonRequest{}
	req.Context = DefaultContext
	req.Header.Appkey = st.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = ST_START_NAME
	req.Header.Namespace = ST_NAMESPACE
	req.Header.TaskId = st.taskId
	req.Payload = st.StartParam

	b, _ := json.Marshal(req)
	st.nls.logger.Println("send:", string(b))
	st.nls.cmd(string(b))
}

func onStCloseHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)
	if st.onClose != nil {
		st.onClose(st.UserParam)
	}

	st.nls.shutdown()
}

func onStStartedHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)
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

func onStSentenceBeginHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)
	if st.onSentenceBegin != nil {
		st.onSentenceBegin(string(text), st.UserParam)
	}
}

func onStSentenceEndHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)
	if st.onSentenceEnd != nil {
		st.onSentenceEnd(string(text), st.UserParam)
	}
}

func onStResultChangedHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)
	if st.onResultChanged != nil {
		st.onResultChanged(string(text), st.UserParam)
	}
}

func onStCompletedHandler(isErr bool, text []byte, proto *nlsProto) {
	st := checkStNlsProto(proto)
	if st.onCompleted != nil {
		st.onCompleted(string(text), st.UserParam)
	}

	st.lk.Lock()
  defer st.lk.Unlock()
  if st.stopCh != nil {
		st.stopCh <- true
		st.stopCh = nil
	}
}

var stProto = commonProto{
	namespace: ST_NAMESPACE,
	handlers: map[string]func(bool, []byte, *nlsProto){
		CLOSE_HANDLER:          onStCloseHandler,
		CONNECTED_HANDLER:      onStConnectedHandler,
		ST_STARTED_NAME:        onStStartedHandler,
		ST_SENTENCE_BEGIN_NAME: onStSentenceBeginHandler,
		ST_SENTENCE_END_NAME:   onStSentenceEndHandler,
		ST_RESULT_CHG_NAME:     onStResultChangedHandler,
		ST_COMPLETED_NAME:      onStCompletedHandler,
		TASK_FAILED_NAME:       onStTaskFailedHandler,
	},
}

func newSpeechTranscriptionProto() *commonProto {
	return &stProto
}

func NewSpeechTranscription(config *ConnectionConfig,
	logger *NlsLogger,
	taskfailed func(string, interface{}),
	started func(string, interface{}),
	sentencebegin func(string, interface{}),
	sentenceend func(string, interface{}),
	resultchanged func(string, interface{}),
	completed func(string, interface{}),
	closed func(interface{}),
	param interface{}) (*SpeechTranscription, error) {
	st := new(SpeechTranscription)
	proto := newSpeechTranscriptionProto()
	if logger == nil {
		logger = DefaultNlsLog()
	}

	nls, err := newNlsProto(config, proto, logger, st)
	if err != nil {
		return nil, err
	}

	st.nls = nls
	st.UserParam = param
	st.onTaskFailed = taskfailed
	st.onStarted = started
	st.onSentenceBegin = sentencebegin
	st.onSentenceEnd = sentenceend
	st.onResultChanged = resultchanged
	st.onCompleted = completed
	st.onClose = closed
	return st, nil
}

func (st *SpeechTranscription) Start(param SpeechTranscriptionStartParam, extra map[string]interface{}) (chan bool, error) {
	if st.nls == nil {
		return nil, errors.New("empty nls: using NewSpeechTranscription to create a valid instance")
	}

	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(b, &st.StartParam)
	if extra != nil {
		if st.StartParam == nil {
			st.StartParam = extra
		} else {
			for k, v := range extra {
				st.StartParam[k] = v
			}
		}
	}
	st.taskId = getUuid()
	err = st.nls.Connect()
	if err != nil {
		return nil, err
	}

  st.lk.Lock()
  defer st.lk.Unlock()

	st.startCh = make(chan bool, 1)
	return st.startCh, nil
}

func (st *SpeechTranscription) Ctrl(param map[string]interface{}) error {
	if st.nls == nil {
		return errors.New("empty nls: using NewSpeechTranscription to create a valid instance")
	}

	req := CommonRequest{}
	req.Context = DefaultContext
	req.Header.Appkey = st.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = ST_CTRL_NAME
	req.Header.Namespace = ST_NAMESPACE
	req.Header.TaskId = st.taskId
	req.Payload = param

	b, _ := json.Marshal(req)
	err := st.nls.cmd(string(b))
	if err != nil {
		return err
	}

	return nil
}

func (st *SpeechTranscription) Stop() (chan bool, error) {
	if st.nls == nil {
		return nil, errors.New("empty nls: using NewSpeechTranscription to create a valid instance")
	}


	req := CommonRequest{}
	req.Context = DefaultContext
	req.Header.Appkey = st.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = ST_STOP_NAME
	req.Header.Namespace = ST_NAMESPACE
	req.Header.TaskId = st.taskId

	b, _ := json.Marshal(req)
	err := st.nls.cmd(string(b))
	if err != nil {
		return nil, err
	}

  st.lk.Lock()
  defer st.lk.Unlock()
	st.stopCh = make(chan bool, 1)
	return st.stopCh, nil
}

func (st *SpeechTranscription) Shutdown() {
	if st.nls == nil {
		return
	}

  st.nls.shutdown()
  st.lk.Lock()
  defer st.lk.Unlock()
	if st.startCh != nil {
		st.startCh <- false
		close(st.startCh)
		st.startCh = nil
	}

	if st.stopCh != nil {
		st.stopCh <- false
		close(st.stopCh)
		st.stopCh = nil
	}
}

func (st *SpeechTranscription) SendAudioData(data []byte) error {
	if st.nls == nil {
		return errors.New("empty nls: using NewSpeechTranscription to create a valid instance")
	}

	return st.nls.sendRawData(data)
}
