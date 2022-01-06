# NLS Go SDK说明

> 本文介绍如何使用阿里云智能语音服务提供的Go SDK，包括SDK的安装方法及SDK代码示例。



## 前提条件

使用SDK前，请先阅读接口说明，详细请参见**接口说明**。

### 下载安装

> 说明
>
> * SDK支持go1.16
> * 请确认已经安装golang环境，并完成基本配置

1. 下载SDK

通过以下命令完成SDK下载和安装：

> go get github.com/aliyun/alibabacloud-nls-go-sdk

2. 导入SDK

在代码中通过将以下字段加入import来导入SDK：

> import ("github.com/aliyun/alibabacloud-nls-go-sdk")



## SDK常量

| 常量               | 常量含义                                                     |
| ------------------ | ------------------------------------------------------------ |
| SDK_VERSION        | SDK版本                                                      |
| PCM                | pcm音频格式                                                  |
| WAV                | wav音频格式                                                  |
| OPUS               | opus音频格式                                                 |
| OPU                | opu音频格式                                                  |
| DEFAULT_DISTRIBUTE | 获取token时使用的默认区域，"cn-shanghai"                     |
| DEFAULT_DOMAIN     | 获取token时使用的默认URL，"nls-meta.cn-shanghai.aliyuncs.com" |
| DEFAULT_VERSION    | 获取token时使用的协议版本，"2019-02-28"                      |
| DEFAULT_URL        | 默认公有云URL，"wss://nls-gateway.cn-shanghai.aliyuncs.com/ws/v1" |



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



## 获取token

### 1. func GetToken(dist string, domain string, akid string, akkey string, version string) (*TokenResultMessage, error)

> 获取访问token

参数说明：

| 参数    | 类型   | 参数说明                                    |
| ------- | ------ | ------------------------------------------- |
| dist    | string | 区域，如果不确定，请使用DEFAULT_DISTRIBUTE  |
| domain  | string | URL，如果不确定，请使用DEFAULT_DOMAIN       |
| akid    | string | 阿里云accessid                              |
| akkey   | string | 阿里云accesskey                             |
| version | string | 协议版本，如果不确定，请使用DEFAULT_VERSION |

返回值：

TokenResultMessage对象指针和错误信息



## 建立连接

### 1. ConnectionConfig

> 用于建立连接的基础参数

参数说明：

| 参数   | 类型   | 参数说明                                         |
| ------ | ------ | ------------------------------------------------ |
| Url    | string | 访问的公有云URL，如果不确定，可以使用DEFAULT_URL |
| Token  | string | 通过GetToken获取的token或者测试token             |
| Akid   | string | 阿里云accessid                                   |
| Akkey  | string | 阿里云accesskey                                  |
| Appkey | string | appkey，可以在控制台中对应项目上看到             |



### 2. func NewConnectionConfigWithAKInfoDefault(url string, appkey string, akid string, akkey string) (*ConnectionConfig, error) 

> 通过url，appkey，akid和akkey创建连接参数，等效于先调用GetToken然后再调用NewConnectionConfigWithToken

参数说明：

| 参数   | 类型   | 参数说明                                         |
| ------ | ------ | ------------------------------------------------ |
| Url    | string | 访问的公有云URL，如果不确定，可以使用DEFAULT_URL |
| Appkey | string | appkey，可以在控制台中对应项目上看到             |
| Akid   | string | 阿里云accessid                                   |
| Akkey  | string | 阿里云accesskey                                  |

返回值：

*ConnectionConfig：连接参数对象指针，用于后续创建语音交互实例

error：异常对象，为nil则无异常



### 3. func NewConnectionConfigWithToken(url string, appkey string, token string) *ConnectionConfig 

> 通过url，appkey和token创建连接参数

参数说明：

| 参数   | 类型   | 参数说明                                         |
| ------ | ------ | ------------------------------------------------ |
| Url    | string | 访问的公有云URL，如果不确定，可以使用DEFAULT_URL |
| Appkey | string | appkey，可以在控制台中对应项目上看到             |
| Token  | string | 已经通过GetToken或其他方式获取的token            |

返回值：

*ConnectionConfig：连接参数对象指针



### 4. func NewConnectionConfigFromJson(jsonStr string) (*ConnectionConfig, error) 

> 通过json字符串来创建连接参数

参数说明

| 参数    | 类型   | 参数说明                                                     |
| ------- | ------ | ------------------------------------------------------------ |
| jsonStr | string | 描述连接参数的json字符串，有效字段如下：url，token，akid，akkey，appkey。其中必须包含url和appkey，如果包含token则不需要包含akid和akkey |

返回值：

*ConnectionConfig：连接对象指针

## 一句话识别

### 1. SpeechRecognitionStartParam

> 一句话识别参数

参数说明:

| 参数                           | 类型   | 参数说明              |
| ------------------------------ | ------ | --------------------- |
| Format                         | string | 音频格式，默认使用pcm |
| SampleRate                     | int    | 采样率，默认16000     |
| EnableIntermediateResult       | bool   | 是否打开中间结果返回  |
| EnablePunctuationPredition     | bool   | 是否打开标点预测      |
| EnableInverseTextNormalization | bool   | 是否打开ITN           |



### 2. func DefaultSpeechRecognitionParam() SpeechRecognitionStartParam

> 返回一个默认的推荐参数，其中format为pcm，采样率为16000，中间结果，标点预测和ITN全开

参数说明：

无

返回值：

默认参数



### 3. func NewSpeechRecognition(...) (*SpeechRecognition, error)

> 创建一个SpeechRecognition实例

参数说明：

