/**
 * sync.Pool是一个可以存或取的临时对象集合; 可以安全被多个线程同时使用，保证线程安全;sync.Pool中保存的任何项都可能随时不做通知的释放掉，所以不适合用于像socket长连接或数据库连接池。
 * sync.Pool主要用途是增加临时对象的重用率，减少GC负担。
 */
package buffer

import (
	"bytes"
	"sync"
)

type Pool struct {
	pool sync.Pool
}

func (p *Pool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

func (p *Pool) Put(buf *bytes.Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}

func NewPool(s int) *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() interface{} {
				b := bytes.NewBuffer(make([]byte, s))
				b.Reset()
				return b
			},
		},
	}
}
