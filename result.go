// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-08-02

package logic

import (
	"encoding/json"
)

const (
	defaultResultCode = 0
	defaultResultText = "success"
)

type (
	// Result
	// 逻辑处理结果.
	Result interface {
		JSON() string
		Success() bool

		SetData(data interface{}) Result
		SetError(code int, err error) Result
	}

	// 结果结构体.
	result struct {
		Code int         `json:"code" yaml:"code"`
		Data interface{} `json:"data" yaml:"data"`
		Text string      `json:"text" yaml:"text"`
	}
)

// NewResult
// 创建结果实例.
func NewResult() Result {
	return (&result{
		Code: defaultResultCode,
		Text: defaultResultText,
	}).init()
}

func (o *result) JSON() string {
	if body, err := json.Marshal(o); err == nil {
		return string(body)
	}
	return "{}"
}

func (o *result) Success() bool { return defaultResultCode == o.Code }

// +---------------------------------------------------------------------------+
// | Set result fields                                                         |
// +---------------------------------------------------------------------------+

func (o *result) SetData(data interface{}) Result {
	o.Data = data
	return o
}

func (o *result) SetError(code int, err error) Result {
	o.Code = code
	o.Data = map[string]interface{}{}
	o.Text = err.Error()
	return o
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *result) init() *result {
	return o
}
