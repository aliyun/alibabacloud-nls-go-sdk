/*
sr.go

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
	SR_NAMESPACE = "SpeechRecognizer"

	//name field
	SR_START_NAME = "StartRecognition"
	SR_STOP_NAME  = "StopRecognition"

	SR_STARTED_NAME    = "RecognitionStarted"
	SR_RESULT_CHG_NAME = "RecognitionResultChanged"
	SR_COMPLETED_NAME  = "RecognitionCompleted"
)

type SpeechRecognitionStartParam struct {
	Format                         string `json:"format,omitempty"`
	SampleRate                     int    `json:"sample_rate,omitempty"`
	EnableIntermediateResult       bool   `json:"enable_intermediate_result"`
	EnablePunctuationPrediction    bool   `json:"enable_punctuation_prediction"`
	EnableInverseTextNormalization bool   `json:"enable_inverse_text_normalization"`
}

func DefaultSpeechRecognitionParam() SpeechRecognitionStartParam {
	return SpeechRecognitionStartParam{
		Format:                         "pcm",
		SampleRate:                     16000,
		EnableIntermediateResult:       true,
		EnablePunctuationPrediction:    true,
		EnableInverseTextNormalization: true,
	}
}

type SpeechRecognition struct {
	nls    *nlsProto
	taskId string

	startCh chan bool
	stopCh  chan bool
	lk      sync.Mutex

	onTaskFailed    func(text string, param interface{})
	onStarted       func(text string, param interface{})
	onResultChanged func(text string, param interface{})
	onCompleted     func(text string, param interface{})
	onClose         func(param interface{})

	StartParam map[string]interface{}
	UserParam  interface{}
}

func checkSrNlsProto(proto *nlsProto) *SpeechRecognition {
	if proto == nil {
		log.Default().Fatal("empty proto check failed")
		return nil
	}

	sr, ok := proto.param.(*SpeechRecognition)
	if !ok {
		log.Default().Fatal("proto param not SpeechRecognition instance")
		return nil
	}

	return sr
}

func onSrTaskFailedHandler(isErr bool, text []byte, proto *nlsProto) {
	sr := checkSrNlsProto(proto)
	if sr.onTaskFailed != nil {
		sr.onTaskFailed(string(text), sr.UserParam)
	}

	sr.lk.Lock()
	defer sr.lk.Unlock()
	if sr.startCh != nil {
		sr.startCh <- false
		close(sr.startCh)
		sr.startCh = nil
	}

	if sr.stopCh != nil {
		sr.stopCh <- false
		close(sr.stopCh)
		sr.stopCh = nil
	}
}

func onSrConnectedHandler(isErr bool, text []byte, proto *nlsProto) {
	sr := checkSrNlsProto(proto)

	req := CommonRequest{}
	req.Context = DefaultContext
	req.Header.Appkey = sr.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = SR_START_NAME
	req.Header.Namespace = SR_NAMESPACE
	req.Header.TaskId = sr.taskId
	req.Payload = sr.StartParam

	b, _ := json.Marshal(req)
	sr.nls.logger.Println("send:", string(b))
	sr.nls.cmd(string(b))
}

func onSrCloseHandler(isErr bool, text []byte, proto *nlsProto) {
	sr := checkSrNlsProto(proto)
	if sr.onClose != nil {
		sr.onClose(sr.UserParam)
	}

	sr.nls.shutdown()
}

func onSrStartedHandler(isErr bool, text []byte, proto *nlsProto) {
	sr := checkSrNlsProto(proto)
	if sr.onStarted != nil {
		sr.onStarted(string(text), sr.UserParam)
	}

	sr.lk.Lock()
	defer sr.lk.Unlock()
	if sr.startCh != nil {
		sr.startCh <- true
		close(sr.startCh)
		sr.startCh = nil
	}
}

func onSrResultChangedHandler(isErr bool, text []byte, proto *nlsProto) {
	sr := checkSrNlsProto(proto)
	if sr.onResultChanged != nil {
		sr.onResultChanged(string(text), sr.UserParam)
	}
}

func onSrCompletedHandler(isErr bool, text []byte, proto *nlsProto) {
	sr := checkSrNlsProto(proto)
	if sr.onCompleted != nil {
		sr.onCompleted(string(text), sr.UserParam)
	}

	sr.lk.Lock()
	defer sr.lk.Unlock()
	if sr.stopCh != nil {
		sr.stopCh <- true
		close(sr.stopCh)
		sr.stopCh = nil
	}
}

var srProto = commonProto{
	namespace: SR_NAMESPACE,
	handlers: map[string]func(bool, []byte, *nlsProto){
		CLOSE_HANDLER:      onSrCloseHandler,
		CONNECTED_HANDLER:  onSrConnectedHandler,
		SR_STARTED_NAME:    onSrStartedHandler,
		SR_RESULT_CHG_NAME: onSrResultChangedHandler,
		SR_COMPLETED_NAME:  onSrCompletedHandler,
		TASK_FAILED_NAME:   onSrTaskFailedHandler,
	},
}

func newSpeechRecognitionProto() *commonProto {
	return &srProto
}

func NewSpeechRecognition(config *ConnectionConfig,
	logger *NlsLogger,
	taskfailed func(string, interface{}),
	started func(string, interface{}),
	resultchanged func(string, interface{}),
	completed func(string, interface{}),
	closed func(interface{}),
	param interface{}) (*SpeechRecognition, error) {
	sr := new(SpeechRecognition)
	proto := newSpeechRecognitionProto()
	if logger == nil {
		logger = DefaultNlsLog()
	}

	nls, err := newNlsProto(config, proto, logger, sr)
	if err != nil {
		return nil, err
	}

	sr.nls = nls
	sr.UserParam = param
	sr.onTaskFailed = taskfailed
	sr.onStarted = started
	sr.onResultChanged = resultchanged
	sr.onCompleted = completed
	sr.onClose = closed
	return sr, nil
}

func (sr *SpeechRecognition) Start(param SpeechRecognitionStartParam, extra map[string]interface{}) (chan bool, error) {
	if sr.nls == nil {
		return nil, errors.New("empty nls: using NewSpeechRecognition to create a valid instance")
	}

	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(b, &sr.StartParam)
	if extra != nil {
		if sr.StartParam == nil {
			sr.StartParam = extra
		} else {
			for k, v := range extra {
				sr.StartParam[k] = v
			}
		}
	}
	sr.taskId = getUuid()
	err = sr.nls.Connect()
	if err != nil {
		return nil, err
	}

  sr.lk.Lock()
  defer sr.lk.Unlock()
	sr.startCh = make(chan bool, 1)
	return sr.startCh, nil
}

func (sr *SpeechRecognition) Stop() (chan bool, error) {
	if sr.nls == nil {
		return nil, errors.New("empty nls: using NewSpeechRecognition to create a valid instance")
	}


	req := CommonRequest{}
	req.Context = DefaultContext
	req.Header.Appkey = sr.nls.connConfig.Appkey
	req.Header.MessageId = getUuid()
	req.Header.Name = SR_STOP_NAME
	req.Header.Namespace = SR_NAMESPACE
	req.Header.TaskId = sr.taskId

	b, _ := json.Marshal(req)
	err := sr.nls.cmd(string(b))
	if err != nil {
		return nil, err
	}

  sr.lk.Lock()
  defer sr.lk.Unlock()
	sr.stopCh = make(chan bool, 1)
	return sr.stopCh, nil
}

func (sr *SpeechRecognition) Shutdown() {
	if sr.nls == nil {
		return
	}

	sr.nls.shutdown()

	sr.lk.Lock()
	defer sr.lk.Unlock()
	if sr.startCh != nil {
		sr.startCh <- false
		close(sr.startCh)
		sr.startCh = nil
	}

	if sr.stopCh != nil {
		sr.stopCh <- false
		close(sr.stopCh)
		sr.stopCh = nil
	}
}

func (sr *SpeechRecognition) SendAudioData(data []byte) error {
	if sr.nls == nil {
		return errors.New("empty nls: using NewSpeechRecognition to create a valid instance")
	}

	return sr.nls.sendRawData(data)
}
