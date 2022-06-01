package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pc-network/go-publicChain/utils"
	"pc-network/go-publicChain/wallet"
)

type CLI struct{}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateWallet -- create wallet")
	fmt.Println("\tgetAddressList -- print all wallet's address")
	fmt.Println("\tcreateBlockChain -address --genesis address")
	fmt.Println("\ttransfer -from FROM -to TO -amount AMOUNT -mine NOTMINE --transaction details")
	fmt.Println("\tprintChain -- output block's information")
	fmt.Println("\tgetBalance -address -- output balance")
	fmt.Println("\treset -- resetDatabase")
	fmt.Println("\tstartNode -miner ADDRESS -- start the node server,and specify the mining reward address")
}
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
func (cli CLI) Run() {

	isValidArgs()

	//get Node ID from env var
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID is not set\n")
		os.Exit(1)
	}
	fmt.Printf("NODE_ID:%s\n", nodeID)

	//custom command
	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
	getAddressListCmd := flag.NewFlagSet("getAddressList", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	transferBlockCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startNode", flag.ExitOnError)

	flagFrom := transferBlockCmd.String("from", "", "origin address")
	flagTo := transferBlockCmd.String("to", "", "destination address")
	flagAmount := transferBlockCmd.String("amount", "", "transfer amount")
	flagSendBlockVerify := transferBlockCmd.Bool("mine", false, "Whether to verify now")
	flagMiner := startNodeCmd.String("miner", "", "reward address")

	flagCreateBlockChainWithAddress := createBlockChainCmd.String("address", "", "create the address of genesis block")
	getBalanceWithAddress := getBalanceCmd.String("address", "", "inquire one's account")

	switch os.Args[1] {
	case "reset":
		err := resetCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getAddressList":
		err := getAddressListCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "transfer":
		err := transferBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockChain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startNode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	if resetCmd.Parsed() {

		cli.ResetDataBase(nodeID)
	}
	if getAddressListCmd.Parsed() {

		cli.GetAddressList(nodeID)
	}
	if createWalletCmd.Parsed() {

		cli.createWallet(nodeID)
	}
	if transferBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := utils.JSONToArray(*flagFrom)
		to := utils.JSONToArray(*flagTo)

		//verify the validity of address before transaction occurs
		for index, fromAddress := range from {
			if wallet.IsValidForAddress([]byte(fromAddress)) == false || wallet.IsValidForAddress([]byte(to[index])) == false {
				fmt.Println("Address is invalid")
				printUsage()
				os.Exit(1)
			}
		}

		amount := utils.JSONToArray(*flagAmount)
		cli.send(from, to, amount, nodeID, *flagSendBlockVerify)
	}
	if printChainCmd.Parsed() {
		//fmt.Println("output all blocks' information")
		cli.printChain(nodeID)
	}
	if createBlockChainCmd.Parsed() {

		if wallet.IsValidForAddress([]byte(*flagCreateBlockChainWithAddress)) == false {
			fmt.Println("address is invalid")
			printUsage()
			os.Exit(1)
		}
		cli.creatGenesisBlockChain(*flagCreateBlockChainWithAddress, nodeID)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceWithAddress == "" {
			fmt.Println("the address shouldn't be null")
			printUsage()
			os.Exit(1)
		}

		if wallet.IsValidForAddress([]byte(*getBalanceWithAddress)) == false {
			fmt.Println("address is invalid")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceWithAddress, nodeID)
	}
	if startNodeCmd.Parsed() {
		cli.StartNode(nodeID, *flagMiner)
	}
}
