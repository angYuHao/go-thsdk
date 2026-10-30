# go-thsdk

[![Go Reference](https://pkg.go.dev/badge/github.com/angYuHao/go-thsdk.svg)](https://pkg.go.dev/github.com/angYuHao/go-thsdk)
[![Go Version](https://img.shields.io/github/go-mod/go-version/angYuHao/go-thsdk)](go.mod)
[![Latest Tag](https://img.shields.io/github/v/tag/angYuHao/go-thsdk)](https://github.com/angYuHao/go-thsdk/tags)
![Platform](https://img.shields.io/badge/platform-Windows%20amd64-blue)

同花顺行情数据的纯 Go SDK：直接调用原生组件 `hq.dll` 的 JSON ABI，覆盖 80 个原生
方法——实时与历史行情、证券与板块、资讯事件、排行统计、问财选股与自选股管理；
既能作为标准库被 `import`，也能一行命令启动成 HTTP 服务。

上游 [panghu11033/thsdk](https://github.com/panghu11033/thsdk) 是 Python 层封装，
数据能力来自随包分发的 c-shared 组件（`hq.dll`）。本项目用纯 Go 直接调用该组件的
`CallAlloc` / `FreeResult` JSON ABI，不依赖 Python，数据权限、限频规则和会话文件
格式与 Python 版保持一致。

```
Go SDK ──JSON──▶ hq.dll ──▶ 同花顺行情网关
```

**当前实测平台：Windows amd64。** 80 个原生方法全部发过真实请求：70 个只读方法
用真实账号验证，10 个账号写方法在试用账号上逐条回读校验。

## 目录

- [能力范围](#能力范围)
- [实测结果](#实测结果)
- [环境要求](#环境要求)
- [安装](#安装)
- [快速开始](#快速开始)
- [认证](#认证)
- [标准层](#标准层)
- [便捷层](#便捷层)
- [HTTP 服务](#http-服务)
- [数据约定（已实测）](#数据约定已实测)
- [进程行为](#进程行为)
- [已知的上游限制](#已知的上游限制)
- [示例](#示例)
- [测试](#测试)
- [文档](#文档)
- [版本](#版本)
- [许可](#许可)

## 能力范围

覆盖组件提供的全部原生方法，分为六类：

| 类别 | 方法数 | 内容 |
| --- | ---: | --- |
| 行情 | 17 | K 线、分时、逐笔、盘口、成交明细、基本面快照 |
| 证券与板块 | 18 | 证券搜索/解析、板块成分、行业映射、A/H 关系 |
| 资讯与事件 | 14 | 快讯、新闻、公告、异动、涨停分析 |
| 排行与统计 | 11 | 涨跌幅/换手率排行、实时统计、成分贡献 |
| 问财 | 7 | 自然语言选股、表达式字段、实时字段 |
| 账号 | 13 | 权限查询、自选股与分组读写 |

## 实测结果

一次完整扫描（真实账号，盘中）：

| 结果 | 数量 | 说明 |
| --- | ---: | --- |
| OK | 61 | 返回真实数据 |
| 服务侧 | 8 | `no_data` / `encode_failed` / `permission_denied` |
| 失败 | 1 | `list_market_call_auction_unusual_events` 超时 |

账号写方法 10 个全部通过：临时分组的建/增/改名/整组替换/删，默认自选股的加/删/
同内容替换，以及 `clear` 清空后恢复。每一步都回读校验，结束时数据与调用前一致。

## 环境要求

- Windows amd64
- Go 1.26 或更高
- `hq.dll`（不在本仓库中，见下一节）

## 安装

```bash
go get github.com/angYuHao/go-thsdk
```

`hq.dll` 只随 PyPI 源码包分发，上游同步流程会删除 GitHub 仓库里的 `libs/**`：

```bash
pip download --no-deps --no-binary :all: thsdk -d ./_pkg
tar -xzf _pkg/thsdk-*.tar.gz -C _pkg
mkdir -p native
cp _pkg/thsdk-*/thsdk/libs/windows/amd64/hq.dll ./native/
```

SDK 按以下顺序查找组件，命中第一个即用：

1. `Options.RuntimePath`
2. 环境变量 `THSDK_RUNTIME_PATH`
3. 可执行文件同级：`hq.dll`、`libs/windows/amd64/hq.dll`、`native/hq.dll`
4. 工作目录同级：`native/hq.dll`、`hq.dll`、`libs/windows/amd64/hq.dll`
5. 交给系统加载器（`PATH`）

部署时建议显式传 `RuntimePath`，不要依赖工作目录。

## 快速开始

```go
package main

import (
	"context"
	"fmt"

	thsdk "github.com/angYuHao/go-thsdk"
)

func main() {
	client, err := thsdk.New(thsdk.Options{})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Auth(ctx, thsdk.Credentials{}); err != nil {
		panic(err)
	}

	data, err := client.ListSecurityKLines(ctx, thsdk.ListSecurityKLinesRequest{
		Security: thsdk.MustSecurity("USHA600519"),
		Start:    -20,
		End:      0,
		Adjust:   thsdk.AdjustForward,
		Period:   thsdk.PeriodDay,
	})
	if err != nil {
		panic(err)
	}

	var klines []thsdk.KLine
	if err := data.Rows(&klines); err != nil {
		panic(err)
	}
	for _, bar := range klines {
		fmt.Println(bar.Time, bar.Close)
	}
}
```

## 认证

`Auth` 按顺序尝试三种凭据，任一成功即返回：

1. 显式传入的 `Credentials{Username, Password}`
2. 环境变量 `THS_USERNAME` / `THS_PASSWORD`
3. 本地会话文件 `account.session`
4. 内置试用账号（设置 `DisableTemporaryAccount: true` 可禁用）

```bash
export THS_USERNAME=你的账号
export THS_PASSWORD=你的密码
```

会话文件写在 `Options.SessionDir`（默认工作目录），**按凭据对待，不要提交**；
`Logout` 会删除它。设备 MAC 由账号名派生（SHA-256 + 本地管理位），可用
`Options.Mac` 覆盖。

## 标准层

80 个方法与原生组件同名同参，返回 `RawData`：

- 必填字段是具体类型（`Security`、`int`、`string` 等），缺失时本地直接报错。
- 可选字段是 `Opt`，用 `thsdk.Num(20)` / `thsdk.Text("D")` / `thsdk.Flag(true)`
  赋值；不赋值就不发送，显式零值也能表达。
- `RawData.Rows` 把裸数组或 `rows` / `items` / `securities` / `events` / `tags` /
  `sections` / `groups` 提升为行集；`Into` 解整个对象；`Metadata` 取分页游标等
  非行集字段。

## 便捷层

对常用查询做了薄封装，下面是已实测的用法：

```go
rows, err := client.GetPrice(ctx, thsdk.PriceQuery{
	Security:  thsdk.MustSecurity("USHA600519"),
	Frequency: "daily",
	Start:     time.Now().AddDate(-1, 0, 0).In(thsdk.Shanghai),
	End:       time.Now().In(thsdk.Shanghai),
	Fields:    []string{"time", "open", "high", "low", "close", "volume"},
})

pick, err := client.Wencai(ctx, "市盈率小于20，总市值大于500亿", 50)
news, err := client.News(ctx, 1, 20, "")
```

## 数据约定（已实测）

- 日线及以上 `KLine.Time` 是 `YYYYMMDD` 整数；分时索引不是时间戳。
- `tick.time` 与 `timeline.time` 是 Unix 秒（09:15 的集合竞价能对上）。
- 市场前缀大写：`USHA` 沪 A、`USZA` 深 A/创业板、`USHT`、`URFI` 行业/概念指数、
  `USHI` 指数、`UHKM` 港股、`UNQQ`/`UNYC` 美股。
- `lead_indicator` 的 `4294967295` 是 `-1` 的无符号哨兵值。

## 进程行为

- 组件按进程限流，提示「请保持至少 50ms 的调用间隔」；SDK 在 `mappedLibrary`
  上补齐间隔，并发调用会排队。
- `rate_limited` 且 `retry_after_ms ≤ 1000` 时自动重试一次。
- Go c-shared 库卸载后无法安全重载，因此 `Close` 只释放句柄，映射保留到进程结束；
  同一进程内多次 `New`/`Close` 是安全的。

## 已知的上游限制

这些现象来自实测，错误由组件或服务端返回，与移植无关（同一时刻 Python 版错误码
一致）：

| 现象 | 原因 |
| --- | --- |
| `list_market_call_auction_unusual_events` 60s 超时 | `check_market_timeline_readiness` 报 `ready:false`（`timeindexinit=fail`），时间轴未初始化 |
| `extended_hours_*` 返回 `permission_denied` | 账号没有扩展行情权益（`a_stock_extended_level2:false`），港股同样受限 |
| `list_security_call_auction_quotes` 对部分历史日期返回 `encode_failed` | 上游对该日响应编码失败；换当天日期可正常返回 |
| 盘口、逐笔、依赖当日快照的接口间歇 `no_data` / `encode_failed` | 非交易时段或当日数据未就绪 |

## HTTP 服务

同一份代码既能当标准库用，也能直接跑成一个 HTTP 服务。服务端只是 `thsdk.Client`
的一层薄包装，路由可直接挂到已有应用里。

```bash
go run ./cmd/thsdk-server                        # 默认 127.0.0.1:8765
go run ./cmd/thsdk-server -addr :9000 -allow-writes
```

| 端点 | 说明 |
| --- | --- |
| `GET /healthz` | 存活状态、是否已认证、当前可用方法数 |
| `GET /v1/methods` | 方法目录（名称、类别、是否写方法） |
| `POST /v1/call/<method>` | 调用一个方法，请求体即该方法的参数对象 |

```bash
curl http://127.0.0.1:8765/healthz

curl -X POST http://127.0.0.1:8765/v1/call/list_security_k_lines \
  -d '{"security":{"code":"600519","market":"USHA"},"start":-3,"end":0,"adjust":1,"period":16384}'
```

请求体为空或 `null` 表示无参数。响应直接是组件返回的 `data`，不做包装。
错误统一为 `{"error":{"code":...,"message":...}}`，HTTP 状态按类型映射：

| 上游错误 | HTTP |
| --- | --- |
| `invalid_request` / 本地参数错误 | 400 |
| 未知方法 | 404 |
| `permission_denied` / `no_permission` | 403 |
| 写方法未开启 | 403 |
| `no_data` | 404 |
| `rate_limited` | 429 |
| `timeout` / `query_failed` | 504 |
| 未认证 / `session_unavailable` | 503 |
| 其它组件业务错误 | 502 |

**写方法默认关闭**：10 个账号写方法只有在 `-allow-writes`（库层为
`Config.AllowWrites`）下才会注册，否则返回 403。服务默认只监听 `127.0.0.1`，
对外暴露前请自行加上认证和 TLS。

在已有应用里挂载：

```go
handler, err := httpserver.New(httpserver.Config{
    Client:      client,          // 已认证的 *thsdk.Client
    AllowWrites: false,
    CallTimeout: 30 * time.Second,
})
if err != nil {
    return err
}
mux.Handle("/thsdk/", http.StripPrefix("/thsdk", handler.Handler()))
```

命令行参数：`-addr`、`-runtime`、`-session-dir`、`-allow-writes`、`-no-temporary`、
`-call-timeout`、`-auth-timeout`、`-shutdown-timeout`。

## 示例

```bash
go run ./examples/quickstart     # 最小登录 + 取数
go run ./examples/typed          # 强类型请求与结果
go run ./examples/convenience    # GetPrice / Wencai / News
go run ./examples/allmethods     # 70 个只读方法全量扫描
go run ./examples/allmethods -writes  # 追加 10 个账号写方法（会改动试用账号）
```

## 测试

```bash
go build ./...
go test -short ./...   # 离线契约测试与 HTTP 路由测试，不访问服务
go test ./...          # 追加真实请求（需要 native/hq.dll）
```

CI 在 `windows-latest` 上跑 `gofmt`、`go build ./...` 与 `go test -short ./...`，
见 [.github/workflows/ci.yml](.github/workflows/ci.yml)。离线用例不需要 `hq.dll`，
所以 CI 不需要凭据也不会外呼。

## 文档

- [接入指南](docs/integration.md)：组件准备、登录、数据约定、限流、错误分类、上线自查
- [API 参考](docs/api-reference.md)：80 个标准方法与 12 个便捷方法的签名与返回形态

## 版本

| 版本 | 内容 |
| --- | --- |
| `v0.1.0` | 首个可用版本：80 个原生方法、便捷层、接入指南与 API 参考、HTTP 服务 |

查看全部标签：<https://github.com/angYuHao/go-thsdk/tags>

版本号遵循语义化版本：`v0.x` 期间 API 仍可能微调，`v1.0.0` 之后不再破坏兼容。

## 许可

本移植沿用上游 `thsdk` 的许可与数据服务条款。原生组件 `hq.dll` 来自 PyPI 发行包，
不在本仓库中；其使用条款与合规性需自行向上游和同花顺确认。
