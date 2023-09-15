# NLS Go SDK 说明

> 本文介绍如何使用阿里云智能语音服务提供的 Go SDK，包括 SDK 的安装方法及 SDK 代码示例。

## 前提条件

使用 SDK 前，请先阅读接口说明，详细请参见**接口说明**。

### 下载安装

> 说明
>
> - SDK 支持 go1.16
> - 请确认已经安装 golang 环境，并完成基本配置

1. 下载 SDK

通过以下命令完成 SDK 下载和安装：

> go get github.com/ikmak/alibabacloud-nls-go-sdk

2. 导入 SDK

在代码中通过将以下字段加入 import 来导入 SDK：

> import ("github.com/ikmak/alibabacloud-nls-go-sdk")

## SDK 常量

| 常量               | 常量含义                                                           |
| ------------------ | ------------------------------------------------------------------ |
| SDK_VERSION        | SDK 版本                                                           |
| PCM                | pcm 音频格式                                                       |
| WAV                | wav 音频格式                                                       |
| OPUS               | opus 音频格式                                                      |
| OPU                | opu 音频格式                                                       |
| DEFAULT_DISTRIBUTE | 获取 token 时使用的默认区域，"cn-shanghai"                         |
| DEFAULT_DOMAIN     | 获取 token 时使用的默认 URL，"nls-meta.cn-shanghai.aliyuncs.com"   |
| DEFAULT_VERSION    | 获取 token 时使用的协议版本，"2019-02-28"                          |
| DEFAULT_URL        | 默认公有云 URL，"wss://nls-gateway.cn-shanghai.aliyuncs.com/ws/v1" |

## SDK 日志

### 1. func DefaultNlsLog() \*NlsLogger

> 用于创建全局唯一的默认日志对象，默认日志以 NLS 为前缀，输出到标准错误

参数说明：

无

返回值：

NlsLogger 对象指针

### 2. func NewNlsLogger(w io.Writer, tag string, flag int) \*NlsLogger

> 创建一个新的日志

参数说明：

| 参数 | 类型      | 参数说明                             |
| ---- | --------- | ------------------------------------ |
| w    | io.Writer | 任意实现 io.Writer 接口的对象        |
| tag  | string    | 日志前缀，会打印到日志行首部         |
| flag | int       | 日志 flag，具体参考 go 官方 log 文档 |

返回值：

NlsLogger 对象指针

### 3. func (logger \*NlsLogger) SetLogSil(sil bool)

> 设置日志是否输出到对应的 io.Writer

参数说明:

| 参数 | 类型 | 参数说明                      |
| ---- | ---- | ----------------------------- |
| sil  | bool | 是否禁止日志输出，true 为禁止 |

返回值：

无

### 4. func (logger \*NlsLogger) SetDebug(debug bool)

> 设置是否打印 debug 日志，仅影响通过 Debugf 或 Debugln 进行输出的日志

参数说明：

| 参数  | 类型 | 参数说明                             |
| ----- | ---- | ------------------------------------ |
| debug | bool | 是否允许 debug 日志输出，true 为允许 |

返回值：

无

### 5. func (logger \*NlsLogger) SetOutput(w io.Writer)

> 设置日志输出方式

参数说明：

| 参数 | 类型      | 参数说明                      |
| ---- | --------- | ----------------------------- |
| w    | io.Writer | 任意实现 io.Writer 接口的对象 |

返回值：

无

### 6. func (logger \*NlsLogger) SetPrefix(prefix string)

> 设置日志行的标签

参数说明：

| 参数   | 类型   | 参数说明                       |
| ------ | ------ | ------------------------------ |
| prefix | string | 日志行标签，会输出在日志行行首 |

返回值：

无

### 7. func (logger \*NlsLogger) SetFlags(flags int)

> 设置日志属性

参数说明：

| 参数  | 类型 | 参数说明                                         |
| ----- | ---- | ------------------------------------------------ |
| flags | int  | 日志属性，见https://pkg.go.dev/log#pkg-constants |

