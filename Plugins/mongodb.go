package Plugins

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/shadow1ng/fscan/common"
)

func MongodbScan(info *common.HostInfo, ch chan int, wg *sync.WaitGroup) {
	MongodbUnauth(info)
	wg.Done()
	<-ch
}

func MongodbUnauth(info *common.HostInfo) (flag bool, err error) {
	flag = false
	send_data := []byte{58, 0, 0, 0, 167, 65, 0, 0, 0, 0, 0, 0, 212, 7, 0, 0, 0, 0, 0, 0, 97, 100, 109, 105, 110, 46, 36, 99, 109, 100, 0, 0, 0, 0, 0, 255, 255, 255, 255, 19, 0, 0, 0, 16, 105, 115, 109, 97, 115, 116, 101, 114, 0, 1, 0, 0, 0, 0}
	getlog_data := []byte{72, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 212, 7, 0, 0, 0, 0, 0, 0, 97, 100, 109, 105, 110, 46, 36, 99, 109, 100, 0, 0, 0, 0, 0, 1, 0, 0, 0, 33, 0, 0, 0, 2, 103, 101, 116, 76, 111, 103, 0, 16, 0, 0, 0, 115, 116, 97, 114, 116, 117, 112, 87, 97, 114, 110, 105, 110, 103, 115, 0, 0}
	realhost := fmt.Sprintf("%s:%d", info.Host, common.PORTList["mgo"])
	conn, err := net.DialTimeout("tcp", realhost, time.Duration(info.Timeout)*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.Write(send_data)
	buf := make([]byte, 1024)
	count, err := conn.Read(buf)
	if err != nil {
		return flag, err
	}
	text := string(buf[0:count])
	if strings.Contains(text, "ismaster") {
		conn.Write(getlog_data)
		count, err := conn.Read(buf)
		if err != nil {
			return flag, err
		}
		text := string(buf[0:count])
		if strings.Contains(text, "totalLinesWritten") {
			flag = true
			result := fmt.Sprintf("Mongodb:%v unauthorized", realhost)
			common.LogSuccess(result)
		}
	}
	return flag, err
}
