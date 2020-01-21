package endpoints


import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"fmt"
	// "github.com/dgraph-io/badger"
	"log"

	"HackathonBlockchain/services"
)


type Set struct {
	NewBlockchainEndpoint endpoint.Endpoint
	PrintBlockchainEndpoint endpoint.Endpoint
	AddBlockEndpoint endpoint.Endpoint
	GenerateAddressEndpoint endpoint.Endpoint
}

func New(svc services.Service) Set {
	var newBlockchainEndpoint endpoint.Endpoint
	{
		newBlockchainEndpoint = MakeNewBlockchainEndpoint(svc)
	}
	var printBlockchainEndpoint endpoint.Endpoint
	{
		printBlockchainEndpoint = MakePrintBlockchainEndpoint(svc)
	}
	var addBlockEndpoint endpoint.Endpoint
	{
		addBlockEndpoint = MakeAddBlockEndpoint(svc)
	}
	return Set {
		NewBlockchainEndpoint: newBlockchainEndpoint,
		PrintBlockchainEndpoint: printBlockchainEndpoint,
		AddBlockEndpoint: addBlockEndpoint,
	}
}

func (s Set) NewBlockchain(ctx context.Context, dataset services.Dataset, objective services.Objective, worldstate services.WorldState) ([]byte, error){
	resp, _ := s.NewBlockchainEndpoint(ctx, NewBlockchainRequest{dataset, objective, worldstate})
	response := resp.(NewBlockchainResponse)
	log.Print("Endpoint: %+v",response)
	return response.Blockchain, response.Err
}

func MakeNewBlockchainEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewBlockchainRequest)
		bc, err := s.NewBlockchain(ctx, services.InitData{req.Dataset,req.Objective,req.WorldState})
		bcData,_ := json.Marshal(bc.Tip)
		return NewBlockchainResponse{Blockchain: bcData, Err: err}, nil
	}
}


type NewBlockchainRequest struct {
	Dataset services.Dataset `json:"dataset,omitempty"` 
	Objective services.Objective `json:"objective,omitempty"`
	WorldState services.WorldState `json:"worldstate,omitempty"`
	
}


type NewBlockchainResponse struct {
	Blockchain []byte `json:"blockchain"`
	Err error `json:"err,omitempty"`
}


func (s Set) PrintBlockchain(ctx context.Context) ([]byte, error){
	resp, err := s.PrintBlockchainEndpoint(ctx, PrintBlockchainRequest{})
	if err != nil {
		return  nil, err
	}
	response := resp.(PrintBlockchainResponse)
	fmt.Printf("Response is %+v", response)
	return response.BlockchainIter, response.Err
}

func MakePrintBlockchainEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		bci, err := s.PrintBlockchain(ctx)
		bciData, _ := json.Marshal(bci.Blocks)
		return PrintBlockchainResponse{BlockchainIter: bciData, Err: err}, nil
	}
}

type PrintBlockchainRequest struct {
}

type PrintBlockchainResponse struct {
	BlockchainIter []byte 
	Err error 
}

func (s Set) AddBlock(ctx context.Context, Address string, Name string, Email string, Preds []byte, LinkToCode string, Description string, PrivKey string, Score float64) (error){
	resp, err := s.AddBlockEndpoint(ctx, AddBlockRequest{Address: Address, Name: Name, Email: Email, Preds: Preds, LinkToCode: LinkToCode, Description: Description, PrivKey: PrivKey, Score:Score})
	if err != nil {
		return err
	}
	response,_ := resp.(AddBlockResponse)
	return response.Err
}

func MakeAddBlockEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddBlockRequest)
		blockData := services.BlockData{req.Address, req.Name, req.Email, req.Preds, req.LinkToCode, req.Description, req.PrivKey}
		err = s.AddBlock(ctx, blockData,  req.Score)
		return AddBlockResponse{Err: err}, nil
	}
}

func GenerateAddressEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// err = s.GenerateAddress(ctx)
		return GenerateAddressResponse{Err: err}, nil
	}
}

type AddBlockRequest struct {
	Address string `json:"address"`
	Name string  `json:"name"`
	Email string `json:"email"`
	Preds []byte `json:"preds"`
	LinkToCode string `json:"linkToCode,omitempty"`
	Description string `json:"description,omitempty"`
	PrivKey string `json:"privateKey"`
	Score float64 `json:"score"`
}

type AddBlockResponse struct {
	Err error `json:"err,omitempty"`
}

type GenerateAddressRequest struct {
}

type GenerateAddressResponse struct{
	Err error `json:"err,omitempty"`
}