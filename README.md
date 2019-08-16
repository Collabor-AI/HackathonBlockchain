## HackathonBlockchain


#Usage
Use `go get` to install all dependencies

Run the server `go run main.go`

Use `curl -d '$mydata' localhost:8081/$method` to call an API method


# Development

Add decode/encode functions to transports/transport.go ; these specify how to read/send data from local to remote

Add endpoints ; these are a connector between transports and your application

Add application logic to services/
