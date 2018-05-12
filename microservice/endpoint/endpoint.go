package endpoint

import "context"

type Endpoint interface {
	Handle(ctx context.Context, request []byte) (response []byte, err error)
}
