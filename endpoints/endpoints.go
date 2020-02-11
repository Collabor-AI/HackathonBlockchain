package endpoints


import (
	"context"
	"encoding/json"
	// "encoding/base64"
	"github.com/go-kit/kit/endpoint"
	"fmt"
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
	var generateAddressEndpoint endpoint.Endpoint
	{
		generateAddressEndpoint = MakeGenerateAddressEndpoint(svc)
	}
	return Set {
		NewBlockchainEndpoint: newBlockchainEndpoint,
		PrintBlockchainEndpoint: printBlockchainEndpoint,
		AddBlockEndpoint: addBlockEndpoint,
		GenerateAddressEndpoint: generateAddressEndpoint,
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





func (s Set) PrintBlockchain(ctx context.Context, Hash string ) ([]byte, error){

	resp, err := s.PrintBlockchainEndpoint(ctx, PrintBlockchainRequest{Hash:Hash})
	if err != nil {
		return  nil, err
	}
	response := resp.(PrintBlockchainResponse)
	fmt.Printf("Response is %+v", response)
	return response.BlockchainIter, response.Err
}

func MakePrintBlockchainEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PrintBlockchainRequest)
		bci, err := s.PrintBlockchain(ctx, req.Hash)
		bciData, _ := json.Marshal(bci.Blocks)
		return PrintBlockchainResponse{BlockchainIter: bciData, Err: err}, nil
	}
}

func (s Set) AddBlock(ctx context.Context, PubKey string, Email string, Preds string,TrainPreds string, Description string, PrivKey string, Score float64, Hash string) (error){
	resp, err := s.AddBlockEndpoint(ctx, AddBlockRequest{PubKey:PubKey, Email: Email, TrainPreds: TrainPreds, Preds: Preds, Description: Description, PrivKey: PrivKey, Score:Score, Hash:Hash})
	if err != nil {
		return err
	}
	response,_ := resp.(AddBlockResponse)
	return response.Err
}

func MakeAddBlockEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddBlockRequest)
		log.Print("Endpoint - AddBlockRequest : %+v",req)
		blockData := services.BlockData{Email:req.Email, Preds:req.Preds, TrainPreds:req.TrainPreds, Description:req.Description}
		err = s.AddBlock(ctx, blockData,  req.Score, req.PubKey, req.PrivKey, req.Hash)
		return AddBlockResponse{Err: err}, nil
	}
}


func (s Set) GenerateAddress(ctx context.Context) ([]byte, error){
	resp, err := s.GenerateAddressEndpoint(ctx, GenerateAddressRequest{})
	if err != nil {
		return  nil, err
	}
	response := resp.(GenerateAddressResponse)
	fmt.Printf("Response is %+v", response)
	return response.Wallet, response.Err
}

func MakeGenerateAddressEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		wallet, err := s.GenerateAddress(ctx)
		fmt.Printf("'Endpoints':%+v",wallet)
		walletData, _ := json.Marshal(wallet)

		return GenerateAddressResponse{Wallet: walletData, Err: err}, nil
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

type PrintBlockchainRequest struct {
	Hash string `json:"hash"`
}

type PrintBlockchainResponse struct {
	BlockchainIter []byte 
	Err error 
}

type AddBlockRequest struct {
	PubKey string `json:"pubKey"`
	Email string `json:"email"`
	TrainPreds string `json:"trainPreds"`
	Preds string `json:"preds"`
	Description string `json:"description,omitempty"`
	PrivKey string `json:"privateKey"`
	Score float64 `json:"score"`
	Hash string `json:"hash"`
}

type AddBlockResponse struct {
	Err error `json:"err,omitempty"`
}

type GenerateAddressRequest struct {
}

type GenerateAddressResponse struct{
	Wallet []byte `json:"wallet"`
	Err error `json:"err,omitempty"`
}