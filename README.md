## HackathonBlockchain

A small library in collaboration with https://github.com/topher-lo for future use as a backend framework in hackathons. The concept is to create a general-use blockchain that can store arbitrary data, but using 'Proof of Machine Learning' (https://github.com/topher-lo/go-tofu)  as a consensus mechanism, which is intended to utilise the blockchain to provide 'useful work'.

# Usage
Use `go get` to install all dependencies

Run the server `go run main.go`

Use `curl -d '$mydata' localhost:8081/{method}` to call an API method

## Methods (Implemented)
Create Blockchain

Add Block

Print Blockchain


# Development
Added communicawtion with Flask microservice, which will facilitate ML operations

# To Do
Improve data schema and block structure

Networking - use a Kafka cluster to provide message brokering as an initial step, then implement real P2P
