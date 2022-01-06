/*
token.go

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
  "github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func GetToken(dist string, domain string, akid string, akkey string, version string) (*TokenResultMessage, error) {
	client, err := sdk.NewClientWithAccessKey(dist, akid, akkey)
	if err != nil {
		return nil, err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Domain = domain
	request.ApiName = "CreateToken"
	request.Version = version
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}

	message := new(TokenResultMessage)
	err = json.Unmarshal(response.GetHttpContentBytes(), message)
	if err != nil {
		return nil, err
	}

	return message, nil
}


