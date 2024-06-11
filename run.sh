#!/bin/bash

export CGO_LDFLAGS="-lcublas -lcudart -L/usr/local/cuda/lib64/"
export LIBRARY_PATH=$PWD
export C_INCLUDE_PATH=$PWD

go run . -m "/home/ctyler/llm_models/mistral-7b-instruct-v0.2.Q3_K_S.gguf"

