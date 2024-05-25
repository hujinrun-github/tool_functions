package task

type TaskInfoST struct {
	TaskID       TASKID_T          // 任务ID
	TaskType     TASKTYPE_T        // 任务类型
	ExtraInfo    map[string][]byte // 额外信息
	WithResponse bool              // 是否需要响应
}

type TaskResponseST struct {
	TaskID   TASKID_T
	RetCode  TASKRET_T
	ErrMsg   string
	ExtraRet map[string][]byte
}
