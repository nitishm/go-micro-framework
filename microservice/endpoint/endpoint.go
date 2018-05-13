package endpoint

import "context"

// Endpoint provides the main interface to be implemented by each endpoint.
// Endpoints are where the user provides all the business logic of the microservice.
type Endpoint interface {
	Handle(ctx context.Context, request []byte) (response []byte, err error)
}
