package model

import (
	"errors"
	"net"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
)

func init() {
	log = logrus.WithField("Exampel", "ExampleSvc")
	logrus.SetLevel(logrus.DebugLevel)
}

// A server wishes to export an object of type ExampleSvc:
type ExampleSvc struct{}

// Method with positional params.
func (*ExampleSvc) Sum(vals [2]int, res *int) error {
	*res = vals[0] + vals[1]
	return nil
}

// Method with positional params.
func (*ExampleSvc) SumAll(vals []int, res *int) error {
	for _, v := range vals {
		*res += v
	}
	return nil
}

// Method with named params.
func (*ExampleSvc) MapLen(m map[string]int, res *int) error {
	*res = len(m)
	return nil
}

type NameArg struct{ Fname, Lname string }
type NameRes struct{ Name string }

// Method with named params.
func (*ExampleSvc) FullName(t NameArg, res *NameRes) error {
	*res = NameRes{t.Fname + " " + t.Lname}
	return nil
}

var RemoteAddrContextKey = "RemoteAddr"

type NameArgContext struct {
	Fname, Lname string
	jsonrpc2.Ctx
}

// Method with named params and TCP context.
func (*ExampleSvc) FullName2(t NameArgContext, res *NameRes) error {
	host, _, _ := net.SplitHostPort(t.Context().Value(RemoteAddrContextKey).(*net.TCPAddr).String())
	log.Debugf("FullName2(): Remote IP is %s", host)
	*res = NameRes{t.Fname + " " + t.Lname}
	return nil
}

// Method with named params and HTTP context.
func (*ExampleSvc) FullName3(t NameArgContext, res *NameRes) error {
	host, _, _ := net.SplitHostPort(jsonrpc2.HTTPRequestFromContext(t.Context()).RemoteAddr)
	log.Debugf("FullName3(): Remote IP is %s", host)
	*res = NameRes{t.Fname + " " + t.Lname}
	return nil
}

// Method returns error with code -32000.
func (*ExampleSvc) Err1(struct{}, *struct{}) error {
	return errors.New("some issue")
}

// Method returns error with code 42.
func (*ExampleSvc) Err2(struct{}, *struct{}) error {
	return jsonrpc2.NewError(42, "some issue")
}

// Method returns error with code 42 and extra error data.
func (*ExampleSvc) Err3(struct{}, *struct{}) error {
	return &jsonrpc2.Error{42, "some issue", []string{"one", "two"}}
}
