package task

import (
	"fmt"
	"sync"
)

type ITaskExec interface {
	Execute(TaskInfoST) (TaskResponseST, error)
	ResHandle(TaskResponseST) error
}

var Tasks = make(map[TASKTYPE_T]ITaskExec)
var lock sync.RWMutex

func RegisterTaskExec(taskType TASKTYPE_T, taskExec ITaskExec) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := Tasks[taskType]; !ok {
		Tasks[taskType] = taskExec
	} else {
		return fmt.Errorf("[ERROR]task type %v already exists", taskType)
	}
	return nil
}

func GetTaskExec(taskType TASKTYPE_T) ITaskExec {
	lock.RLock()
	defer lock.RUnlock()
	if task, ok := Tasks[taskType]; ok {
		return task
	} else {
		return nil
	}
}
