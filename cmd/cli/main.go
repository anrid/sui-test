package main

import (
	"fmt"
	"strconv"

	"github.com/anrid/sui-test/pkg/sui"
)

func main() {
	sui.WaitForServer()

	{
		// Initialize client (if not yet initialized) and
		// show all configured environments
		sui.Exec("sui", "client", "-y", "envs")

		res := sui.CLI("client", "envs")
		envs := res["result"].([]interface{})
		if len(envs) < 1 {
			panic("should find at least 1 env")
		}
		println(sui.ToPrettyJSON(envs))
	}

	{
		// Setup local testnet
		sui.CLI("client", "new-env", "--alias", "local", "--rpc", sui.LocalValidatorRPCEndpoint)
		sui.CLI("client", "switch", "--env", "local")

		res := sui.CLI("client", "envs")
		envs := res["result"].([]interface{})
		if len(envs) < 2 {
			panic("should find at least 2 envs")
		}
		selectedEnv := envs[1].(string)
		fmt.Printf("Selected env: %s\n", selectedEnv)
	}

	var foundAddrs int
	var addr1 string
	var addr2 string
	{
		// Show all available wallet addresses
		res := sui.CLI("client", "addresses")
		addrs := res["addresses"].([]interface{})
		foundAddrs = len(addrs)
		if foundAddrs < 1 {
			panic("should find at least one wallet address")
		}
		println(sui.ToPrettyJSON(addrs))

		addr1 = (addrs[0].([]interface{}))[1].(string)
		if foundAddrs > 1 {
			addr2 = (addrs[1].([]interface{}))[1].(string)
		}
	}

	fmt.Printf("Addr 1: %s\n", addr1)

	if foundAddrs < 2 {
		// Create another wallet address if there's
		// currently only one
		res := sui.CLI("client", "new-address", "ed25519")
		println(sui.ToPrettyJSON(res))
		addr2 = res["address"].(string)

		sui.CLI("client", "switch", "--address", addr1)
	}

	fmt.Printf("Addr 2: %s\n", addr2)

	sui.CallFaucet(addr1)
	sui.CallFaucet(addr2)

	addr1Sui := new(SuiCoin)
	addr2Sui := new(SuiCoin)
	{
		// Ensure both wallets have some gas
		{
			res := sui.CLI("client", "gas", addr1)
			println(sui.ToPrettyJSON(res))
			coins := res["result"].([]interface{})
			coin1 := coins[0].(map[string]interface{})
			addr1Sui.ID = coin1["gasCoinId"].(string)
			addr1Sui.Balance = uint64(coin1["gasBalance"].(float64))
		}
		{
			res := sui.CLI("client", "gas", addr2)
			println(sui.ToPrettyJSON(res))
			coins := res["result"].([]interface{})
			coin1 := coins[0].(map[string]interface{})
			addr2Sui.ID = coin1["gasCoinId"].(string)
			addr2Sui.Balance = uint64(coin1["gasBalance"].(float64))
		}

		fmt.Printf("Addr 1 SUI: %s (%d)\n", addr1Sui.ID, addr1Sui.Balance)
		fmt.Printf("Addr 2 SUI: %s (%d)\n", addr2Sui.ID, addr2Sui.Balance)
	}

	// Transfer some SUI addr1 -> addr2
	var trans1CreatedID string
	{
		res := sui.CLI(
			"client", "transfer-sui", "--to", addr2,
			"--sui-coin-object-id", addr1Sui.ID,
			"--amount", "100000000",
			"--gas-budget", "10000000",
		)
		changes := res["objectChanges"].([]interface{})
		// println(sui.ToPrettyJSON(changes))

		mutated := changes[0].(sui.Map)
		if mutated["type"].(string) != "mutated" {
			panic("expected mutated object")
		}

		created := changes[1].(sui.Map)
		if created["type"].(string) != "created" {
			panic("expected created object")
		}

		trans1CreatedID = created["objectId"].(string)
	}

	// Merge newly transferred coin to main gas coin in addr2
	{
		sui.CLI("client", "switch", "--address", addr2)

		res := sui.CLI(
			"client", "merge-coin",
			"--primary-coin", addr2Sui.ID,
			"--coin-to-merge", trans1CreatedID,
			"--gas-budget", "10000000",
		)
		if _, found := res["balanceChanges"]; found {
			fmt.Printf("Merged newly transferred coin in addr %s\n", addr2)
		}
	}

	// Merge all small coins in addr2
	{
		res := sui.CLI("client", "objects", addr2)
		objs := res["result"].([]interface{})
		// println(sui.ToPrettyJSON(objs[0]))

		for _, obj := range objs {
			data := (obj.(sui.Map))["data"].(sui.Map)
			fields := (data["content"].(sui.Map))["fields"].(sui.Map)

			balance, err := strconv.ParseUint(fields["balance"].(string), 10, 64)
			if err != nil {
				panic(err)
			}
			id := data["objectId"].(string)
			typ := data["type"].(string)

			fmt.Printf("%s %s %d\n", id, typ, balance)

			if balance <= 100_000_000 && typ == sui.SuiCoinType {
				res := sui.CLI(
					"client", "merge-coin",
					"--primary-coin", addr2Sui.ID,
					"--coin-to-merge", id,
					"--gas-budget", "10000000",
				)
				if _, found := res["balanceChanges"]; found {
					fmt.Printf("Merged coins")
				}
			}
		}
	}
}

type SuiCoin struct {
	ID      string
	Balance uint64
}
