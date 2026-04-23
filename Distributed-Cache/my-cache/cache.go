// 实例化lru，对其封装一层互斥锁，并用饿汉模式，直到使用时才创建lru
package mycache

import (
	"mycache/lru"
	"sync"
)

type cache struct {
	lru        *lru.Cache
	mutex      sync.Mutex
	cacheBytes int64
}

// 懒汉模式实现函数
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCacheLru(c.cacheBytes, nil)
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), true
	}
	return
}

func (c *cache) add(key string, value ByteView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.lru == nil {
		c.lru = lru.NewCacheLru(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}
