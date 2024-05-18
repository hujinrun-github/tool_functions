package task

type TaskInfoST struct {
	TaskId   string // 任务ID
	TaskType TASKTYPE // 任务类型
	ExtraInfo map[string][]byte{} // 额外信息
}

type TaskResultST struct {
	
}