返回值：

无

### 8. 日志打印

日志打印方法：

| 方法名                                                       | 方法说明                                                         |
| ------------------------------------------------------------ | ---------------------------------------------------------------- |
| func (l \*NlsLogger) Print(v ...interface{})                 | 标准日志输出                                                     |
| func (l \*NlsLogger) Println(v ...interface{})               | 标注日志输出，行尾自动换行                                       |
| func (l \*NlsLogger) Printf(format string, v ...interface{}) | 带 format 的日志输出，format 方式见 go 官方文档                  |
| func (l \*NlsLogger) Debugln(v ...interface{})               | debug 信息日志输出，行尾自动换行                                 |
| func (l \*NlsLogger) Debugf(format string, v ...interface{}) | 带 format 的 debug 信息日志输出                                  |
| func (l \*NlsLogger) Fatal(v ...interface{})                 | 致命错误日志输出，输出后自动进程退出                             |
| func (l \*NlsLogger) Fatalln(v ...interface{})               | 致命错误日志输出，行尾自动换行，输出后自动进程退出               |
| func (l \*NlsLogger) Fatalf(format string, v ...interface{}) | 带 format 的致命错误日志输出，输出后自动进程退出                 |
| func (l \*NlsLogger) Panic(v ...interface{})                 | 致命错误日志输出，输出后自动进程退出并打印崩溃信息               |
| func (l \*NlsLogger) Panicln(v ...interface{})               | 致命错误日志输出，行尾自动换行，输出后自动进程退出并打印崩溃信息 |
| func (l \*NlsLogger) Panicf(format string, v ...interface{}) | 带 format 的致命错误日志输出，输出后自动进程退出并打印崩溃信息   |

## 获取 token

### 1. func GetToken(dist string, domain string, akid string, akkey string, version string) (\*TokenResultMessage, error)

> 获取访问 token

参数说明：

| 参数    | 类型   | 参数说明                                     |
| ------- | ------ | -------------------------------------------- |
| dist    | string | 区域，如果不确定，请使用 DEFAULT_DISTRIBUTE  |
| domain  | string | URL，如果不确定，请使用 DEFAULT_DOMAIN       |
| akid    | string | 阿里云 accessid                              |
| akkey   | string | 阿里云 accesskey                             |
| version | string | 协议版本，如果不确定，请使用 DEFAULT_VERSION |

返回值：

TokenResultMessage 对象指针和错误信息

## 建立连接

### 1. ConnectionConfig

> 用于建立连接的基础参数

参数说明：

| 参数   | 类型   | 参数说明                                           |
| ------ | ------ | -------------------------------------------------- |
| Url    | string | 访问的公有云 URL，如果不确定，可以使用 DEFAULT_URL |
| Token  | string | 通过 GetToken 获取的 token 或者测试 token          |
| Akid   | string | 阿里云 accessid                                    |
| Akkey  | string | 阿里云 accesskey                                   |
| Appkey | string | appkey，可以在控制台中对应项目上看到               |

### 2. func NewConnectionConfigWithAKInfoDefault(url string, appkey string, akid string, akkey string) (\*ConnectionConfig, error)

> 通过 url，appkey，akid 和 akkey 创建连接参数，等效于先调用 GetToken 然后再调用 NewConnectionConfigWithToken

参数说明：

| 参数   | 类型   | 参数说明                                           |
| ------ | ------ | -------------------------------------------------- |
| Url    | string | 访问的公有云 URL，如果不确定，可以使用 DEFAULT_URL |
| Appkey | string | appkey，可以在控制台中对应项目上看到               |
| Akid   | string | 阿里云 accessid                                    |
| Akkey  | string | 阿里云 accesskey                                   |

返回值：

\*ConnectionConfig：连接参数对象指针，用于后续创建语音交互实例

error：异常对象，为 nil 则无异常

