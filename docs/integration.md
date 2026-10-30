# 接入指南

从空目录到能稳定跑在生产里的同花顺数据通道。方法级字段清单见
[API 参考](api-reference.md)，本文只讲怎么接、怎么跑对。

## 1. 它到底是什么

`go-thsdk` 不是重新实现了一个行情协议。上游 Python 包 `thsdk` 本身也只是封装，
真正干活的是一个原生动态库（`hq.dll`，内部同样是 Go 写的），
对外只暴露一个 JSON 调用入口：

```
你的 Go 程序
   │  struct → JSON
   ▼
internal/native ──CallAlloc(json)──▶ hq.*  ──TCP──▶  同花顺行情服务器
                    ◀─FreeResult(ptr──       9602 等端口
```

因此：**数据权限、限频规则、账号体系、可查到的市场范围，都和 Python 版完全一致**。
本仓库不含这个动态库，需要自己取（见下）。

## 2. 准备工作

### 2.1 取运行组件

组件不在 GitHub 仓库里（上游同步流程会主动删掉 `libs/**`），只随 PyPI 源码包分发：

```bash
pip download --no-deps --no-binary :all: thsdk -d ./_pkg
tar -xzf _pkg/thsdk-*.tar.gz -C _pkg
mkdir -p libs
cp -r _pkg/thsdk-*/thsdk/libs/* ./libs/   # 整包复制，SDK 按当前平台取用
```

| 平台 | 文件 |
| --- | --- |
| Windows amd64 | `libs/windows/amd64/hq.dll` |

> 当前只在 Windows amd64 上实测过。

### 2.2 让 SDK 找到它

按顺序命中第一个：

1. `Options.RuntimePath`
2. 环境变量 `THSDK_RUNTIME_PATH`
3. 可执行文件同级：`hq.*`、`libs/<os>/<arch>/hq.*`
4. 工作目录同级：`libs/<os>/<arch>/hq.*`、`hq.*`
5. 裸文件名交给系统加载器（`PATH` / `LD_LIBRARY_PATH`）

部署时建议**显式传 `RuntimePath`**，别依赖 cwd——守护进程的工作目录经常不是你以为的那个。

### 2.3 验证连通

这一步能同时验掉「组件加载」「ABI 握手」「真实出网」三件事：

```bash
go run ./examples/quickstart  # 登录 + 拉真实 K 线
```

握手失败时的报错会直接告诉你是哪一层的问题：

| 报错 | 原因 |
| --- | --- |
| `运行组件加载失败` | 路径不对 / 平台架构不匹配 |
| `ABI 版本 … 不兼容` | 组件版本与 SDK 不同源 |
| `请求协议过旧` | 组件是旧版，重新取包 |

## 3. 登录

两种方式，按你的场景选一个：

```go
// A. 自动：环境变量 → 本地会话 → 试用账号池
client.Auth(ctx, thsdk.Credentials{})

// B. 自己的账号
client.Auth(ctx, thsdk.Credentials{Username: "xxx", Password: "yyy"})
```

要点：

* **A 的兜底是上游打包的 52 个共享试用账号**，随机取一个。它的权限较窄，需要额外
  权益的接口会返回 `permission_denied`。不想用它就设
  `Options.DisableTemporaryAccount = true`。
* **会话文件**写在 `Options.SessionDir`（默认当前工作目录）下的 `account.session`，
  是敏感凭据：加进 `.gitignore`，别打进镜像。`Logout` 会删掉它。
* **设备 MAC 由账号名派生**（SHA-256 后设本地管理位），与 Python 算法一致；换机器
  重登不会因 MAC 变化触发风控。要固定就传 `Options.Mac`。
* **预算来自 context**：`Auth` 用 `Options.AuthTimeout`（默认 60s）或 `ctx` 的 deadline。

生产上更稳的做法是启动时登录一次，之后复用同一个 `Client`：

```go
client, err := thsdk.New(thsdk.Options{
    SessionDir: `C:\ProgramData\myapp\ths`,
    RuntimePath: `C:\ths\hq.dll`,
    DisableTemporaryAccount: true,
})
```

## 4. 两种调用姿势

### 4.1 便捷层：拿现成的行

```go
rows, err := client.GetPrice(ctx, thsdk.PriceQuery{
    Security:  thsdk.MustSecurity("USHA600519"),
    Count:     20,
    Frequency: "daily",
    FQ:        "pre",
})
// rows[i]["close"] 是 json.RawMessage，键与组件返回一致
```

适合快速取数、字段不固定的场景（问财、快讯这类动态列本来就只能是 map）。

### 4.2 标准层：强类型

```go
data, err := client.ListSecurityKLines(ctx, thsdk.ListSecurityKLinesRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Start:    -20, End: 0,
    Adjust:   thsdk.AdjustForward,
    Period:   thsdk.PeriodDay,
})

var klines []thsdk.KLine
if err := data.Rows(&klines); err != nil { … }
```

80 个方法一一对应原生方法，签名和字段见 [API 参考](api-reference.md)。

**`Rows` / `Into` / `Metadata` 怎么选**

