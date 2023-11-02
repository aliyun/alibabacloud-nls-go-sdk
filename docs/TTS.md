# NLS Go Dashscope SDK说明

> 本文介绍如何使用阿里云智能语音服务为DashScope平台提供的Go SDK，包括SDK的安装方法及SDK代码示例。



## 前提条件

使用SDK前，请先阅读接口说明，详细请参见**接口说明**。

### 下载安装

> 说明
>
> * SDK支持go1.16
> * 请确认已经安装golang环境，并完成基本配置

1. 下载SDK

通过以下命令获取SDK

> go get -u github.com/aliyun/alibabacloud-nls-go-sdk@v0.1.0


2. 导入SDK

在代码中通过将以下字段加入import来导入SDK：

> import ("github.com/aliyun/alibabacloud-nls-go-sdk")



## SDK常量

| 常量               | 常量含义                                                     |
| ------------------ | ------------------------------------------------------------ |
| SDK_VERSION        | SDK版本                                                      |
| PCM                | pcm音频格式                                                  |
| WAV                | wav音频格式                                                  |
| DEFAULT_URL        | 默认公有云URL，"wss://dashscope.aliyuncs.com/api-ws/v1/inference"" |
| DEFAULT_WS_RBUFFER_SIZE | 默认Websocket读Buffer长度，4096 |
| DEFAULT_WS_WBUFFER_SIZE | 默认Websocket写Buffer长度，4096 |



## SDK日志

### 1. func DefaultNlsLog() *NlsLogger

> 用于创建全局唯一的默认日志对象，默认日志以NLS为前缀，输出到标准错误

参数说明：

无

返回值：

NlsLogger对象指针



### 2. func NewNlsLogger(w io.Writer, tag string, flag int) *NlsLogger 

> 创建一个新的日志

参数说明：

| 参数 | 类型      | 参数说明                        |
| ---- | --------- | ------------------------------- |
| w    | io.Writer | 任意实现io.Writer接口的对象     |
| tag  | string    | 日志前缀，会打印到日志行首部    |
| flag | int       | 日志flag，具体参考go官方log文档 |

返回值：

NlsLogger对象指针



### 3. func (logger *NlsLogger) SetLogSil(sil bool) 

> 设置日志是否输出到对应的io.Writer

参数说明:

| 参数 | 类型 | 参数说明                     |
| ---- | ---- | ---------------------------- |
| sil  | bool | 是否禁止日志输出，true为禁止 |

返回值：

无



### 4. func (logger *NlsLogger) SetDebug(debug bool)

> 设置是否打印debug日志，仅影响通过Debugf或Debugln进行输出的日志

参数说明：

| 参数  | 类型 | 参数说明                          |
| ----- | ---- | --------------------------------- |
| debug | bool | 是否允许debug日志输出，true为允许 |

返回值：

无



### 5. func (logger *NlsLogger) SetOutput(w io.Writer)

> 设置日志输出方式

参数说明：

| 参数 | 类型      | 参数说明                    |
| ---- | --------- | --------------------------- |
| w    | io.Writer | 任意实现io.Writer接口的对象 |

返回值：

无



### 6. func (logger *NlsLogger) SetPrefix(prefix string)

> 设置日志行的标签

参数说明：

| 参数   | 类型   | 参数说明                       |
| ------ | ------ | ------------------------------ |
| prefix | string | 日志行标签，会输出在日志行行首 |

返回值：

无



### 7. func (logger *NlsLogger) SetFlags(flags int)

> 设置日志属性

参数说明：

| 参数  | 类型 | 参数说明                                         |
| ----- | ---- | ------------------------------------------------ |
| flags | int  | 日志属性，见https://pkg.go.dev/log#pkg-constants |

返回值：

无



### 8. 日志打印

日志打印方法：

