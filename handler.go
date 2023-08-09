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
	"context"
	"github.com/kataras/iris/v12"
)

type (
	// Handler
	// 执行器回调.
	Handler func() HandlerRunner

	// HandlerRunner
	// 执行器接口.
	HandlerRunner interface {
		// Clean
		// 清理数据.
		Clean(ctx context.Context)

		// Ready
		// 准备处理.
		Ready(ctx context.Context)

		// Run
		// 执行逻辑.
		Run(ctx context.Context, i iris.Context) (res Result)
	}
)
