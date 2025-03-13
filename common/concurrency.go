package common

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Task 表示要执行的任务
type Task func() (interface{}, error)

// TaskWithContext 表示带有上下文的任务
type TaskWithContext func(ctx context.Context) (interface{}, error)

// TaskResult 表示任务的执行结果
type TaskResult struct {
	Value interface{}
	Error error
	Index int
}

// SafeCounter 线程安全的计数器
type SafeCounter struct {
	value int64
}

// Increment 增加计数器的值
func (c *SafeCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement 减少计数器的值
func (c *SafeCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Get 获取计数器的当前值
func (c *SafeCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set 设置计数器的值
func (c *SafeCounter) Set(value int64) {
	atomic.StoreInt64(&c.value, value)
}

// SafeMap 线程安全的映射
type SafeMap struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewSafeMap 创建一个新的线程安全映射
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]interface{}),
	}
}

// Set 设置键值对
func (m *SafeMap) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get 获取值，如果键不存在则返回nil
func (m *SafeMap) Get(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

// GetWithDefault 获取值，如果键不存在则返回默认值
func (m *SafeMap) GetWithDefault(key string, defaultValue interface{}) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if value, ok := m.data[key]; ok {
		return value
	}
	return defaultValue
}

// Has 检查键是否存在
func (m *SafeMap) Has(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.data[key]
	return ok
}

// Delete 删除键
func (m *SafeMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len 返回映射的大小
func (m *SafeMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Keys 返回所有键的列表
func (m *SafeMap) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有值的列表
func (m *SafeMap) Values() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	values := make([]interface{}, 0, len(m.data))
	for _, v := range m.data {
		values = append(values, v)
	}
	return values
}

// Clear 清空映射
func (m *SafeMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]interface{})
}

// ThreadPool 表示一个线程池
type ThreadPool struct {
	workers    int
	tasks      chan Task
	taskCount  SafeCounter
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
	stopped    bool
	mu         sync.Mutex
}

// NewThreadPool 创建一个新的线程池
func NewThreadPool(workers int, queueSize int) *ThreadPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &ThreadPool{
		workers:    workers,
		tasks:      make(chan Task, queueSize),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	// 启动工作线程
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go func(workerID int) {
			defer pool.wg.Done()
			pool.worker(workerID)
		}(i)
	}

	return pool
}

// worker 工作线程函数
func (p *ThreadPool) worker(id int) {
	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task()
			p.taskCount.Decrement()
		}
	}
}

// Submit 提交一个任务到线程池
func (p *ThreadPool) Submit(task Task) error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return errors.New("线程池已停止")
	}
	p.mu.Unlock()

	p.taskCount.Increment()
	select {
	case p.tasks <- task:
		return nil
	case <-p.ctx.Done():
		p.taskCount.Decrement()
		return errors.New("线程池已停止")
	}
}

// Stop 停止线程池
func (p *ThreadPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped {
		return
	}
	p.stopped = true
	p.cancelFunc()
	close(p.tasks)
	p.wg.Wait()
}

