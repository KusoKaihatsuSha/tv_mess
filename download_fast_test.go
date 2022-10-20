package main

import (
	"fmt"
	"testing"
)

func Test_tasker_downloading_fast(t *testing.T) {
	fmt.Println(t.Name())
	Tr := new(Tasker).Init(2, 512)
	Tr.Add([]any{9999, t}, downloading_tasker_tester, &Message{})
	defer Tr.Wg.Wait()
}

func Test_tasker_checker_fast(t *testing.T) {
	fmt.Println(t.Name())
	Tr := new(Tasker).Init(2, 512)
	mes := &Message{}
	Tr.Add([]any{9999, t}, hash_files_tester, mes)
	Tr.Wg.Wait()
	//os.RemoveAll(testid + sprintf("__%d", 9999))
}
