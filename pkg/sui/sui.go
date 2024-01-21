package sui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const (
	DockerImage               = "sui-local"
	CLIBinary                 = "sui"
	LocalFaucetEndpoint       = "http://127.0.0.1:9123/gas"
	LocalValidatorRPCEndpoint = "http://127.0.0.1:9000"
)

func Request(method string) Map {
	return Map{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  nil,
	}
}

func WaitForServer() {
	req := ToJSON(Request("sui_getTotalTransactionBlocks"))

	for i := 0; i < 30; i++ {
		res, err := DockerExec("curl", "-d", req, "-H", "Content-Type: application/json", LocalValidatorRPCEndpoint)
		if err != nil {
			// if strings.Contains(string(res), "exit code:") {
			// 	break
			// }
			fmt.Printf("PING: error %s\n", string(res))
			time.Sleep(2_000 * time.Millisecond)
			continue
		}

		fmt.Printf("Server is ready: %s\n", string(res))
		break
	}
}

func CallFaucet(addr string) {
	res := CLI("client", "gas", addr)
	if len(res["result"].([]interface{})) > 0 {
		return
	}

	req := ToJSON(Map{"FixedAmountRequest": Map{"recipient": addr}})
	_, err := DockerExec("curl", "-d", req, "-H", "Content-Type: application/json", LocalFaucetEndpoint)
	if err != nil {
		panic(err)
	}
}

func CLI(command ...string) (res Map) {
	command = append([]string{"sui"}, command...)
	command = append(command, "--json")

	res = make(Map)

	out, err := DockerExec(command...)
	if err != nil {
		res["error"] = err.Error()
		return
	}

	if string(out)[0] == '[' {
		arr := make([]interface{}, 0)
		err := json.Unmarshal(out, &arr)
		if err != nil {
			panic(err)
		}

		res["result"] = arr
		return
	}

	if string(out)[0] == '"' {
		res["result"] = strings.Trim(string(out), "\"\n\r\t")
		return
	}

	err = json.Unmarshal(out, &res)
	if err != nil {
		fmt.Printf("Error: could not marshal to JSON: %s\n", string(out))
		panic(err)
	}

	return
}

func DockerExec(command ...string) (output []byte, err error) {
	command = append([]string{"exec", DockerImage}, command...)

	cmd := exec.Command("docker", command...)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error executing: %s\n", cmd.String())
	}
	return
}

func FromJSON(jsonBytes []byte) Map {
	m := make(Map)
	err := json.Unmarshal(jsonBytes, &m)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON:\n%s\n\n", string(jsonBytes))
		panic(err)
	}
	return m
}

type Map map[string]interface{}

func PostJSON(url string, data Map) (resBody Map, err error) {
	var body io.Reader
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewReader(b)
	}

	res, err := http.Post(url, "application/json", body)
	if err != nil {
		// fmt.Printf("ERROR: Sent JSON payload:\n%s\n\n", ToPrettyJSON(data))
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return FromJSON(resBytes), nil
}

func ToJSON(o interface{}) string {
	b, _ := json.Marshal(o)
	return string(b)
}

func ToPrettyJSON(o interface{}) string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}