| 返回形态 | 用法 |
| --- | --- |
| 顶层数组（`list[Xxx]`） | `data.Rows(&[]Xxx{})` |
| 对象 + 行集键（`rows`/`items`/`securities`/…） | 同样 `data.Rows(...)`，SDK 自己找行集 |
| 对象且要读游标 | `data.Metadata()` 拿除行集外的键 |
| 单个对象（如 `ResolveBlock`） | `data.Into(&block)` |

**可选字段用 `Opt`**

```go
SortBegin: thsdk.Num(0),   // 发送 "sort_begin": 0
SortCount: thsdk.Num(20),
// 不赋值 = 该键完全不出现在请求里
```

`Num` / `Text` / `Flag` 三选一即可。不要写 `SortBegin: thsdk.Num(nil)`。

## 5. 数据约定（已实测）

### 5.1 证券标识

完整代码 = **4 位字母市场 + 代码**，例如 `USHA600519`。市场前缀实测可用：

| 前缀 | 覆盖 | 样例 |
| --- | --- | --- |
| `USHA` | 沪 A | `USHA600519` |
| `USZA` | 深 A（主板 000/002、创业板 300） | `USZA000001`、`USZA300750` |
| `USHD` | 债券 | `USHD110100` |
| `UHKM` | 港股 | `UHKMHK0700` |
| `UNQQ` / `UNYC` | 美股 / 美股 ETF | `UNQQAAPL`、`UNYCSMQ` |
| `USHJ` | 场内基金 | `USHJ513300` |
| `USHI` | 指数 | `USHI1B0006` |
| `URFI` | 同花顺行业/概念指数 | `URFI881273` |

这张表不必背——拿不准时先补全再调用，别猜。

```go
rows, err := client.CompleteTHSCodes(ctx, "600519", "AAPL", "00700")
// 返回 market / code / full_code
```

注意 `resolve_securities` **只吃短代码**，传 `USHA600519` 反而匹配不到。

### 5.2 时间字段——最容易踩的坑

| 字段 | 单位 | 说明 |
| --- | --- | --- |
| K 线 `time`（日及以上周期） | `YYYYMMDD` 整数 | `KLine.TradeDate()` 直接转 `time.Time` |
| `tick.time`、`timeline.time` | **Unix 秒** | 实测 `1793322900` → 09:15:00，正是开盘 |
| 分时 `intraday.time` / `time_index` | **组件内部分钟序号，不是时间戳** | 不能直接格式化 |
| `trade_date` | `YYYYMMDD` 整数 | |

分时那条要说清楚：相邻行确实每分钟 +1，但盘中会跳变（实测到一处 +5），所以它
不是墙钟。两种正确用法：

```go
// 要分钟对齐：让 SDK 按 09:30 起、跳过午休，回填真实钟点
bars, err := client.IntradayBarsClock(ctx, sec, time.Time{})
// bars[i].At = 09:30, 09:31, … 11:29, 13:00, …

// 要精确时间：用逐笔或交易时间线，它们的 time 是 Unix 秒
ticks, err := client.Ticks(ctx, thsdk.TickQuery{Security: sec, Count: 500})
```

`IntradayBarsClock` 不建模半日市；港股/美股请改走逐笔。

### 5.3 数值与哨兵

* `volume` / `turnover`：整数，随周期语义变化（日 K 是当日总量）。
* `lead_indicator` 之类可能为空的小整数，组件用 **无符号 32 位 `4294967295`（即 `-1`）**
  表示「无值」，别当真实读数。
* 价格已含复权结果，组件不会给未复权原值加标记——复权方式由你请求时的 `adjust` 决定。

### 5.4 日期参数格式混用

这点和上游一致，只能按方法查：

| 形式 | 用在 |
| --- | --- |
| 整数 `YYYYMMDD` | K 线 `start` / `end`（相对条数时为负数） |
| 字符串 `"YYYYMMDD"` | `date`、`trade_date`、`time_index_day`（`"0"` = 今天） |
| ISO-8601 带时区 | `window.start` / `window.end`、集合竞价 `date` |

## 6. 限流与性能

实测出来的三条硬约束：

* **任意两次调用最少间隔 50ms**，超了返回 `rate_limited: 兄弟，太快啦。。。请保持至少 50ms 的调用间隔`。
  这是**按进程**的，不是按 `Client`——所以 SDK 把节流状态挂在共享的组件句柄上，
  同进程多 `Client`、多 goroutine 并发都会被排队，无需你自己加 sleep。
* **`rate_limited` 且 `retry_after_ms ≤ 1000` 时自动重试一次**（与 Python 版一致）；
  再失败就把错误抛给你。
* **请求体上限 1 MiB，嵌套 ≤ 64 层**，字段名必须 `snake_case`——这些在本地拦，
  不消耗限流额度。

吞吐参考：实测 70 个只读方法全扫一遍约 **65 秒**（`go run ./examples/allmethods`），
其中一半时间花在排队。要更快只能减少调用次数——批量传多只证券，不能绕过 50ms。

