package main

import (
	pb "github.com/bukhavtsov/gateway-jwt-auth/pkg/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	proxyAddr := ":8081"
	serviceAddr := "127.0.0.1:8082"
	HTTPProxy(proxyAddr, serviceAddr)
}

func HTTPProxy(proxyAddr string, serviceAddr string) {
	grpcConn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("failed to connect to grpc ", err)
	}
	defer grpcConn.Close()
	grpcGWMux := runtime.NewServeMux()
	err = pb.RegisterRestAppHandler(context.Background(), grpcGWMux, grpcConn)
	if err != nil {
		log.Fatalln("failed to start HTTP server", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/developers", grpcGWMux)
	log.Fatalln(http.ListenAndServe(proxyAddr, mux))
}

