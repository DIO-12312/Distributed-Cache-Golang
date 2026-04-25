package mycache

// 1. 选择节点的策略
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 2. 获取节点数据的方法
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
