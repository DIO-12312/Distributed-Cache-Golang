package mycache

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	//f的类型写成接口getter而非具体实现GetterFunc的原因
	//更灵活：可以自由转换成实现该接口的类型
	//类似于其他的语言中基类指针指向子类
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	except := []byte("key")
	//DeepEqual可以用于处理切片，map，结构体等复杂的内容的相等，以及nil不等于空切片
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, except) {
		t.Fatal("Not Passing,call failed")
	}
}

func TestNewGroup(t *testing.T) {
	_, err := NewGroup("test", 32, nil)
	if err == nil {
		t.Fatal("getter not be nil")
	}

}

func TestGetGroup(t *testing.T) {
	// 测试 1: 获取不存在的 Group，应该返回 nil 和 false
	g, ok := GetGroup("non-existent")
	if g != nil || ok {
		t.Errorf("GetGroup non-existent: got group=%v ok=%v, want nil false", g, ok)
	}

	// 测试 2: 创建一个 Group，然后成功获取它
	testGetter := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	createGroup, err := NewGroup("testGroup", 1024, testGetter)
	if err != nil {
		t.Fatalf("NewGroup failed: %v", err)
	}
	if createGroup == nil {
		t.Fatal("NewGroup failed: returned nil")
	}

	// 测试 3: 获取已存在的 Group
	retrievedGroup, ok := GetGroup("testGroup")
	if !ok {
		t.Error("GetGroup failed: should find 'testGroup' but got false")
	}
	if retrievedGroup == nil {
		t.Error("GetGroup failed: returned nil for existing group")
	}
	if retrievedGroup.name != "testGroup" {
		t.Errorf("GetGroup: got name=%s, want 'testGroup'", retrievedGroup.name)
	}

	// 测试 4: 验证返回的是同一个对象
	if retrievedGroup != createGroup {
		t.Error("GetGroup: returned different group object")
	}
}
