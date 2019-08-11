package main
//https://github.com/go-kit/kit/blob/master/examples/addsvc/cmd/addsvc/addsvc.go
import (
	"github.com/dgraph-io/badger"
	"github.com/go-kit/kit/log"
	"flag"
	"github.com/oklog/oklog/pkg/group"
	"fmt"
	// "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"HackathonBlockchain/endpoints"
	"HackathonBlockchain/services"
	"HackathonBlockchain/transports"
)

	
func main(){

 	fs := flag.NewFlagSet("svc", flag.ExitOnError)
	var (
		httpAddr       = fs.String("http-addr", ":8081", "HTTP listen address")
		//grpcAddr       = fs.String("grpc-addr", ":8082", "gRPC listen address")
	)



	db, err := badger.Open(badger.DefaultOptions("tmp/badger"))
	if err != nil {
		  fmt.Print("Failed to connect to db")
	}
 	defer db.Close()
 
	var (
		service        = services.New(db)
		endpoints      = endpoints.New(service)
		httpHandler    = transports.NewHTTPHandler(endpoints)
		// grpcServer     = addtransport.NewGRPCServer(endpoints, tracer, zipkinTracer, logger)

	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())


}