| 方法名                                                      | 方法说明                                                     |
| ----------------------------------------------------------- | ------------------------------------------------------------ |
| func (l *NlsLogger) Print(v ...interface{})                 | 标准日志输出                                                 |
| func (l *NlsLogger) Println(v ...interface{})               | 标注日志输出，行尾自动换行                                   |
| func (l *NlsLogger) Printf(format string, v ...interface{}) | 带format的日志输出，format方式见go官方文档                   |
| func (l *NlsLogger) Debugln(v ...interface{})               | debug信息日志输出，行尾自动换行                              |
| func (l *NlsLogger) Debugf(format string, v ...interface{}) | 带format的debug信息日志输出                                  |
| func (l *NlsLogger) Fatal(v ...interface{})                 | 致命错误日志输出，输出后自动进程退出                         |
| func (l *NlsLogger) Fatalln(v ...interface{})               | 致命错误日志输出，行尾自动换行，输出后自动进程退出           |
| func (l *NlsLogger) Fatalf(format string, v ...interface{}) | 带format的致命错误日志输出，输出后自动进程退出               |
| func (l *NlsLogger) Panic(v ...interface{})                 | 致命错误日志输出，输出后自动进程退出并打印崩溃信息           |
| func (l *NlsLogger) Panicln(v ...interface{})               | 致命错误日志输出，行尾自动换行，输出后自动进程退出并打印崩溃信息 |
| func (l *NlsLogger) Panicf(format string, v ...interface{}) | 带format的致命错误日志输出，输出后自动进程退出并打印崩溃信息 |



## 建立连接

### 1. NewConnectionConfigDefault

> 使用默认参数创建连接配置，该接口会通过”DASHSCOPE_API_KEY“环境变量来获取APIKEY，如果没有设置，该接口会报错，同时该接口使用默认URL，和1024读Buffer及4096写Buffer长度

参数说明：

无



### 2. func NewConnectionConfigWithUrlApiKey(url string, apikey string, rbuffer int, wbuffer int) (*ConnectionConfig, error) 

> 通过url，apikey，rbuffer和wbuffer创建连接参数

参数说明：

| 参数   | 类型   | 参数说明                                         |
| ------ | ------ | ------------------------------------------------ |
| url    | string | 访问的DashScope URL，如果不确定，可以使用DEFAULT_URL |
| apikey | string | apikey，可以在控制台中对应项目上看到，不建议在使用中明文暴露，会有泄露风险    |
| rbuffer   | int | Websocket读缓冲长度，默认为4096，大多数情况可以不修改               |
| wbuffer  | int |  Websocket写缓冲长度，默认为4096，大多数情况可以不修改               |

返回值：

*ConnectionConfig：连接参数对象指针，用于后续创建语音交互实例

error：异常对象，为nil则无异常


## 语音合成

### 1. SpeechSynthesisStartParam


参数说明:

| 参数           | 类型   | 参数说明                      |
| -------------- | ------ | ----------------------------- |
| TextType       | string | 不需要修改，为“PlainText”     |
| Format         | string | 音频格式，默认使用wav         |
| SampleRate     | int    | 采样率，默认16000             |
| Volume         | int    | 音量，范围为0-100，默认50     |
| Rate            | float32    | 语速，范围为0.5-2.0，默认为1.0 |
| Pitch         | float32    | 音高，范围为0.5-2.0，默认为1.0 |
| EnableWordTimestamp | bool   | 字幕功能，默认为false         |
| EnablePhonemeTimestamp | bool   | 音素字幕功能，默认为false，开启同时需要先开启字幕功能  |

### 2.  func DefaultSpeechSynthesisParam() SpeechSynthesisStartParam

> 创建一个默认的语音合成参数

参数说明：

无

返回值：

SpeechSynthesisStartParam：语音合成参数

### 3. func NewSpeechSynthesis(...) (*SpeechSynthesis, error)

> 创建一个新的语音合成对象

参数说明：

