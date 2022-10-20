package main

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

// if more than 8, may block on google server
func Test_tasker_downloading(t *testing.T) {
	fmt.Println(t.Name())
	Tr := new(Tasker).Init(runtime.NumCPU(), 512)
	for i := 1; i <= runtime.NumCPU(); i++ {
		Tr.Add([]any{i, t}, downloading_tasker_tester, &Message{})
		<-time.After(500 * time.Millisecond)
	}
	defer Tr.Wg.Wait()
}

func Test_tasker_checker(t *testing.T) {
	fmt.Println(t.Name())
	Tr := new(Tasker).Init(runtime.NumCPU(), 512)
	mes := &Message{}
	for i := 1; i <= runtime.NumCPU(); i++ {
		Tr.Add([]any{i, t}, hash_files_tester, mes)
		<-time.After(500 * time.Millisecond)
	}
	Tr.Wg.Wait()
	for i := 1; i <= runtime.NumCPU(); i++ {
		os.RemoveAll(testid + sprintf("__%d", i))
	}
}
