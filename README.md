<!---
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

## Simple Buffered Data Processor in Go

[![GoDoc Widget](https://godoc.org/github.com/oleewere/go-buffered-processor/processor?status.svg)](https://godoc.org/github.com/oleewere/go-buffered-processor/processor)
[![Build Status](https://travis-ci.org/oleewere/go-buffered-processor.svg?branch=master)](https://travis-ci.org/oleewere/go-buffered-processor)
[![Go Report Card](https://goreportcard.com/badge/github.com/oleewere/go-buffered-processor)](https://goreportcard.com/report/github.com/oleewere/go-buffered-processor)
![license](http://img.shields.io/badge/license-Apache%20v2-blue.svg)

### Description

Small framework to send data to a buffer, and process that data if the buffer is full or after some time (check the buffer with a task in the background).

### Install & Build
Requires go 1.10.x+
- Install the source
```bash
make install
```
- Build the source
```bash
make build
```

### Usage

```go
package main

import (
    "fmt"
    "sync"
    "time"
    "github.com/oleewere/go-buffered-processor/processor"
)

//  example implementation of processor interface
type ExampleDataProcessor struct {
    Mutex *sync.Mutex
}

// example implementation of process function for Processor interface
func (p ExampleDataProcessor) Process(batchContext *processor.BatchContext) error {
    fmt.Println("Processing...")
    t := batchContext.LastChanged
    fmt.Println(t.Format("2006-01-02-15:04:05"))
    p.Mutex.Lock()
    defer p.Mutex.Unlock()
    for _, data := range *batchContext.BufferData {
        fmt.Println(data)
    }
	return nil
}

// example implementation of error handling
func (p ExampleDataProcessor) HandleError(batchContext *processor.BatchContext, err error) {
    fmt.Println(err)
}

func main() { 
    // ... 
    proc := ExampleDataProcessor{Mutex: &sync.Mutex{}}
    b := processor.CreateDefaultBatchContext()
    b.ProcessTimeInterval = 2 * time.Second
    data := "logging something"
    // simply process some data
    processor.ProcessData(data, b, proc)
    // start a background task to process data if it reached the process time interval
    go processor.StartTimeBasedProcessing(b, proc, 1)
    // ...
}
```