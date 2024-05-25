package task

import (
	"fmt"
	"testing"
	"time"

	"github.com/hujinrun-github/tool_functions/go/log"
)

type TaskTest struct {
}

func (t *TaskTest) Execute(TaskInfoST) (TaskResponseST, error) {
	time.Sleep(time.Second * 5)
	return TaskResponseST{ExtraRet: map[string][]byte{
		"test": []byte("dfadf"),
	}}, nil
}

func (t *TaskTest) ResHandle(tr TaskResponseST) error {
	fmt.Printf("res:%+v\n", tr)
	return nil
}

func TestTask(t *testing.T) {
	tt := &TaskTest{}
	RegisterTaskExec(TASKTYPE_T(1), tt)

	// address := "127.0.0.1:8999"
	s := NewTaskServer(WithLogLevel(log.LEVEL_DEBUG))
	s.Run()

	c := NewTaskClientST("127.0.0.1:9008", WithClientDebugLevel(log.LEVEL_DEBUG))
	c.Run()
	c.AddTask(&TaskInfoST{
		TaskType: 1,
		TaskID:   "112",
	})

	c.AddTask(&TaskInfoST{
		TaskType: 1,
		TaskID:   "113",
	})

	c.AddTask(&TaskInfoST{
		TaskType: 1,
		TaskID:   "114",
	})
	time.Sleep(time.Second * 10)
}
