// singleflight
package helper

import "sync"

// call 代表正在进行中，或已经结束的请求
type call struct {
	wg  sync.WaitGroup // 使用 sync.WaitGroup 锁避免重入
	val interface{}
	err error
}

// SingleFlightGroup singleflight 的主数据结构，管理不同 key 的请求(call)
type SingleFlightGroup struct {
	mu sync.Mutex // 保护 SingleFlightGroup 的成员变量 m 不被并发读写而加上的锁
	m  map[string]*call
}

// Do 针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误
func (g *SingleFlightGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		// 延迟初始化，提高内存使用效率
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果有其他请求正在进行中，则等待
		return c.val, c.err // 请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁，wg.Add(1) 锁加1；wg.Wait() 阻塞，直到锁被释放；wg.Done() 锁减1
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m
	g.mu.Unlock()

	return c.val, c.err // 返回结果
}
