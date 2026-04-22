// 这是 Go 的惯例：模块根目录下的 .go 文件，package名通常和模块名保持一致。
package mycache

import (
	"errors"
	"sync"
)

// 定义了当缓存未命中时如何获取数据的方法
type Getter interface {
	Get(key string) ([]byte, error)
}

// 设置回调函数
type GetterFunc func(key string) ([]byte, error)

//函数类型实现某一个接口，称之为接口型函数
//调用时能将传入函数作为参数，也能够传入实现接口的结构体作为参数

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/*核心数据结构Group*/
/*
	- 隔离不同缓存 - 不同的业务数据用不同的 Group
    - 独立管理 - 每个 Group 有自己的容量、过期策略等
	- 灵活扩展 - 后续可以为不同 Group 配置不同的分布式节点
*/

type Group struct {
	name      string //每个group都有自己的命名空间
	mainCache cache
	getter    Getter
}

var (
	//全局锁：保护不同的groups访问的线程安全
	//封装锁：保护Group中lru的不同缓存的线程安全
	mutex  sync.RWMutex
	groups = make(map[string]*Group) //指针确保能修改到真实对象
)

func NewGroup(name string, cacheBytes int64, getter Getter) (*Group, error) {
	if getter == nil {
		return nil, errors.New("No nil getter")
	}
	mutex.Lock()
	defer mutex.Unlock()
	g := &Group{
		name:      name,
		mainCache: cache{cacheBytes: cacheBytes},
		getter:    getter,
	}
	groups[name] = g
	return g, nil
}

func GetGroup(name string) (*Group, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	if g, ok := groups[name]; ok {
		return g, true
	}
	return nil, false
}
