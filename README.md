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


	{
	"dataset":{
		trainURL: "", //Link to download training data
		testURL: "", //Link to download testing data
		name: "" //name 
		"description": "" //Description of Data, Data Specification
		}, //This represents the 'shared knowledge of the problem. It should contain a description and specific requirements for the data, and a link to a dataset on distributed /cloud storage (not on chain), although whether any data is shared will depend on the participants predefined rules
	"objective":{ 
		"baseline": 0.0, float64 `json:"baseline,omitempty"` //baseline score, reject if scores below this
		scoring:"" //scoring method, string
		} 
	"worldstate":{ 
		"ensembleMethod":"" // Global Rules, choice of ensemble method
		} 
	}


#### Add Block
	{
		"blockData":{
			"address":"", //hash of pub key
			"email":"", //email for Kaggle Clone
			"trainPreds":"", //stored submitted predictions
			"description": "", // optional description
			},
		"pubKey": "", //'username'  on network
		"pubKey": "", // effectively 'password' on network (validate user)
		"Hash": ""  // identifier on 'leaderboard' and blockchain, hash of public key
	}
This represents proposing a block, wherein a participant proposes an improvement in model performance. It should contain either the model, a way to access the model or an API endpoint for the specific model, or some means to validate the claim.


#### Print Blockchain

This returns an array of blockchain data for front-end /API usage

```
[{"blockData":{""}}, {"blockData":{""}}]
```
