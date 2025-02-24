package jaeger

import (
	"fmt"
	"log"
	"net"

	jaeger2 "github.com/VictoriaMetrics/VictoriaMetrics/app/vlselect/jaeger"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/jaeger/proto"

	"google.golang.org/grpc"
)

// MustInit - init for Jaeger gRPC storage backend
func MustInit() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 17271))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterSpanWriterPluginServer(grpcServer, &SpanWriterPluginServer{})
	proto.RegisterSpanReaderPluginServer(grpcServer, &jaeger2.SpanReaderPluginServer{})

	go grpcServer.Serve(lis)
}

// MustStop - stop for Jaeger gRPC storage backend
func MustStop() {}
