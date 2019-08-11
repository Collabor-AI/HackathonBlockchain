package endpoints


import (
	"context"
	"github.com/go-kit/kit/endpoint"
	// "github.com/dgraph-io/badger"
	"log"

	"HackathonBlockchain/services"
)


type Set struct {
	NewBlockchainEndpoint endpoint.Endpoint
}

func New(svc services.Service) Set {
	var newBlockchainEndpoint endpoint.Endpoint
	{
		newBlockchainEndpoint = MakeNewBlockchainEndpoint(svc)
	}
	return Set {
		NewBlockchainEndpoint: newBlockchainEndpoint,
	}
}

func (s Set) NewBlockchain(ctx context.Context, startingData services.InitData) (*services.Blockchain, error){
	resp, err := s.NewBlockchainEndpoint(ctx, NewBlockchainRequest{InitData:startingData,})
	if err != nil {
		log.Print("Failed to make new blockchain at endpoint")
		return  &services.Blockchain{}, err
	}
	response := resp.(NewBlockchainResponse)
	return response.Blockchain, response.Err
}

func MakeNewBlockchainEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewBlockchainRequest)
		bc, err := s.NewBlockchain(ctx, req.InitData)
		return NewBlockchainResponse{Blockchain: bc, Err: err}, nil
	}
}


type NewBlockchainRequest struct {
	InitData services.InitData
}

type NewBlockchainResponse struct {
	Blockchain *services.Blockchain `json:"blockchain"`
	Err error `json:"err",omitempty`
}