/*Package conner 读写数据

读函数，写函数，消息分发，打包机制，解析函数

*/
package conner

import (
	"net"

	"github.com/ranxx/ztcp/pkg/io/read"
	"github.com/ranxx/ztcp/pkg/io/write"
)

// Conner conn
type Conner interface {
	net.Conn
	ID() int64
	Start()
	Writer() write.Writer
	Reader() read.Reader
}
