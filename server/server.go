package server

import (
	"context"
	"jrgo-demo/model"
	"net"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/sirupsen/logrus"
)

var (
	log  *logrus.Entry
	addr = ":8086"
)

func init() {
	// Server export an object of type ExampleSvc.
	rpc.Register(&model.ExampleSvc{})

	logrus.SetLevel(logrus.DebugLevel)
}

func ServerTCP() {
	log = logrus.WithField("server", "tcp")

	// Server provide a TCP transport.
	lnTCP, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	//defer lnTCP.Close()
	log.Debugln("Start Listen:", addr)
	go func() {
		for {
			conn, err := lnTCP.Accept()
			if err != nil {
				log.Errorln("Accept err:", err.Error())
				return
			}
			ctx := context.WithValue(context.Background(), model.RemoteAddrContextKey, conn.RemoteAddr())
			go jsonrpc2.ServeConnContext(ctx, conn)
		}
	}()
}

func ServerHTTP() {
	log = logrus.WithField("server", "http")
	// Server provide a HTTP transport on /rpc endpoint.
	http.Handle("/rpc", jsonrpc2.HTTPHandler(nil))

	lnHTTP, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	log.Debugln("Start Listen:", addr)

	go http.Serve(lnHTTP, nil)
}
