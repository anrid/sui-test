package main

import (
	"fmt"

	"github.com/anrid/sui-test/pkg/sui"
)

func main() {
	sui.WaitForServer()

	{
		// Initialize client (if not yet initialized) and
		// show all configured environments
		sui.DockerExec("sui", "client", "-y", "envs")

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
	var myAddr string
	var otherAddr string
	{
		// Show all available wallet addresses
		res := sui.CLI("client", "addresses")
		addrs := res["addresses"].([]interface{})
		foundAddrs = len(addrs)
		if foundAddrs < 1 {
			panic("should find at least one wallet address")
		}
		println(sui.ToPrettyJSON(addrs))

		myAddr = (addrs[0].([]interface{}))[1].(string)
		if foundAddrs > 1 {
			otherAddr = (addrs[1].([]interface{}))[1].(string)
		}
	}

	fmt.Printf("My addr    : %s\n", myAddr)

	if foundAddrs < 2 {
		// Create another wallet address if there's
		// currently only one
		res := sui.CLI("client", "new-address", "ed25519")
		println(sui.ToPrettyJSON(res))
		otherAddr = res["address"].(string)

		sui.CLI("client", "switch", "--address", myAddr)
	}

	fmt.Printf("Other addr : %s\n", otherAddr)

	sui.CallFaucet(myAddr)
	sui.CallFaucet(otherAddr)

	{
		// Ensure both wallets have some gas
		res := sui.CLI("client", "gas", myAddr)
		println(sui.ToPrettyJSON(res))
	}
}
