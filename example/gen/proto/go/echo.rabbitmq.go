package v1

import "context"
import "github.com/autom8ter/protoc-gen-rabbitmq/rpc"
import "github.com/golang/protobuf/proto"

// EchoServiceServer is a type safe RabbitMQ rpc server
type EchoServiceServer interface {
	// Echo returns the same message it receives.
	Echo(ctx context.Context, in *EchoRequest) (*EchoResponse, error)
}

// Serve starts the server and blocks until the context is canceled or the deadline is exceeded
func Serve(ctx context.Context, srv rpc.IServer, handler EchoServiceServer) error {
	return srv.Serve(ctx, func(ctx context.Context, msg rpc.Message) rpc.Message {
		meta := msg.Metadata
		switch msg.Method {
		case "Echo":
			var in EchoRequest
			if err := proto.Unmarshal(msg.Body, &in); err != nil {
				return rpc.Message{
					ID:       msg.ID,
					Method:   msg.Method,
					Metadata: meta,
					Error:    err,
				}
			}
			out, err := handler.Echo(rpc.NewContextWithMetadata(ctx, meta), &in)
			if err != nil {
				return rpc.Message{
					ID:       msg.ID,
					Method:   msg.Method,
					Metadata: meta,
					Error:    err,
				}
			}
			body, err := proto.Marshal(out)
			if err != nil {
				return rpc.Message{
					ID:       msg.ID,
					Method:   msg.Method,
					Metadata: meta,
					Error:    err,
				}
			}
			return rpc.Message{
				ID:       msg.ID,
				Method:   msg.Method,
				Metadata: meta,
				Body:     body,
			}
		}
		return rpc.Message{
			ID:       msg.ID,
			Method:   msg.Method,
			Metadata: meta,
			Error:    rpc.ErrUnsupportedMethod,
		}
	})
}

// EchoServiceClient is a type safe RabbitMQ rpc client
type EchoServiceClient struct {
	client rpc.IClient
}

// NewEchoServiceClient returns a new EchoServiceClientwith the given rpc client
func NewEchoServiceClient(client rpc.IClient) *EchoServiceClient {
	return &EchoServiceClient{client: client}
}

// Echo returns the same message it receives.
func (c *EchoServiceClient) Echo(ctx context.Context, in *EchoRequest) (*EchoResponse, error) {
	meta := rpc.MetadataFromContext(ctx)
	var out EchoResponse
	body, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}
	msg, err := c.client.Request(ctx, rpc.Message{Method: "Echo", Body: body, Metadata: meta})
	if err != nil {
		return nil, err
	}
	if msg.Error != nil {
		return nil, msg.Error
	}
	if err := proto.Unmarshal(msg.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
