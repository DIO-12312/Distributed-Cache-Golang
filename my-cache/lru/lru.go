package lru

import "container/list"

type Cache struct {
	maxBytes   int64
	nowBytes   int64
	ll         *list.List
	node_table map[string]*list.Element

	//某条记录被移除时的回调函数
	OnEvicted func(key string, vaue Value)
}

// 存入链表中，方便同时查看key和value两个值
type entry struct {
	key   string
	value Value
}

// 只要重写了这个函数的数据结构就能是Value
type Value interface {
	Len() int
}

// 相当于简化struct的构造
func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		ll:         list.New(),
		node_table: make(map[string]*list.Element),
		OnEvicted:  onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if e, ok := c.node_table[key]; ok {
		c.ll.MoveToFront(e)
		//更新value
		kv := e.Value.(*entry)
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		e := c.ll.PushFront(&entry{key, value})
		c.node_table[key] = e
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if e, ok := c.node_table[key]; ok {
		c.ll.PushFront(e)
		kv := e.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	e := c.ll.Back()
	if e != nil {
		c.ll.Remove(e)
		kv := e.Value.(*entry)
		delete(c.node_table, kv.key)
		c.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
