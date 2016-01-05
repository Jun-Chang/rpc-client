package main

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
	"sync"

	libhttp "github.com/CyberAgent/car-golib/http"
	"github.com/Jun-Chang/rpc-client/proto"
	"github.com/Jun-Chang/rpc-client/service"
	msgpackrpc "github.com/msgpack-rpc/msgpack-rpc-go/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	http.HandleFunc("/monolithic/", monolithic)
	http.HandleFunc("/microservice/grpc/", grpcH)
	http.HandleFunc("/microservice/messagepack_rpc/", messagepackH)
	http.HandleFunc("/microservice/http/", httpH)
	http.ListenAndServe(":8080", nil)
}

func monolithic(w http.ResponseWriter, r *http.Request) {
	seq := service.Run()
	fmt.Fprint(w, "monolithic ", seq)
}

var grpcClient proto.TestServiceClient
var grpcOnce sync.Once

func grpcH(w http.ResponseWriter, r *http.Request) {
	grpcOnce.Do(func() {
		conn, err := grpc.Dial("127.0.0.1:11111", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		grpcClient = proto.NewTestServiceClient(conn)
	})
	res, err := grpcClient.Call(context.Background(), &proto.RequestType{})
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, "grpc ", res.Seq)
}

//var msgpackConn net.Conn
//var msgpackOnce sync.Once
func messagepackH(w http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("tcp", "127.0.0.1:11112")
	if err != nil {
		panic(err)
	}
	defer func() {
		conn.Close()
	}()
	msgpackClient := msgpackrpc.NewSession(conn, true)

	res, err := msgpackClient.Send("call")
	if err != nil {
		panic(err)
	}
	m := res.Interface().(map[interface{}]reflect.Value)
	var v reflect.Value
	for _, _v := range m {
		v = _v
		break
	}

	fmt.Fprint(w, "messagepack rpc ", v)
}

func httpH(w http.ResponseWriter, r *http.Request) {
	b, err := libhttp.Get("http://127.0.0.1/", map[string]string{})
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s %+v", "http ", b)
}
