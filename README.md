## HackathonBlockchain

Blockchain Framework and Protocol in collaboration with <a href = 'https://github.com/topher-lo'></a> . The aim is to create a framework that supports the creation of sub-blockchain networks, which represent a specific domain problem, usinga "Proof of Training" scheme (https://github.com/topher-lo/pyPoT)  as a consensus mechanism, purposing each blockchain as a distributed ledger of 'useful work' towards a problem.


# Usage
Use `go get` to install all dependencies

Run the server `go run main.go`


# Implemented Features
Create New Blockchain, Add Block, View Blockchain endpoints (still Proof Of Work)

## Methods (Implemented)

###Create Blockchain

This takes in a set of initial parameters, this is the problem definition that the blockchain network will aim to solve

```javascript
{
	"dataset":{}, //This represents the 'shared knowledge of the problem. It should contain a description and specific requirements for the data, and a link to a dataset on distributed /cloud storage (not on chain), although whether any data is shared will depend on the participants predefined rules
	"objective":{} // The objective of this blockchain, it should contain a specification of the task, and a scoring mechanism and baseline score if it is a supervised task
}



#### Add Block
{
}
This represents proposing a block, wherein a participant proposes an improvement in model performance. It should contain either the model, a way to access the model or an API endpoint for the specific model, or some means to validate the claim.


#### Print Blockchain

This returns an array of blockchain data for front-end /API usage

```javascript
{

}
```

# Development
Added communicawtion with Flask microservice, which will facilitate ML operations

# To Do
Improve data schema and block structure

Add application logic to services/


# Planned Features 
use GRPC to call a Python microservice; the Python microservice will fetch a model from IPFS/Google cloud and execute it on and obfuscated/homomorphic encrypted validation set
Docker
Replace PoW with PoML

Networking - use a Kafka cluster to provide message brokering as an initial step, then implement real P2P
