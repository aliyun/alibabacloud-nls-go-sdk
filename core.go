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

package dash

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	CONNECTED_HANDLER = "CONNECTED_HANDLER"
	CLOSE_HANDLER     = "CLOSE_HANDLER"
	RAW_HANDLER       = "RAW_HANDLER"
)

type ConnectionConfig struct {
	Url     string
	Apikey  string
	Rbuffer int
	Wbuffer int
}

func NewConnectionConfigDefault() (*ConnectionConfig, error) {
	var apikey string
	apikey = os.Getenv(APIKEY_ENV_KEY_NAME)

	if apikey == "" {
		str := fmt.Sprintf("obtain apikey from env %s failed.", APIKEY_ENV_KEY_NAME)
		return nil, errors.New(str)
	}

	return NewConnectionConfigWithUrlApiKey(DEFAULT_URL, apikey, DEFAULT_WS_RBUFFER_SIZE, DEFAULT_WS_WBUFFER_SIZE), nil
}

func NewConnectionConfigWithUrlApiKey(url string, apikey string, rbuffer int, wbuffer int) *ConnectionConfig {
	config := new(ConnectionConfig)
	config.Url = url
	config.Apikey = apikey
	config.Rbuffer = rbuffer
	config.Wbuffer = wbuffer
	return config
}

type dashProto struct {
	proto      *commonProto
	conn       *wsConnection
	connConfig *ConnectionConfig
	logger     *NlsLogger
	param      interface{}
}

type commonProto struct {
	Model     string
	TaskGroup string
	Task      string
	Function  string
	handlers  map[string]func(isErr bool, text []byte, proto *dashProto)
}

func newDashProto(connConfig *ConnectionConfig,
	proto *commonProto, logger *NlsLogger, param interface{}) (*dashProto, error) {
	if connConfig == nil || proto == nil {
		return nil, errors.New("connConfig or proto is nil")
	}
	if proto.handlers == nil {
		return nil, errors.New("invalid proto: nil handler")
	}

	nls := new(dashProto)
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

func (nls *dashProto) Connect() error {
	if nls.conn != nil {
		nls.conn.shutdown()
		time.Sleep(time.Millisecond * 100)
	}

	ws, err := newWsConnection(nls.connConfig.Url,
		nls.connConfig.Apikey, 10*time.Second, nls.connConfig.Rbuffer,
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
				nls.logger.Debugf("recv none raw frame:%s", string(data))
				resp := CommonResponse{}
				err := json.Unmarshal(data, &resp)
				if err != nil {
					nls.logger.Println("OCCUR UNKNOWN PROTO:", err)
					return
				}

				handler, ok := nls.proto.handlers[resp.Header.Event]
				if !ok {
					nls.logger.Printf("no handler for %s", resp.Header.Event)
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

func (nls *dashProto) shutdown() error {
	if nls.conn == nil {
		return errors.New("nls proto is nil")
	}
	return nls.conn.shutdown()
}

func (nls *dashProto) cmd(cmd string) error {
	if nls.conn == nil {
		return errors.New("nls proto is nil")
	}

	return nls.conn.sendTextData(cmd)
}

func (nls *dashProto) sendRawData(data []byte) error {
	if nls.conn == nil {
		return errors.New("nls proto is nil")
	}

	return nls.conn.sendBinary(data)
}
