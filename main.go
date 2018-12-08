// Copyright 2018 Oliver Szabo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/oleewere/go-buffered-processor/processor"
	"sync"
	"time"
)

// ExampleDataProcessor example implementation of processor interface
type ExampleDataProcessor struct {
	Mutex *sync.Mutex
}

// Process example implementation of process function for Processor interface
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

// HandleError example implementation of error handling
func (p ExampleDataProcessor) HandleError(batchContext *processor.BatchContext, err error) {
	fmt.Println(err)
}

func main() {
	proc := ExampleDataProcessor{Mutex: &sync.Mutex{}}
	b := processor.CreateDefaultBatchContext()
	b.ProcessTimeInterval = 2 * time.Second
	data := "logging something"
	go processor.StartTimeBasedProcessing(b, proc, 1)
	go processor.StartTimeBasedProcessing(b, proc, 2)
	go processor.StartTimeBasedProcessing(b, proc, 3)
	processor.ProcessData(data, b, proc)
	processor.ProcessData(data, b, proc)
	processor.ProcessData(data, b, proc)
	processor.ProcessData(data, b, proc)
	processor.StartTimeBasedProcessing(b, proc, 10)
}
