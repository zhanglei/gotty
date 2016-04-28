package main

import (
	"encoding/binary"
	"github.com/sumory/gotty/client"
	"github.com/sumory/gotty/codec"
	"github.com/sumory/gotty/config"
	log "github.com/sumory/log4go"
	"net"
	"time"
)

func handler(c *client.GottyClient, p *codec.Packet) {
	log.Info("客户端收到包, %v", p)
}

func dial(hostport string) (*net.TCPConn, error) {
	//连接
	remoteAddr, err := net.ResolveTCPAddr("tcp4", hostport)
	if nil != err {
		log.Error("ResolveTCPAddr err:", err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp4", nil, remoteAddr)
	if nil != err {
		log.Error("DiaTcp err:", hostport, err)
		return nil, err
	}

	return conn, nil
}

func main() {
	gottyConfig := config.NewDefaultGottyConfig()
	conn, _ := dial("localhost:6789")
	lengthBasedCodec := codec.NewLengthBasedCodec(binary.BigEndian, 64*1024, nil, nil)
	client := client.NewGottyClient(conn, lengthBasedCodec, gottyConfig, handler)
	client.Start()

	ch := make(chan int, 20)
	for {
		ch <- 1
		time.Sleep(1 * time.Second)
		go func() {
			header := &codec.PacketHeader{
				Sequence:  123,
				Operation: 1,
				Version:   0,
				Extra:     []byte("this is header extra"),
			}
			body := &codec.PacketBody{
				Data: []byte("this is body"),
			}
			meta := &codec.PacketMeta{
				TotalLen:  uint32(8 + header.Len() + body.Len()),
				HeaderLen: uint32(header.Len()),
			}

			p := &codec.Packet{
				Meta:   meta,
				Header: header,
				Body:   body,
			}

			log.Info("Client will write, totalLen: %d  headerLen: %d", p.Meta.TotalLen, p.Meta.HeaderLen)

			err := client.Write(p)
			if nil != err {
				log.Error("Client write failed: ", err)
			} else {
			}
			<-ch
		}()
	}

}
