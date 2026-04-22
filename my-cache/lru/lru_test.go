package lru

import (
	"fmt"
	"testing"
)

// 实现 Value 接口的测试类型
type String string

func (s String) Len() int {
	return len(s)
}

// 测试 Add 和 Get 的正常使用
func TestAddAndGet(t *testing.T) {
	c := NewCacheLru(1024, nil)

	// 添加一条记录并读取
	c.Add("key1", String("hello"))
	if v, ok := c.Get("key1"); !ok || string(v.(String)) != "hello" {
		t.Fatal("获取 key1 失败")
	}

	// 获取不存在的 key
	if _, ok := c.Get("notExist"); ok {
		t.Fatal("不存在的 key 不应该返回 ok=true")
	}

	// 更新已有的 key
	c.Add("key1", String("world"))
	if v, ok := c.Get("key1"); !ok || string(v.(String)) != "world" {
		t.Fatal("key1 应该被更新为 world")
	}
}

// 测试超出容量时 RemoveOldest 自动触发 + 回调函数执行
func TestRemoveOldestAndCallback(t *testing.T) {
	evictedKeys := make([]string, 0)

	// 回调函数：记录被淘汰的 key
	onEvicted := func(key string, value Value) {
		evictedKeys = append(evictedKeys, key)
		fmt.Printf("回调触发: key=%s, value=%s 被淘汰\n", key, value.(String))
	}

	// maxBytes 设为 key+value 刚好容纳两条记录的大小
	// "k1" + "aa" = 2+2 = 4 字节
	// "k2" + "bb" = 2+2 = 4 字节
	// 总共 8 字节，maxBytes 设为 8
	c := NewCacheLru(8, onEvicted)

	c.Add("k1", String("aa")) // nowBytes = 4
	c.Add("k2", String("bb")) // nowBytes = 8

	// 再添加一条，容量不够，应淘汰最久未使用的 k1
	c.Add("k3", String("cc")) // 需要淘汰 k1 腾出空间

	if _, ok := c.Get("k1"); ok {
		t.Fatal("k1 应该已被淘汰")
	}

	if len(evictedKeys) != 1 || evictedKeys[0] != "k1" {
		t.Fatalf("回调应记录 k1 被淘汰，实际记录: %v", evictedKeys)
	}

	// k2 和 k3 应该还在
	if _, ok := c.Get("k2"); !ok {
		t.Fatal("k2 应该还在缓存中")
	}
	if _, ok := c.Get("k3"); !ok {
		t.Fatal("k3 应该还在缓存中")
	}
}
