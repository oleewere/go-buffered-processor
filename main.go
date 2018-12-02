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

type ExampleDataProcessor struct {
	BatchContext *processor.BatchContext
	Mutex sync.Mutex
}

func (p ExampleDataProcessor) Process() error {
	fmt.Println("Processing...")
	t := p.BatchContext.LastChanged
	fmt.Println(t.Format("2006-01-02-15:04:05"))
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	for _, data := range *p.BatchContext.BufferData {
		fmt.Println(data)
	}
	return nil
}

func (s ExampleDataProcessor) GetBatchContext() *processor.BatchContext {
	return s.BatchContext
}

func (s ExampleDataProcessor) HandleError(err error) {
	fmt.Println(err)
}

func main() {
	proc := ExampleDataProcessor{}
	b := processor.CreateDefaultBatchContext()
	b.ProcessTimeInterval = 2 * time.Second
	proc.BatchContext = b
	data := "logging something"
	go processor.StartTimeBasedProcessing(proc, 1)
	go processor.StartTimeBasedProcessing(proc, 2)
	go processor.StartTimeBasedProcessing(proc, 3)
	processor.ProcessData(data, proc)
	processor.ProcessData(data, proc)
	processor.ProcessData(data, proc)
	processor.ProcessData(data, proc)
	processor.StartTimeBasedProcessing(proc, 10)
}