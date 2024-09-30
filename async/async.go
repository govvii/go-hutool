package async

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

// Task 表示一个异步任务
type Task func(ctx context.Context) (interface{}, error)

// Result 表示异步任务的结果
type Result struct {
	Value interface{}
	Err   error
}

// AsyncExecutor 处理异步任务执行
type AsyncExecutor struct {
	workerPool chan struct{}
	results    chan Result
	wg         sync.WaitGroup
	logger     *log.Logger
	ctx        context.Context
	cancel     context.CancelFunc
}

// ExecutorOption 是设置 AsyncExecutor 选项的函数类型
type ExecutorOption func(*AsyncExecutor)

// WithLogger 为 AsyncExecutor 设置自定义日志记录器
func WithLogger(logger *log.Logger) ExecutorOption {
	return func(e *AsyncExecutor) {
		e.logger = logger
	}
}

// WithContext 为 AsyncExecutor 设置自定义上下文
func WithContext(ctx context.Context) ExecutorOption {
	return func(e *AsyncExecutor) {
		e.ctx, e.cancel = context.WithCancel(ctx)
	}
}

// NewAsyncExecutor 创建一个新的 AsyncExecutor，指定工作者数量
func NewAsyncExecutor(workers int, options ...ExecutorOption) *AsyncExecutor {
	ctx, cancel := context.WithCancel(context.Background())
	e := &AsyncExecutor{
		workerPool: make(chan struct{}, workers),
		results:    make(chan Result),
		logger:     log.New(log.Writer(), "AsyncExecutor: ", log.LstdFlags),
		ctx:        ctx,
		cancel:     cancel,
	}

	for _, option := range options {
		option(e)
	}

	return e
}

// Execute 异步执行任务
func (e *AsyncExecutor) Execute(task Task) {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		select {
		case e.workerPool <- struct{}{}:
			defer func() { <-e.workerPool }()
			result := e.executeWithRecover(task)
			e.results <- result
		case <-e.ctx.Done():
			e.logger.Printf("由于上下文结束，任务执行被取消")
			e.results <- Result{Err: e.ctx.Err()}
		}
	}()
}

// executeWithRecover 执行任务并从 panic 中恢复
func (e *AsyncExecutor) executeWithRecover(task Task) (result Result) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			e.logger.Printf("发生 panic：%v\n%s", r, stack)
			result.Err = fmt.Errorf("发生 panic：%v", r)
		}
	}()

	value, err := task(e.ctx)
	return Result{Value: value, Err: err}
}

// ExecuteWithTimeout 在指定超时时间内执行任务
func (e *AsyncExecutor) ExecuteWithTimeout(task Task, timeout time.Duration) Result {
	ctx, cancel := context.WithTimeout(e.ctx, timeout)
	defer cancel()

	resultChan := make(chan Result, 1)
	go func() {
		result := e.executeWithRecover(func(innerCtx context.Context) (interface{}, error) {
			return task(innerCtx)
		})
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result
	case <-ctx.Done():
		return Result{Err: fmt.Errorf("任务在 %v 后超时", timeout)}
	}
}

// Wait 等待所有任务完成并关闭结果通道
func (e *AsyncExecutor) Wait() {
	go func() {
		e.wg.Wait()
		close(e.results)
	}()
}

// Results 返回提供所有任务结果的通道
func (e *AsyncExecutor) Results() <-chan Result {
	return e.results
}

// ExecuteAll 异步执行多个任务并返回所有结果
func (e *AsyncExecutor) ExecuteAll(tasks ...Task) []Result {
	for _, task := range tasks {
		e.Execute(task)
	}
	e.Wait()

	var results []Result
	for result := range e.Results() {
		results = append(results, result)
	}
	return results
}

// ExecuteAllWithTimeout 在指定超时时间内执行多个任务并返回可用结果
func (e *AsyncExecutor) ExecuteAllWithTimeout(timeout time.Duration, tasks ...Task) []Result {
	ctx, cancel := context.WithTimeout(e.ctx, timeout)
	defer cancel()

	for _, task := range tasks {
		e.Execute(task)
	}

	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range e.Results() {
			results = append(results, result)
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		e.logger.Printf("ExecuteAllWithTimeout：在 %v 后发生超时", timeout)
		return results
	case <-done:
		return results
	}
}

// Shutdown 优雅地关闭 AsyncExecutor
func (e *AsyncExecutor) Shutdown(timeout time.Duration) error {
	e.cancel() // 通知所有正在进行的任务停止

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("关闭在 %v 后超时", timeout)
	}
}