### 3. func NewConnectionConfigWithToken(url string, appkey string, token string) \*ConnectionConfig

> 通过 url，appkey 和 token 创建连接参数

参数说明：

| 参数   | 类型   | 参数说明                                           |
| ------ | ------ | -------------------------------------------------- |
| Url    | string | 访问的公有云 URL，如果不确定，可以使用 DEFAULT_URL |
| Appkey | string | appkey，可以在控制台中对应项目上看到               |
| Token  | string | 已经通过 GetToken 或其他方式获取的 token           |

返回值：

\*ConnectionConfig：连接参数对象指针

### 4. func NewConnectionConfigFromJson(jsonStr string) (\*ConnectionConfig, error)

> 通过 json 字符串来创建连接参数

参数说明

| 参数    | 类型   | 参数说明                                                                                                                                         |
| ------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| jsonStr | string | 描述连接参数的 json 字符串，有效字段如下：url，token，akid，akkey，appkey。其中必须包含 url 和 appkey，如果包含 token 则不需要包含 akid 和 akkey |

返回值：

\*ConnectionConfig：连接对象指针

## 语音合成

### 1. SpeechSynthesisStartParam

参数说明:

| 参数           | 类型   | 参数说明                       |
| -------------- | ------ | ------------------------------ |
| Voice          | string | 发音人，默认“xiaoyun”          |
| Format         | string | 音频格式，默认使用 wav         |
| SampleRate     | int    | 采样率，默认 16000             |
| Volume         | int    | 音量，范围为 0-100，默认 50    |
| SpeechRate     | int    | 语速，范围为-500-500，默认为 0 |
| PitchRate      | int    | 音高，范围为-500-500，默认为 0 |
| EnableSubtitle | bool   | 字幕功能，默认为 false         |

### 2. func DefaultSpeechSynthesisParam() SpeechSynthesisStartParam

> 创建一个默认的语音合成参数

参数说明：

无

返回值：

SpeechSynthesisStartParam：语音合成参数

### 3. func NewSpeechSynthesis(...) (\*SpeechSynthesis, error)

> 创建一个新的语音合成对象

参数说明：

| 参数            | 类型                      | 参数说明                                              |
| --------------- | ------------------------- | ----------------------------------------------------- |
| config          | \*ConnectionConfig        | 见上文建立连接相关内容                                |
| logger          | \*NlsLogger               | 见 SDK 日志相关内容                                   |
| realtime        | bool                      | 是否使用实时长文本，默认为短文本                      |
| taskfailed      | func(string, interface{}) | 识别过程中的错误处理回调，interface{}为用户自定义参数 |
| synthesisresult | func([]byte, interface{}) | 语音合成数据回调                                      |
| metainfo        | func(string, interface{}) | 字幕数据回调，需要参数中 EnableSubtitle 为 true       |
| completed       | func(string, interface{}) | 合成完毕结果回调                                      |
| closed          | func(interface{})         | 连接断开回调                                          |
| param           | interface{}               | 用户自定义参数                                        |

返回值：

无

### 4. func (tts \*SpeechSynthesis) Start(text string, param SpeechSynthesisStartParam, extra map[string]interface{}) (chan bool, error)

> 给定文本和参数进行语音合成

参数说明：

| 参数  | 类型                          | 参数说明            |
| ----- | ----------------------------- | ------------------- |
| text  | string                        | 待合成文本          |
| param | SpeechTranscriptionStartParam | 语音合成参数        |
| extra | map[string]interface{}        | 额外 key value 参数 |

返回值：

chan bool：语音合成完成通知管道

error：错误异常

### 5. func (tts \*SpeechSynthesis) Shutdown()

> 强制停止语音合成

参数说明：

无

返回值：

无

### 代码示例：

