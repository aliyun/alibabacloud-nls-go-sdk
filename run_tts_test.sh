#!/bin/bash
APIKEY=$1
MODEL=$2
export DASHSCOPE_API_KEY=$APIKEY

go run tests/tts/test_tts.go -num 1 -model $MODEL
