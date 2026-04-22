package mycache

import (
	"fmt"
	"testing"
)

func TestCacheAddAndGet(t *testing.T) {
	c := &cache{cacheBytes: 100}

	// 测试 Add 和 Get
	key := "key1"
	value := NewByteView([]byte("value1"))
	c.add(key, value)

	// 获取存在的键
	got, ok := c.get(key)
	if !ok {
		t.Errorf("expected to get value for key %s", key)
	}
	if got.String() != value.String() {
		t.Errorf("expected %s, got %s", value.String(), got.String())
	}
}

func TestCacheGetNotFound(t *testing.T) {
	c := &cache{cacheBytes: 100}

	// 测试获取不存在的键
	got, ok := c.get("nonexistent")
	if ok {
		t.Errorf("expected not to find key, but got %v", got)
	}
}

func TestCacheMultipleKeys(t *testing.T) {
	c := &cache{cacheBytes: 1000}

	// 添加多个键值对
	testCases := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
	}

	for _, tc := range testCases {
		c.add(tc.key, NewByteView([]byte(tc.value)))
	}

	// 验证所有键都能正确检索
	for _, tc := range testCases {
		got, ok := c.get(tc.key)
		if !ok {
			t.Errorf("expected to get value for key %s", tc.key)
		}
		if got.String() != tc.value {
			t.Errorf("key %s: expected %s, got %s", tc.key, tc.value, got.String())
		}
	}
}

func TestCacheUpdateValue(t *testing.T) {
	c := &cache{cacheBytes: 100}

	key := "key1"
	value1 := NewByteView([]byte("value1"))
	value2 := NewByteView([]byte("value2"))

	// 添加初始值
	c.add(key, value1)
	got1, _ := c.get(key)
	if got1.String() != value1.String() {
		t.Errorf("expected %s, got %s", value1.String(), got1.String())
	}

	// 更新值
	c.add(key, value2)
	got2, _ := c.get(key)
	if got2.String() != value2.String() {
		t.Errorf("expected %s, got %s", value2.String(), got2.String())
	}
}

func TestCacheConcurrency(t *testing.T) {
	c := &cache{cacheBytes: 10000}
	done := make(chan bool)

	// 并发写入
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				value := NewByteView([]byte(key))
				c.add(key, value)
			}
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证数据完整性（随机检查几个）
	//避免检查"key_0_0"这类早期数据，10KB的cache会导致早期数据被清除
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_5_%d", i)
		got, ok := c.get(key)
		if !ok {
			t.Errorf("expected to get value after concurrent writes")
		}
		if got.String() != key {
			t.Errorf("data corruption: expected key_0_0, got %s", got.String())
		}
	}
}
