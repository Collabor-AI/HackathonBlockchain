package main

import (
	"github.com/dgraph-io/badger"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

	
func main(){
	db, err := badger.Open(badger.DefaultOptions("tmp/badger"))
	  if err != nil {
		  log.Fatal(err)
	}
	bc := NewBlockChain(db)
	cli := CLI{bc}
	cli.Run()
	// bc.AddBlock("Block 1")
	// bc.AddBlock("Block 2")

	// bci := bc.Iterator()
	// while 
	// // bc.Print()

	

	  defer db.Close()
	

}


type CLI struct {
	bc *Blockchain
}


func (cli *CLI) Run(){
	// cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		_ = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		_ = printChainCmd.Parse(os.Args[2:])
	default:
		// cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock([]byte(data))
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	
	for {
		block := bci.Next()
		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("%+v",block)
			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Data: %+v\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			pow := NewProofOfWork(block)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			break
		}

		

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		
	}
}