package client

import (
	"io"
	"jrgo-demo/model"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/sirupsen/logrus"
)

var (
	addr = "localhost:8086"
	log  *logrus.Entry
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func ClientTCP() {
	log = logrus.WithField("client", "tcp")
	// Client use TCP transport.
	clientTCP, err := jsonrpc2.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer clientTCP.Close()

	var reply int

	// Synchronous call using positional params and TCP.
	err = clientTCP.Call("ExampleSvc.Sum", [2]int{3, 5}, &reply)
	log.Debugf("Sum(3,5)=%d", reply)

	// Asynchronous call using named params and TCP.
	startCall := clientTCP.Go("ExampleSvc.MapLen",
		map[string]int{"a": 10, "b": 20, "c": 30}, &reply, nil)
	replyCall := <-startCall.Done
	log.Debugf("MapLen({a:10,b:20,c:30})=%d", *replyCall.Reply.(*int))

	// Synchronous call using named params and TCP with context.
	clientTCP.Call("ExampleSvc.FullName2", model.NameArg{"First", "Last"}, nil)

	// Correct error handling.
	err = clientTCP.Call("ExampleSvc.Err1", nil, nil)
	if err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
		log.Errorf("Err1(): %q\n", err)
	} else if err != nil {
		rpcerr := jsonrpc2.ServerError(err)
		log.Errorf("Err1(): code=%d msg=%q data=%v\n", rpcerr.Code, rpcerr.Message, rpcerr.Data)
	}

}

func ClientHTTP() {
	log = logrus.WithField("client", "http")
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://" + addr + "/rpc")
	defer clientHTTP.Close()

	var (
		reply int
		err   error
	)
	// Synchronous call using positional params and HTTP.
	err = clientHTTP.Call("ExampleSvc.SumAll", []int{3, 5, -2}, &reply)
	if err != nil {
		log.Panicln(err)
	}
	log.Debugf("SumAll(3,5,-2)=%d", reply)

	// Notification using named params and HTTP.
	clientHTTP.Notify("ExampleSvc.FullName", model.NameArg{"First", "Last"})

	// Synchronous call using named params and HTTP with context.
	clientHTTP.Call("ExampleSvc.FullName3", model.NameArg{"First", "Last"}, nil)

	err = clientHTTP.Call("ExampleSvc.Err3", nil, nil)
	if err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
		log.Errorf("Err3(): %q\n", err)
	} else if err != nil {
		rpcerr := jsonrpc2.ServerError(err)
		log.Errorf("Err3(): code=%d msg=%q data=%v\n", rpcerr.Code, rpcerr.Message, rpcerr.Data)
	}
}

func ClientCusHTTP() {
	log = logrus.WithField("client", "cus-http")
	// Custom client use HTTP transport.
	clientCustomHTTP := jsonrpc2.NewCustomHTTPClient(
		"http://"+addr+"/rpc",
		jsonrpc2.DoerFunc(func(req *http.Request) (*http.Response, error) {
			// Setup custom HTTP client.
			client := &http.Client{}
			// Modify request as needed.
			req.Header.Set("Content-Type", "application/json-rpc")
			return client.Do(req)
		}),
	)
	defer clientCustomHTTP.Close()

	var err error

	err = clientCustomHTTP.Call("ExampleSvc.Err2", nil, nil)
	if err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
		log.Errorf("Err2(): %q\n", err)
	} else if err != nil {
		rpcerr := jsonrpc2.ServerError(err)
		log.Errorf("Err2(): code=%d msg=%q data=%v\n", rpcerr.Code, rpcerr.Message, rpcerr.Data)
	}
}
