package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Usage = "command-line client connect with Ethereum smart contracts"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "apikey",
			Usage: "etherscan.io apikey",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "show",
			Usage:  "Show available function",
			Action: showABI,
			ArgsUsage: `<address>
			Example : show 0xdac17f958d2ee523a2206206994597c13d831ec7`,
		},
		{
			Name:   "call",
			Usage:  "Call function",
			Action: callContract,
			ArgsUsage: `<address> <function name> <arguments...>
			Example : call 0xdac17f958d2ee523a2206206994597c13d831ec7 balanceOf 0x0000000000000000000000000000000000000000`,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showABI(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) < 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}
	var apikey string
	if ctx.GlobalIsSet("apikey") {
		apikey = ctx.GlobalString("apikey")
	}
	fmt.Println("Contract Address", args[0])
	var response EtherscanApi
	resp, err := http.Get("https://api.etherscan.io/api?module=contract&action=getabi&address=" + args[0] + "&apikey=" + apikey)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	if errUnmarshal := json.Unmarshal(body, &response); errUnmarshal != nil {
		fmt.Println("errorUnmarshal", errUnmarshal)
		return err
	}
	if response.Status == "0" {
		fmt.Println("Status", response.Status)
		fmt.Println("Message", response.Message)
		fmt.Println("Result", response.Result)
		return fmt.Errorf(response.Message)
	}
	abiPack, errPack := abi.JSON(strings.NewReader(response.Result))
	if errPack != nil {
		fmt.Println("errPack", errPack)
		return errPack
	}
	fmt.Println("Available method")
	for _, functionAbi := range abiPack.Methods {
		fmt.Println(functionAbi)
	}
	return err
}

func callContract(ctx *cli.Context) error {
	client := connect()
	args := ctx.Args()
	if len(args) < 2 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}
	var apikey string
	if ctx.GlobalIsSet("apikey") {
		apikey = ctx.GlobalString("apikey")
	}
	fmt.Println("Contract Address", args[0])
	contractAddress := args[0]
	var response EtherscanApi
	resp, err := http.Get("https://api.etherscan.io/api?module=contract&action=getabi&address=" + args[0] + "&apikey=" + apikey)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	if errUnmarshal := json.Unmarshal(body, &response); errUnmarshal != nil {
		fmt.Println("errorUnmarshal", errUnmarshal)
		return err
	}
	if response.Status == "0" {
		fmt.Println("Status", response.Status)
		fmt.Println("Message", response.Message)
		fmt.Println("Result", response.Result)
		return fmt.Errorf(response.Message)
	}
	abiPack, errPack := abi.JSON(strings.NewReader(response.Result))
	if errPack != nil {
		fmt.Println("errPack", errPack)
		return errPack
	}
	choseFunc := abiPack.Methods[args[1]]
	fmt.Println(choseFunc)
	var result string
	var req request
	var output []interface{}
	if len(choseFunc.Inputs) != 0 {
		var words []string
		if len(args) > 2 {
			for i := 2; i < len(args); i++ {
				fmt.Println("input", args[i])
				words = append(words, args[i])
			}
		}
		var resss []interface{}
		for i, parameter := range choseFunc.Inputs {
			cuyConvert, err := convertType(parameter.Type, words[i])
			if err != nil {
				log.Fatal(err)
			}
			resultValue := reflect.ValueOf(cuyConvert)
			resss = append(resss, resultValue.Interface())
		}
		jancokKoe, err := abiPack.Pack(choseFunc.Name, resss...)
		if err != nil {
			log.Fatal(err)
		}
		req = request{contractAddress, "0x" + common.Bytes2Hex(jancokKoe)}
		if err := client.Call(&result, "eth_call", req, "latest"); err != nil {
			log.Fatal(err)
		}
		output, err = abiPack.Unpack(choseFunc.Name, common.Hex2Bytes(result[2:]))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		sig := choseFunc.ID
		req = request{contractAddress, "0x" + common.Bytes2Hex(sig)}
		if err := client.Call(&result, "eth_call", req, "latest"); err != nil {
			log.Fatal(err)
		}
		output, err = abiPack.Unpack(choseFunc.Name, common.Hex2Bytes(result[2:]))
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("result", output)
	return err
}

func connect() *rpc.Client {
	client, err := rpc.DialHTTP("https://eth.llamarpc.com")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	return client
}