```python
package main

import (
        "errors"
        "flag"
        "fmt"
        "io"
        "log"
        "os"
        "os/signal"
        "sync"
        "time"

        "github.com/ikmak/alibabacloud-nls-go-sdk"
)

const (
  		AKID  = "Your AKID"
        AKKEY = "Your AKKEY"
        //online key
        APPKEY = "Your APPKEY"
        TOKEN  = "Your TOKEN"
)

type TtsUserParam struct {
        F           io.Writer
        Logger      *nls.NlsLogger
}

func onTaskFailed(text string, param interface{}) {
        p, ok := param.(*TtsUserParam)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        p.Logger.Println("TaskFailed:", text)
}

func onSynthesisResult(data []byte, param interface{}) {
        p, ok := param.(*TtsUserParam)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }
        p.F.Write(data)
}

func onCompleted(text string, param interface{}) {
        p, ok := param.(*TtsUserParam)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        p.Logger.Println("onCompleted:", text)
}


func onClose(param interface{}) {
        p, ok := param.(*TtsUserParam)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        p.Logger.Println("onClosed:")
}

func waitReady(ch chan bool, logger *nls.NlsLogger) error {
        select {
        case done := <-ch:
                {
                        if !done {
                                logger.Println("Wait failed")
                                return errors.New("wait failed")
                        }
                        logger.Println("Wait done")
                }
        case <-time.After(60 * time.Second):
                {
                        logger.Println("Wait timeout")
                        return errors.New("wait timeout")
                }
        }
        return nil
}

var lk sync.Mutex
var fail = 0
var reqNum = 0

const (
        TEXT = "你好小德，今天天气怎么样。"
)

func testMultiInstance(num int) {
        param := nls.DefaultSpeechSynthesisParam()
    	  config,_ := nls.NewConnectionConfigWithAKInfoDefault(nls.DEFAULT_URL, APPKEY, AKID, AKKEY)
        var wg sync.WaitGroup
        for i := 0; i < num; i++ {
                wg.Add(1)
                go func(id int) {
                        defer wg.Done()
                        strId := fmt.Sprintf("ID%d   ", id)
                        fname := fmt.Sprintf("ttsdump%d.wav", id)
                        ttsUserParam := new(TtsUserParam)
                        fout, err := os.OpenFile(fname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
                        logger := nls.NewNlsLogger(os.Stderr, strId, log.LstdFlags|log.Lmicroseconds)
                        logger.SetLogSil(false)
                        logger.SetDebug(true)
                        logger.Printf("Test Normal Case for SpeechRecognition:%s", strId)
                        ttsUserParam.F = fout
                        ttsUserParam.Logger = logger
                        tts, err := nls.NewSpeechSynthesis(config, logger, false,
                                onTaskFailed, onSynthesisResult, nil,
                                onCompleted, onClose, ttsUserParam)
                        if err != nil {
                                logger.Fatalln(err)
                                return
                        }

                        for {
                                lk.Lock()
                                reqNum++
                                lk.Unlock()
                                logger.Println("SR start")
                                ch, err := tts.Start(TEXT, param, nil)
                                if err != nil {
                                        lk.Lock()
                                        fail++
                                        lk.Unlock()
                                        tts.Shutdown()
                                        continue
                                }

                                err = waitReady(ch, logger)
                                if err != nil {
                                        lk.Lock()
                                        fail++
                                        lk.Unlock()
                                        tts.Shutdown()
                                        continue
                                }
                                logger.Println("Synthesis done")
                                tts.Shutdown()
                        }
                }(i)
        }

        wg.Wait()
}

func main() {
        coroutineId := flag.Int("num", 1, "coroutine number")
        flag.Parse()
        log.Default().Printf("start %d coroutines", *coroutineId)

        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        go func() {
                for range c {
                        lk.Lock()
                        log.Printf(">>>>>>>>REQ NUM: %d>>>>>>>>>FAIL: %d", reqNum, fail)
                        lk.Unlock()
                        os.Exit(0)
                }
        }()
        testMultiInstance(*coroutineId)
}

```