| 参数            | 类型                      | 参数说明                                              |
| --------------- | ------------------------- | ----------------------------------------------------- |
| config          | *ConnectionConfig         | 见上文建立连接相关内容                                |
| logger          | *NlsLogger                | 见SDK日志相关内容                                     |
| started        | func(string, interface{})  | 连接建立回调，第一个参数为taskid，第二个参数为用户自定义参数  |
| taskfailed      | func(string, interface{}) | 识别过程中的错误处理回调，interface{}为用户自定义参数 |
| synthesisresult | func([]byte, interface{}) | 语音合成数据回调                                      |
| metainfo        | func(string, interface{}) | 字幕数据回调，需要参数中EnableSubtitle为true          |
| completed       | func(string, interface{}) | 合成完毕结果回调                                      |
| closed          | func(interface{})         | 连接断开回调                                          |
| param           | interface{}               | 用户自定义参数                                        |

返回值：

无

### 4. func (tts *SpeechSynthesis) Start(model string, text string, param SpeechSynthesisStartParam, extra map[string]interface{}) (chan bool, error) 

> 给定文本和参数进行语音合成

参数说明：

| 参数  | 类型                          | 参数说明          |
| ----- | ----------------------------- | ----------------- |
| model | string                        | 调用模型名        |
| text  | string                        | 待合成文本        |
| param | SpeechTranscriptionStartParam | 语音合成参数      |
| extra | map[string]interface{}        | 额外key value参数 |

返回值：

chan bool：语音合成完成通知管道

error：错误异常

### 5. func (tts *SpeechSynthesis) Shutdown()

> 强制停止语音合成

参数说明：

无

返回值：

无



### 代码示例：

下面代码通过num传入并行实例数，通过model传入模型名称，可以参考run_tts_test.sh脚本。
```go
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"time"

	dash "github.com/aliyun/alibabacloud-nls-go-sdk"
)

type TtsUserParam struct {
	F      io.Writer
	Logger *dash.NlsLogger
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

func onTaskStarted(taskid string, param interface{}) {
	p, of := param.(*TtsUserParam)
	if !of {
		log.Default().Fatal("invalid logger")
		return
	}
	p.Logger.Println("onTaskStarted:", taskid)
}

func waitReady(ch chan bool, logger *dash.NlsLogger) error {
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

func testMultiInstance(num int, model string) {
	param := dash.DefaultSpeechSynthesisParam()
	param.EnableWordTimestamp = true
	param.EnablePhonemeTimestamp = true
	config, e := dash.NewConnectionConfigDefault()
	if e != nil {
		log.Fatal(e)
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			strId := fmt.Sprintf("ID%d   ", id)
			fname := fmt.Sprintf("ttsdump%d.wav", id)
			ttsUserParam := new(TtsUserParam)
			fout, err := os.OpenFile(fname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
			logger := dash.NewNlsLogger(os.Stderr, strId, log.LstdFlags|log.Lmicroseconds)
			logger.SetLogSil(false)
			logger.SetDebug(true)
			logger.Printf("Test Normal Case for SpeechRecognition:%s", strId)
			ttsUserParam.F = fout
			ttsUserParam.Logger = logger
			//third param control using realtime long text tts
			tts, err := dash.NewSpeechSynthesis(config, logger, onTaskStarted,
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
				logger.Printf("TTS start: model=%s", model)
				ch, err := tts.Start(model, TEXT, param, nil)
				if err != nil {
					lk.Lock()
					fail++
					lk.Unlock()
					tts.Shutdown()
					time.Sleep(time.Second * 2)
					continue
				}

				err = waitReady(ch, logger)
				if err != nil {
					lk.Lock()
					fail++
					lk.Unlock()
					tts.Shutdown()
					time.Sleep(time.Second * 2)
					continue
				}
				logger.Println("Synthesis done")
				tts.Shutdown()
				break
			}
		}(i)
	}

	wg.Wait()
}

func main() {
	go func() {
		log.Default().Println(http.ListenAndServe(":6060", nil))
	}()
	coroutineId := flag.Int("num", 1, "coroutine number")
	modelId := flag.String("model", "sambert-zhimao-v1", "model id")
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
	testMultiInstance(*coroutineId, *modelId)
}


```



