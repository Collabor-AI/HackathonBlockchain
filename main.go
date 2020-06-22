package main
//https://github.com/go-kit/kit/blob/master/examples/addsvc/cmd/addsvc/addsvc.go
import (
	"context"
	"github.com/dgraph-io/badger"
	"github.com/go-kit/kit/log"
	"flag"
	"github.com/oklog/oklog/pkg/group"
	"fmt"
	// "log"
	firebase "firebase.google.com/go"
	firebaseDB "firebase.google.com/go/db"
	// "firebase.google.com/go/auth"

	"google.golang.org/api/option"

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

 	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}




	var (
		httpAddr       = fs.String("http-addr", ":"+port, "HTTP listen address")
		// httpAddr       = fs.String("http-addr", "0.0.0.0:"+port, "HTTP listen address")
		//grpcAddr       = fs.String("grpc-addr", ":8082", "gRPC listen address")
	)
	fmt.Print(httpAddr)

	config := &firebase.Config{
	  DatabaseURL: "https://block-8f42c.firebaseio.com/",
	}

	CONF_JSON = "block-8f42c-firebase-adminsdk-5xrxz-0c78a97fb6.json"
	
	opt := option.WithCredentialsFile("block-8f42c-firebase-adminsdk-5xrxz-0c78a97fb6.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
	  fmt.Errorf("error initializing app: %v", err)
	}
	var client *firebaseDB.Client
	client,err = app.Database(ctx)


	db, err := badger.Open(badger.DefaultOptions("tmp/badger"))
	if err != nil {
		  fmt.Print("Failed to connect to db")
	}
 	defer db.Close()
 
	var (
		service        = services.New(db, app, client)
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
