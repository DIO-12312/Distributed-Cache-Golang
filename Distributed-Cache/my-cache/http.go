package mycache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

/*	HTTP服务端	*/

const defaultBasePath = "/_geecache/"

// self用于记录自己的地址，包括主机名/IP和端口
// basePath作为节点间通讯地址的前缀
/*
一个主机上还可能承载其他的服务
加一段 Path 是一个好习惯
比如，大部分网站的 API 接口，一般以 /api 作为前缀。
*/
type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 用interface输出可变参数模板v ...interface{}
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s]:%s", p.self, fmt.Sprintf(format, v...)) //格式化日志消息
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.SplitN(path[len(p.basePath):], "/", -1)

	if len(parts) > 2 {
		http.Error(w, "bad Request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group, ok := GetGroup(groupName)
	if !ok {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	value, err := group.Get(key)
	if err != nil {
		http.Error(w, "no such Key:"+key, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlices())
}

/*	HTTP客户端	*/
type httpGetter struct {
	baseURL string //需要访问的数据的远程节点地址
}

// 访问标定格式的地址，获取响应体的数据
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%s%s/%s", h.baseURL, group, key)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned:%s", res.Status)
	}
	//ioutil.ReadAll在1.16版本已被弃用
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil

}

// Go语言编译时检查接口的常用语句，检查特定类型是否真的实现了接口
var _ PeerGetter = (*httpGetter)(nil)
