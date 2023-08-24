package rpcwrapper

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mindwingx/abstraction"
	"github.com/mindwingx/go-helper"
	"net"
	"net/rpc"
)

type (
	rpcEngine struct {
		config   rpcConfig
		locale   abstraction.Locale
		entities []interface{}
	}

	rpcConfig struct {
		Network string
		Port    string
	}
)

func New(registry abstraction.Registry, locale abstraction.Locale) abstraction.RpcService {
	serviceRpc := new(rpcEngine)
	err := registry.Parse(&serviceRpc.config)
	if err != nil {
		helper.CustomPanic("", err)
	}

	serviceRpc.locale = locale

	return serviceRpc
}

func (r *rpcEngine) InitRpcService(rpcEntities []interface{}) {
	r.entities = append(r.entities, rpcEntities...)
}

func (r *rpcEngine) StartRpc() {
	get := r.locale.Get("rpc_start")
	color.Cyan(get)

	for _, entity := range r.entities {
		err := rpc.Register(entity)

		if err != nil {
			helper.CustomPanic(r.locale.Get("rpc_init_err"), err)
		}
	}

	listener, err := net.Listen(r.config.Network, fmt.Sprintf(":%s", r.config.Port))
	if err != nil {
		helper.CustomPanic(r.locale.Get("rpc_listen_err"), err)
	}

	defer listener.Close()

	for {
		rpcConn, acceptErr := listener.Accept()

		if acceptErr != nil {
			//todo: handle logger
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func (r *rpcEngine) Caller(destinationPort string, rpcMethod string, args interface{}, reply interface{}) (err error) {
	port := fmt.Sprintf(":%s", destinationPort)
	dial, err := rpc.DialHTTP(r.config.Network, port)
	if err != nil {
		//todo: call logger
		return
	}

	defer dial.Close()

	err = dial.Call(rpcMethod, args, reply)
	if err != nil {
		//todo: call logger
		return
	}

	return
}