| 参数          | 类型                      | 参数说明                                              |
| ------------- | ------------------------- | ----------------------------------------------------- |
| config        | *ConnectionConfig         | 见上文建立连接相关内容                                |
| logger        | *NlsLogger                | 见SDK日志相关内容                                     |
| taskfailed    | func(string, interface{}) | 识别过程中的错误处理回调，interface{}为用户自定义参数 |
| started       | func(string, interface{}) | 建连完成回调                                          |
| resultchanged | func(string, interface{}) | 识别中间结果回调                                      |
| completed     | func(string, interface{}) | 最终识别结果回调                                      |
| closed        | func(interface{})         | 连接断开回调                                          |
| param         | interface{}               | 用户自定义参数                                        |

返回值：

*SpeechRecognition：识别对象指针

error：错误异常



### 4. func (sr *SpeechRecognition) Start(param SpeechRecognitionStartParam, extra map[string]interface{}) (chan bool, error)

> 根据param发起一次一句话识别

参数说明：

| 参数  | 类型                        | 参数说明          |
| ----- | --------------------------- | ----------------- |
| param | SpeechRecognitionStartParam | 一句话识别参数    |
| extra | map[string]interface{}      | 额外key value参数 |

返回值：

chan bool：同步start完成的管道

error：错误异常

### 5. func (sr *SpeechRecognition) Stop() (chan bool, error)

> 停止一句话识别

参数说明：

无

返回值：

chan bool：同步stop完成的管道

error：错误异常



### 6. func (sr *SpeechRecognition) Shutdown()

> 强制断开连接

参数说明：

无

返回值：

无



### 7. func (sr *SpeechRecognition) SendAudioData(data []byte) error

> 发送音频，音频格式必须和参数中一致

参数说明

| 参数 | 类型   | 参数说明 |
| ---- | ------ | -------- |
| data | []byte | 音频数据 |

返回值：

error：异常错误



### 一句话识别代码示例：

```python
package main

import (
        "errors"
        "flag"
        "fmt"
        "log"
        "os"
        "os/signal"
        "sync"
        "time"

        "github.com/aliyun/alibabacloud-nls-go-sdk"
)

const (
  		AKID  = "Your AKID"
        AKKEY = "Your AKKEY"
        //online key
        APPKEY = "Your APPKEY"
        TOKEN  = "Your TOKEN"
)

func onTaskFailed(text string, param interface{}) {
        logger, ok := param.(*nls.NlsLogger)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        logger.Println("TaskFailed:", text)
}

func onStarted(text string, param interface{}) {
        logger, ok := param.(*nls.NlsLogger)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        logger.Println("onStarted:", text)
}

func onResultChanged(text string, param interface{}) {
        logger, ok := param.(*nls.NlsLogger)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        logger.Println("onResultChanged:", text)
}

func onCompleted(text string, param interface{}) {
        logger, ok := param.(*nls.NlsLogger)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        logger.Println("onCompleted:", text)
}

func onClose(param interface{}) {
        logger, ok := param.(*nls.NlsLogger)
        if !ok {
                log.Default().Fatal("invalid logger")
                return
        }

        logger.Println("onClosed:")
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
        case <-time.After(20 * time.Second):
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

func testMultiInstance(num int) {
        pcm, err := os.Open("tests/test1.pcm")
        if err != nil {
                log.Default().Fatalln(err)
        }

        buffers := nls.LoadPcmInChunk(pcm, 320)
        param := nls.DefaultSpeechRecognitionParam()
    	config := nls.NewConnectionConfigWithAKInfoDefault(nls.DEFAULT_URL, APPKEY, AKID, AKKEY)
        var wg sync.WaitGroup
        for i := 0; i < num; i++ {
                wg.Add(1)
                go func(id int) {
                        defer wg.Done()
                        strId := fmt.Sprintf("ID%d   ", id)
                        logger := nls.NewNlsLogger(os.Stderr, strId, 			log.LstdFlags|log.Lmicroseconds)
                        logger.SetLogSil(false)
                        logger.SetDebug(true)
      logger.Printf("Test Normal Case for SpeechRecognition:%s", strId)
                        sr, err := nls.NewSpeechRecognition(config, logger,
                                onTaskFailed, onStarted, onResultChanged,
                                onCompleted, onClose, logger)
                        if err != nil {
                                logger.Fatalln(err)
                                return
                        }

      test_ex := make(map[string]interface{})
      test_ex["test"] = "hello"

                        for {
                                lk.Lock()
                                reqNum++
                                lk.Unlock()
                                logger.Println("SR start")
                                ready, err := sr.Start(param, test_ex)
                                if err != nil {
                                        lk.Lock()
                                        fail++
                                        lk.Unlock()
                                        sr.Shutdown()
                                        continue
                                }

                                err = waitReady(ready, logger)
                                if err != nil {
                                        lk.Lock()
                                        fail++
                                        lk.Unlock()
                                        sr.Shutdown()
                                        continue
                                }

                                for _, data := range buffers.Data {
                                        if data != nil {
                                                sr.SendAudioData(data.Data)
                                                time.Sleep(10 * time.Millisecond)
                                        }
                                }

                                logger.Println("send audio done")
                                ready, err = sr.Stop()
                                if err != nil {
                                        lk.Lock()
                                        fail++
                                        lk.Unlock()
                                        sr.Shutdown()
                                        continue
                                }

                                err = waitReady(ready, logger)
                                if err != nil {
                                        lk.Lock()
                                        fail++
                                        lk.Unlock()
                                        sr.Shutdown()
                                        continue
                                }

                                logger.Println("Sr done")
                                sr.Shutdown()
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





