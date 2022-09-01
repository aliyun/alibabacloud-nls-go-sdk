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

	"github.com/aliyun/alibabacloud-nls-go-sdk"
)

const (
	AKID  = "Your AKID"
	AKKEY = "Your AKKEY"
	//online key
	APPKEY = "Your APPKEY"
	TOKEN  = "TEST TOKEN"
)

type TtsUserParam struct {
	F      io.Writer
	Logger *nls.NlsLogger
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
	config := nls.NewConnectionConfigWithToken(nls.DEFAULT_URL,
		APPKEY, TOKEN)
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
			//third param control using realtime long text tts
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
	go func() {
		log.Default().Println(http.ListenAndServe(":6060", nil))
	}()
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
