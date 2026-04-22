package byteview

//对真正存储底层数据的b []byte进行封装，提供一个只读的拷贝切片给外部
//实现真正的只读保护

type ByteView struct {
	b []byte
}

func NewByteView(data []byte) ByteView {
	b := make([]byte, len(data))
	copy(b, data)
	return ByteView{b: b}
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlices() []byte {
	res := make([]byte, len(v.b))
	copy(res, v.b)
	return res
}

// 输出时输出v.String()确保输出字符串而非底层的字节
func (v ByteView) String() string {
	return string(v.b)
}
