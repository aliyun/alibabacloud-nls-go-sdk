/*
core.go

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
	"fmt"
	"time"
)

const (
	CONNECTED_HANDLER = "CONNECTED_HANDLER"
	CLOSE_HANDLER     = "CLOSE_HANDLER"
	RAW_HANDLER       = "RAW_HANDLER"
)

type ConnectionConfig struct {
	Url     string `json:"url"`
	Token   string `json:"token"`
	Akid    string `json:"akid"`
	Akkey   string `json:"akkey"`
	Appkey  string `json:"appkey"`
	Rbuffer int    `json:"rbuffer"`
	Wbuffer int    `json:"wbuffer"`
}

func NewConnectionConfigWithAKInfoDefault(url string, appkey string,
	akid string, akkey string) (*ConnectionConfig, error) {
	tokenMsg, err := GetToken(DEFAULT_DISTRIBUTE, DEFAULT_DOMAIN, akid, akkey, DEFAULT_VERSION)
	if err != nil {
		return nil, err
	}

	if tokenMsg.TokenResult.Id == "" {
		str := fmt.Sprintf("obtain empty token err:%s", tokenMsg.ErrMsg)
		return nil, errors.New(str)
	}

	return NewConnectionConfigWithToken(url, appkey, tokenMsg.TokenResult.Id), nil
}

func NewConnectionConfigWithToken(url string, appkey string, token string) *ConnectionConfig {
	config := new(ConnectionConfig)
	config.Url = url
	config.Appkey = appkey
	config.Token = token
	config.Rbuffer = 1024
	config.Wbuffer = 4096
	return config
}

func NewConnectionConfigFromJson(jsonStr string) (*ConnectionConfig, error) {
	config := ConnectionConfig{}
	err := json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		return nil, err
	}

	if config.Url == "" || config.Appkey == "" {
		return nil, errors.New("invalid connection config: no url or appkey")
	}

	if config.Token == "" {
		if config.Akid == "" || config.Akkey == "" {
			return nil, errors.New("invalid connection config: if no token provided, must provide akid and akkey")
		}
		return NewConnectionConfigWithAKInfoDefault(config.Url, config.Appkey, config.Akid, config.Akkey)
	} else {
		return NewConnectionConfigWithToken(config.Url, config.Appkey, config.Token), nil
	}
}

type nlsProto struct {
	proto      *commonProto
	conn       *wsConnection
	connConfig *ConnectionConfig
	logger     *NlsLogger
	taskId     string
	param      interface{}
}

type commonProto struct {
	namespace string
	handlers  map[string]func(isErr bool, text []byte, proto *nlsProto)
}

func newNlsProto(connConfig *ConnectionConfig,
	proto *commonProto, logger *NlsLogger, param interface{}) (*nlsProto, error) {
	if connConfig == nil || proto == nil {
		return nil, errors.New("connConfig or proto is nil")
	}
	if proto.handlers == nil {
		return nil, errors.New("invalid proto: nil handler")
	}

	nls := new(nlsProto)
	nls.connConfig = connConfig
	nls.proto = proto
	if logger == nil {
		nls.logger = DefaultNlsLog()
	} else {
		nls.logger = logger
	}

	nls.param = param
	return nls, nil
}

func (nls *nlsProto) Connect() error {
	if nls.conn != nil {
		nls.conn.shutdown()
		time.Sleep(time.Millisecond * 100)
	}

	ws, err := newWsConnection(nls.connConfig.Url,
		nls.connConfig.Token, 10*time.Second, nls.connConfig.Rbuffer,
		nls.connConfig.Wbuffer, nls.logger,
		//recv frame
		func(rawData bool, data []byte) {
			if rawData {
				handler, ok := nls.proto.handlers[RAW_HANDLER]
				if !ok {
					nls.logger.Fatal("NO RAW_HANDLER BUT recv RAW FRAME")
					return
				} else {
					handler(false, data, nls)
				}
			} else {
				nls.logger.Debugf("recv raw frame:%s", string(data))
				resp := CommonResponse{}
				err := json.Unmarshal(data, &resp)
				if err != nil {
					nls.logger.Println("OCCUR UNKNOWN PROTO:", err)
					return
				}

				if resp.Header.Namespace != "Default" && resp.Header.Namespace != nls.proto.namespace {
					nls.logger.Fatalf("WTF namespace mismatch expect %s but %s", nls.proto.namespace, resp.Header.Namespace)
					return
				}
				handler, ok := nls.proto.handlers[resp.Header.Name]
				if !ok {
					nls.logger.Printf("no handler for %s", resp.Header.Name)
					if cust_handler, ok := nls.proto.handlers[CUSTOM_DEFINED_NAME]; ok {
						nls.logger.Println("using custom handler for", resp.Header.Name)
						cust_handler(false, data, nls)
					} else {
						nls.logger.Println("no custom handler for", resp.Header.Name)
					}
					return
				}
				handler(false, data, nls)
			}
		},
		//close
		func(code int, text string, err error) {
			handler, ok := nls.proto.handlers[CLOSE_HANDLER]
			if ok {
				handler(true, []byte(text), nls)
			}
		})
	if err != nil {
		return err
	}

	nls.conn = ws
	nls.logger.Println("connect done")
	handler, ok := nls.proto.handlers[CONNECTED_HANDLER]
	if ok {
		handler(false, nil, nls)
	} else {
		nls.logger.Println("no onConnected handler")
	}

	return nil
}

func (nls *nlsProto) shutdown() error {
	if nls.conn == nil {
		return errors.New("nls proto is nil")
	}
	return nls.conn.shutdown()
}

func (nls *nlsProto) cmd(cmd string) error {
	if nls.conn == nil {
		return errors.New("nls proto is nil")
	}

	return nls.conn.sendTextData(cmd)
}

func (nls *nlsProto) sendRawData(data []byte) error {
	if nls.conn == nil {
		return errors.New("nls proto is nil")
	}

	return nls.conn.sendBinary(data)
}
