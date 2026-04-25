package mycache

// 1. 选择节点的策略
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 2. 获取节点数据的方法
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

//可以分别改进选择策略和通信方式 ，降低耦合性
