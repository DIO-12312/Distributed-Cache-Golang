package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 哈希算法，默认为crc32.ChecksumIEEE
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int   //虚拟节点的个数
	keys     []int // 节点所对应的哈希值，
	hashMap  map[int]string
}

func NewMap(replicas int, hash Hash) *Map {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	//Ints对整数切片进行排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		//返回满足m.keys[i] >= hash的第一个下标
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
