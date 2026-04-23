// 这是 Go 的惯例：模块根目录下的 .go 文件，package名通常和模块名保持一致。
package mycache

import (
	"errors"
	"fmt"
	"log"
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

/*
                            是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
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

func (g *Group) Get(key string) (value ByteView, err error) {
	if key == "" {
		//不建议用errors.New()过于轻量化
		return ByteView{}, fmt.Errorf("key must not nil")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

// 在未命中时候获取缓存
func (g *Group) load(key string) (value ByteView, err error) {
	//先只获取本地缓存
	return g.getLocally(key)
}

// 从用户设置的回调函数中获取缓存
func (g *Group) getLocally(key string) (value ByteView, err error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value = ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
