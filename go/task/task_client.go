package task

import (
	"context"
	"fmt"

	"github.com/hujinrun-github/tool_functions/go/concurrent"
	"github.com/hujinrun-github/tool_functions/go/log"
	pb "github.com/hujinrun-github/tool_functions/go/task/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// task client负责分配任务，并等待结果

type TaskClientST struct {
	bufferSize    int
	serverAddress string
	debugLevel    log.LogLevel

	ch     chan *TaskInfoST
	logger *log.Logger
}

func NewTaskClientST(address string, opts ...ClientOptions) *TaskClientST {
	tc := &TaskClientST{
		serverAddress: address,
		bufferSize:    DEFAULT_BUFFER_SIZE,
	}

	for _, opt := range opts {
		opt(tc)
	}

	tc.logger = log.NewLogHandler(log.WithLevel(tc.debugLevel))
	tc.ch = make(chan *TaskInfoST, tc.bufferSize)

	return tc
}

func (t *TaskClientST) Run() {
	go concurrent.Try(func() {
		for info := range t.ch {
			t.logger.Debug("task client received task", fmt.Sprintf("%+v", info))
			concurrent.Try(func() {
				t.logger.Debug("task client received task", fmt.Sprintf("%+v", info))
				executer := GetTaskExec(info.TaskType)
				if executer == nil {
					t.logger.Errorf("no task executer for type: %d", info.TaskType)
					return
				}

				conn, err := grpc.NewClient(t.serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					t.logger.Errorf("failed to connect server: %s", err.Error())
					return
				}
				defer conn.Close()
				client := pb.NewTaskServiceClient(conn)
				res, err := client.SubmitTask(context.Background(), &pb.TaskRequest{
					TaskId:        string(info.TaskID),
					TaskType:      uint32(info.TaskType),
					TaskExtraInfo: info.ExtraInfo,
					WithResponse:  info.WithResponse,
				})

				if err != nil {
					t.logger.Errorf("failed to submit task: %s", err.Error())
					return
				}

				err = executer.ResHandle(TaskResponseST{
					TaskID:   TASKID_T(res.TaskId),
					RetCode:  TASKRET_T(res.RetCode),
					ErrMsg:   res.ErrMsg,
					ExtraRet: res.ExtraRet,
				})
				if err != nil {
					t.logger.Errorf("failed to handle task response: %s", err.Error())
					return
				}
			})
		}
	})
}

func (t *TaskClientST) AddTask(task *TaskInfoST) {
	// t.logger.Debug("task client received task", fmt.Sprintf("%+v", *task))
	t.ch <- task
}

// stop前需要保证，所有的AddTask都执行完毕
func (t *TaskClientST) Stop() {
	close(t.ch)
}

// =====================
type ClientOptions func(client *TaskClientST)

func WithServerAddress(serverAddress string) ClientOptions {
	return func(client *TaskClientST) {
		client.serverAddress = serverAddress
	}
}

func WithClientDebugLevel(debugLevel log.LogLevel) ClientOptions {
	return func(client *TaskClientST) {
		client.debugLevel = debugLevel
	}
}