// Wait 等待所有任务完成
func (p *ThreadPool) Wait() {
	for p.taskCount.Get() > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

// Parallelizer 并行执行多个任务
type Parallelizer struct {
	MaxGoroutines int
	Timeout       time.Duration
}

// NewParallelizer 创建一个新的并行执行器
func NewParallelizer(maxGoroutines int, timeout time.Duration) *Parallelizer {
	return &Parallelizer{
		MaxGoroutines: maxGoroutines,
		Timeout:       timeout,
	}
}

// Run 并行执行任务列表
func (p *Parallelizer) Run(tasks []Task) []TaskResult {
	return p.RunWithContext(context.Background(), tasks)
}

// RunWithContext 带有上下文的并行执行任务列表
func (p *Parallelizer) RunWithContext(ctx context.Context, tasks []Task) []TaskResult {
	var (
		taskCount  = len(tasks)
		results    = make([]TaskResult, taskCount)
		wg         sync.WaitGroup
		sem        = make(chan struct{}, p.MaxGoroutines)
		ctxWithTO  context.Context
		cancelFunc context.CancelFunc
	)

	// 如果设置了超时，创建一个带超时的上下文
	if p.Timeout > 0 {
		ctxWithTO, cancelFunc = context.WithTimeout(ctx, p.Timeout)
	} else {
		ctxWithTO, cancelFunc = context.WithCancel(ctx)
	}
	defer cancelFunc()

	// 并行执行任务
	for i, task := range tasks {
		wg.Add(1)
		sem <- struct{}{}

		go func(index int, t Task) {
			defer func() {
				<-sem
				wg.Done()
			}()

			// 检查上下文是否已取消
			select {
			case <-ctxWithTO.Done():
				results[index] = TaskResult{
					Value: nil,
					Error: ctxWithTO.Err(),
					Index: index,
				}
				return
			default:
			}

			// 执行任务并捕获panic
			defer func() {
				if r := recover(); r != nil {
					results[index] = TaskResult{
						Value: nil,
						Error: fmt.Errorf("任务执行时发生panic: %v", r),
						Index: index,
					}
				}
			}()

			value, err := t()
			results[index] = TaskResult{
				Value: value,
				Error: err,
				Index: index,
			}
		}(i, task)
	}

	// 等待所有任务完成
	wg.Wait()
	return results
}

// RunTasksWithTimeout 在指定的超时时间内并行执行多个任务
func RunTasksWithTimeout(timeout time.Duration, tasks ...Task) []TaskResult {
	parallelizer := NewParallelizer(len(tasks), timeout)
	return parallelizer.Run(tasks)
}

// RunTasksConcurrently 并发执行多个任务
func RunTasksConcurrently(maxGoroutines int, tasks ...Task) []TaskResult {
	parallelizer := NewParallelizer(maxGoroutines, 0)
	return parallelizer.Run(tasks)
}

// Semaphore 信号量实现
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore 创建一个新的信号量
func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, size),
	}
}

// Acquire 获取信号量
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

// AcquireWithTimeout 在指定的超时时间内获取信号量
func (s *Semaphore) AcquireWithTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case s.sem <- struct{}{}:
		return true
	case <-timer.C:
		return false
	}
}

// Release 释放信号量
func (s *Semaphore) Release() {
	<-s.sem
}

// TryAcquire 尝试获取信号量，如果不可用则立即返回false
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// Available 返回当前可用的信号量数量
func (s *Semaphore) Available() int {
	return cap(s.sem) - len(s.sem)
}

// RateLimiter 速率限制器
type RateLimiter struct {
	interval time.Duration
	tokens   chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewRateLimiter 创建一个新的速率限制器
func NewRateLimiter(maxRequestsPerSecond int) *RateLimiter {
	interval := time.Second / time.Duration(maxRequestsPerSecond)
	ctx, cancel := context.WithCancel(context.Background())
	limiter := &RateLimiter{
		interval: interval,
		tokens:   make(chan struct{}, 1),
		ctx:      ctx,
		cancel:   cancel,
	}

	// 启动令牌生成器
	go limiter.generateTokens()

	return limiter
}

// generateTokens 定期生成令牌
func (r *RateLimiter) generateTokens() {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			select {
			case r.tokens <- struct{}{}:
			default:
				// 如果无法放入令牌（通道已满），则忽略
			}
		}
	}
}

// Wait 等待获取令牌
func (r *RateLimiter) Wait() {
	<-r.tokens
}

// TryWait 尝试获取令牌，如果不可用则立即返回false
func (r *RateLimiter) TryWait() bool {
	select {
	case <-r.tokens:
		return true
	default:
		return false
	}
}

// Close 关闭速率限制器
func (r *RateLimiter) Close() {
	r.cancel()
}

// RunWithRateLimit 使用速率限制执行函数
func RunWithRateLimit(maxRequestsPerSecond int, fn func()) {
	limiter := NewRateLimiter(maxRequestsPerSecond)
	defer limiter.Close()

	limiter.Wait()
	fn()
}

// BatchProcess 批量处理数据，将数据分成多个批次并行处理
func BatchProcess[T any](items []T, batchSize int, maxGoroutines int, processBatch func(batch []T) error) error {
	// 计算批次数
	batchCount := (len(items) + batchSize - 1) / batchSize
	batches := make([][]T, 0, batchCount)

	// 分割数据为多个批次
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}

	// 创建任务
	tasks := make([]Task, len(batches))
	for i, batch := range batches {
		batchCopy := batch // 创建副本避免闭包问题
		tasks[i] = func() (interface{}, error) {
			return nil, processBatch(batchCopy)
		}
	}

	// 并行执行任务
	results := RunTasksConcurrently(maxGoroutines, tasks...)

	// 检查是否有错误
	for _, result := range results {
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
