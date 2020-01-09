## HackathonBlockchain

A small library in collaboration with https://github.com/topher-lo for future use as a backend framework in hackathons. The concept is to create a general-use blockchain that can store arbitrary data, but using 'Proof of Machine Learning' (https://github.com/topher-lo/go-tofu)  as a consensus mechanism, which is intended to utilise the blockchain to provide 'useful work'.

# Usage
Use `go get` to install all dependencies

Run the server `go run main.go`

Use `curl -d '$mydata' localhost:8081/$method` to call an API method

# Implemented Features
Create New Blockchain, Add Block, View Blockchain endpoints (still Proof Of Work)


# Development

Add decode/encode functions to transports/transport.go ; these specify how to read/send data from local to remote

Add endpoints ; these are a connector between transports and your application

Add application logic to services/


# Planned Features 
use GRPC to call a Python microservice; the Python microservice will fetch a model from IPFS/Google cloud and execute it on and obfuscated/homomorphic encrypted validation set
Replace PoW with PoML
Networking - use a Kafka cluster to provide message brokering as an initial step, then implement real P2P
