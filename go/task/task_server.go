package task

import (
	"context"
	"fmt"
	"net"

	"github.com/hujinrun-github/tool_functions/go/concurrent"
	"github.com/hujinrun-github/tool_functions/go/log"
	pb "github.com/hujinrun-github/tool_functions/go/task/proto"
	"github.com/panjf2000/ants/v2"
	"google.golang.org/grpc"
)

// task server负责处理任务，并返回结果

type TaskServerST struct {
	pb.UnimplementedTaskServiceServer
	poolSize int
	logLevel log.LogLevel
	port     int

	pool   *ants.Pool
	logger *log.Logger
	server *grpc.Server
}

// 需要判断TaskServerST是否为空
func NewTaskServer(options ...ServerOptions) *TaskServerST {
	server := &TaskServerST{
		poolSize: DEFAULT_POOL_SIZE,
		port:     DEFAULT_SEVER_PORT,
	}
	for _, opt := range options {
		opt(server)
	}

	server.logger = log.NewLogHandler(log.WithLevel(server.logLevel))

	var err error
	server.pool, err = ants.NewPool(server.poolSize, ants.WithNonblocking(true))
	if err != nil {
		server.logger.Errorf("create pool failed, err::%v", err)
		return nil
	}
	return server
}

// Run 函数用于启动TaskServerST的grpc服务
func (s *TaskServerST) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Errorf("listen failed, err::%v", err)
		return err
	}

	s.logger.Errorf("server start at:%+v", lis.Addr())

	s.server = grpc.NewServer()
	pb.RegisterTaskServiceServer(s.server, s)

	go concurrent.Try(func() {
		err = s.server.Serve(lis)
		if err != nil {
			s.logger.Errorf("start server failed, err::%v", err)
		}
	})

	return nil
}

func (s *TaskServerST) Stop() error {
	s.server.Stop()
	s.pool.Release()
	return nil
}

// SubmitTask 是TaskServerST结构体的一个方法，用于提交任务，只是内部使用，一般不做外部调用
//
// 参数：
// ctx：上下文对象，用于控制任务的执行环境
// req：指向pb.TaskRequest类型的指针，表示要提交的任务请求
//
// 返回值：
// *pb.TaskResponse：指向pb.TaskResponse类型的指针，表示任务执行的结果
// error：任务执行过程中可能产生的错误
func (s *TaskServerST) SubmitTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	// 获取执行器
	res := &pb.TaskResponse{}
	executor := GetTaskExec(TASKTYPE_T(req.TaskType))
	if executor == nil {
		res.ErrMsg = fmt.Sprintf("get task executor failed, type::%v", req.TaskType)
		res.RetCode = RET_CANNOT_GET_TASK_HANDLE
		s.logger.Errorf(res.ErrMsg)
		return res, fmt.Errorf(res.ErrMsg)
	}

	s.logger.Debug("get task executor success, taskInfo::", req.String())

	// 1.如果需要等待返回，则等待结果，否则直接返回
	if req.WithResponse {
		innnerRes, err := executor.Execute(TaskInfoST{
			TaskID:       TASKID_T(req.TaskId),
			TaskType:     TASKTYPE_T(req.TaskType),
			ExtraInfo:    req.TaskExtraInfo,
			WithResponse: req.WithResponse,
		})
		if err != nil {
			s.logger.Errorf("execute task failed, err::%v", err)
		}

		res.TaskId = string(innnerRes.TaskID)
		res.ExtraRet = innnerRes.ExtraRet
		res.ErrMsg = innnerRes.ErrMsg
		res.RetCode = int32(innnerRes.RetCode)
		return res, err
	} else {
		s.pool.Submit(func() {
			concurrent.Try(func() {
				_, err := executor.Execute(TaskInfoST{
					TaskID:       TASKID_T(req.TaskId),
					TaskType:     TASKTYPE_T(req.TaskType),
					ExtraInfo:    req.TaskExtraInfo,
					WithResponse: req.WithResponse,
				})
				if err != nil {
					s.logger.Errorf("execute task failed, err::%v", err)
				}
			})
		})
		res.TaskId = string(req.TaskId)
		res.ErrMsg = ""
		res.RetCode = int32(RET_OK)
		return res, nil
	}
}

// ================server options=================
type ServerOptions func(server *TaskServerST)

func WithPoolSize(poolSize int) ServerOptions {
	return func(server *TaskServerST) {
		server.poolSize = poolSize
	}
}

func WithLogLevel(level log.LogLevel) ServerOptions {
	return func(server *TaskServerST) {
		server.logLevel = level
	}
}

func WithPort(port int) ServerOptions {
	return func(server *TaskServerST) {
		server.port = port
	}
}
