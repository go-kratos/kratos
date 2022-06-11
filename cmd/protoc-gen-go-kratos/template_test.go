package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecute(t *testing.T) {

	tables := []struct {
		name   string
		data   ClientTemplate
		expect string
	}{
		{
			name: "single service",
			data: ClientTemplate{
				ClientInfoList: []ClientInfo{
					{
						ServiceName: "LibraryService",
						Endpoint:    "discovery://service_name",
					},
				},
			},
			expect: `type LibraryServiceGRPCClient struct {
	cli LibraryServiceClient
}

//NewLibraryServiceGRPCClient create grpc client for kratos
func NewLibraryServiceGRPCClient(opts ...grpc.ClientOption) (cli *LibraryServiceGRPCClient, err error) {
	opts = append(opts, grpc.WithBalancerName(wrr.Name))
	conn, ok := connMap["discovery://service_name"]
	if !ok {
		opts = append(opts, grpc.WithEndpoint("discovery://service_name"))
		conn, err = grpc.DialInsecure(context.Background(), opts...)
		if err != nil {
			return nil, err
		}
		connMap["discovery://service_name"] = conn
	}
	if err != nil {
		return nil, err
	}
	client := NewLibraryServiceClient(conn)
	return &LibraryServiceGRPCClient{cli:client}, nil
}`,
		},
		{
			name: "missing host",
			data: ClientTemplate{
				[]ClientInfo{
					{
						ServiceName: "LibraryService",
					},
				},
			},
			expect: `type LibraryServiceGRPCClient struct {
	cli LibraryServiceClient
}

//NewLibraryServiceGRPCClient create grpc client for kratos
func NewLibraryServiceGRPCClient(opts ...grpc.ClientOption) (cli *LibraryServiceGRPCClient, err error) {
	opts = append(opts, grpc.WithBalancerName(wrr.Name))
	conn, err := grpc.DialInsecure(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	client := NewLibraryServiceClient(conn)
	return &LibraryServiceGRPCClient{cli:client}, nil
}`,
		},
	}
	for _, item := range tables {
		t.Run(item.name, func(t *testing.T) {
			execute := item.data.execute()
			assert.Equalf(t, item.expect, execute, execute)
		})
	}
}
