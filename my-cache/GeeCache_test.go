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
