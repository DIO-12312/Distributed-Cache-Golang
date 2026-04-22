// 这是 Go 的惯例：模块根目录下的 .go 文件，package名通常和模块名保持一致。
package mycache

// 定义传入结构体接口
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
