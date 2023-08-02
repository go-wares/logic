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
	"fmt"
	"github.com/go-wares/log"
	"github.com/kataras/iris/v12"
	"net/http"
	"sync"
)

var (
	// Manager
	// 管理实例.
	Manager Management
)

type (
	// Management
	// 管理接口.
	Management interface {
		// Add
		// 添加执行器.
		Add(name string, handler Handler) Management

		// Call
		// 调用逻辑.
		Call(i iris.Context, name string) (res Result)
	}

	// 管理结构体.
	management struct {
		pools    map[string]*sync.Pool
		registry map[string]Handler
	}
)

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

// Add
// 添加执行.
func (o *management) Add(name string, handler Handler) Management {
	o.pools[name] = &sync.Pool{}
	o.registry[name] = handler
	return o
}

// Call
// 调用逻辑.
func (o *management) Call(i iris.Context, name string) (res Result) {
	var (
		code    int
		err     error
		handler Handler
		span    = log.NewSpanFromRequest(i.Request(), i.Request().URL.Path)
	)

	span.Info("[logic=%s] 开始请求", name)

	// 1. 监听结束.
	defer func() {
		// 构建结果.
		if res == nil {
			res = NewResult()
		}

		// 捕获异常.
		if r := recover(); r != nil {
			span.Fatal("[logic=%s] 请求异常: %v", name, r)

			// 覆盖错误.
			if err == nil {
				code = http.StatusInternalServerError
				err = fmt.Errorf("%v", r)
			}
		}

		// 设置错误.
		if err != nil {
			res.SetError(code, err)
		}

		// 结束请求.
		span.Info("[logic=%s] 请求结束: %s", name, res.JSON())
		span.End()

		// 释放实例.
		if handler != nil {
			if _, ok := o.pools[name]; ok {
				o.pools[name].Put(handler)
			}
		}
	}()

	// 2. 取出逻辑.
	if _, ok := o.pools[name]; ok {
		if v := o.pools[name].Get(); v != nil {
			handler = v.(Handler)
		}
	}

	// 3. 构造逻辑.
	if handler == nil {
		if v, ok := o.registry[name]; ok {
			handler = v
		}
	}

	// 4. 无效逻辑.
	if handler == nil {
		code = http.StatusNotAcceptable
		err = ErrUnregisteredHandler
		return
	}

	// 5. 处理逻辑.
	res = handler(span.Context(), i)
	return
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *management) init() *management {
	o.pools = make(map[string]*sync.Pool)
	o.registry = make(map[string]Handler)
	return o
}