`ctx` 的 deadline 会换算成组件侧的 `timeout_ms`，剩余预算内重试也会重算，所以
`context.WithTimeout` 是真能兜住上层的，不会出现「deadline 过了还在等」。

## 7. 错误处理

```go
switch {
case err == thsdk.ErrNotAuthenticated:
    // 忘了登录，或 Client 已 Close
case errors.As(err, &paramErr):   // *thsdk.ParameterError
    // 本地校验失败：请求根本没发出去，不消耗额度
case errors.As(err, &apiErr):     // *thsdk.Error
    switch apiErr.Code {
    case "no_data", "encode_failed":        // 服务端此刻没数据
    case "permission_denied", "no_permission": // 账号权限不够
    case "rate_limited":                     // 撞限流（SDK 已重试过一次）
    case "timeout", "query_failed":          // 上游超时
    case "invalid_request":                  // 请求形状不对，看 API 参考改参数
    }
}
```

**怎么区分「我的错」还是「服务的错」**：`invalid_request` 和 `ParameterError` 是前者，
`no_data` / `encode_failed` / `timeout` / `query_failed` / `permission_denied` 是后者。
本项目测试对后者一律 `skip` 而非 `fail`，正是基于这个划分——实测同一时刻用上游
Python 版发同样的请求，错误码逐条一致。

## 8. 常见任务

### 增量同步日 K

```go
// 相对条数取最近 N 根，避免自己维护交易日历
rows, err := client.GetPrice(ctx, thsdk.PriceQuery{
    Security: sec, Count: 300, Frequency: "daily", FQ: "pre",
    Fields: []string{"time", "open", "high", "low", "close", "volume"},
})
```

### 批量看盘口

一次请求传多只，别循环单只——省的是限流额度：

```go
book, err := client.Depth(ctx, thsdk.MustSecurity("USHA600519"),
                               thsdk.MustSecurity("USZA300033"))
```

### 翻页

组件不帮你翻页，游标在 `Metadata` 里：

```go
page, err := client.News(ctx, 1, 50, "")      // 第 1 页
// 继续翻：用标准方法带上 metadata 里的游标字段
data, err := client.ListNewsFlashItems(ctx, thsdk.ListNewsFlashItemsRequest{
    Traversal: map[string]any{"mode": "older", /* 上游返回的游标 */},
})
```

`traversal` 这类字段是自由结构（`any`），照上游返回的形状填。

### 自然语言选股

```go
sel, err := client.Wencai(ctx, "市盈率小于20，总市值大于500亿", 50)
// sel.Rows[i] 里固定有 full_code / code / market，其余列是问财动态字段名
for _, row := range sel.Rows {
    fmt.Println(row["full_code"], row["总市值[20261029]"])
}
```

动态列名带日期后缀（如 `市盈率(pe)[20261029]`），别硬编码列名，先读 `sel.Fields`。

## 9. 可测试性

`*Client` 是具体类型，要 mock 就抽接口，只暴露你用到的方法：

```go
type quotes interface {
    GetPrice(context.Context, thsdk.PriceQuery) ([]thsdk.PriceRow, error)
}
```

不联网也能测的部分已经内置：请求编码契约有离线测试（`wire_contract_test.go`），
分时钟点映射、字段校验、`Opt` 省略语义都不需要组件。

需要真实数据时，`native/` 下有组件就能跑，没有则自动 `skip`：

```bash
go test ./...            # 离线语义 + 有组件时附带实请求
go test ./... -short     # 只跑离线
```

## 10. 已知边界

* **没有订阅、没有 Level-2、没有任意原始方法分发**。组件有方法 allow-list，越界调用
  返回 `method_not_allowed`。这与上游 Python 版一致，不是本移植的限制。
* **组件加载后不卸载**。Go 的 c-shared 库卸载后再加载会崩，所以映射常驻进程；
  同进程反复 `New` / `Close` 是安全的，但别指望 `Close` 释放那 11MB 内存。
* **实时快照不保证可用**。盘口、分时、逐笔会返回 `encode_failed` / `no_data`，
  同一请求可能几十秒内时好时坏；历史数据接口通常不受影响。
* **试用账号池是共享的**，权限、频率、可查数据集都比个人账号窄。
* 交易日历、停牌、退市等状态**没有专门接口**，只能从返回的行数与 `trade_date` 反推。

## 11. 自查清单

上线前对着过一遍：

- [ ] `RuntimePath` 显式指定，不依赖 cwd
- [ ] 组件文件与目标平台匹配（Windows amd64 放 `libs/windows/amd64/hq.dll`）
- [ ] `account.session` 目录可写，且不在版本库 / 镜像层里
- [ ] 生产用 `DisableTemporaryAccount: true` + 自己的账号
- [ ] 批量传多只证券，而不是循环单只
- [ ] 所有调用带 `ctx` deadline
- [ ] 错误分支区分了「本地校验」与「服务端可用」
- [ ] 分时数据没有把 `time` 当时间戳格式化
- [ ] `go run ./examples/allmethods` 在你自己的账号上跑过一遍，确认可见范围
