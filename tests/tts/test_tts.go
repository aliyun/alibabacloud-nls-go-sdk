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

	dash "github.com/aliyun/alibabacloud-dash-go-sdk"
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
