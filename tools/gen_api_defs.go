package main

// This code is seriously ugly but it'll do for now

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type oas struct {
	Paths map[string]any `json:"paths"`
}

type Methods map[string][]string

func isValidServiceName(inputString string) bool {
	pattern := "^[a-z-]*$"
	matched, _ := regexp.MatchString(pattern, inputString)
	return matched
}

func printHelp() {
	fmt.Println("Usage: gen_api_defs <service_name>")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("len: %d\n", len(os.Args))
		fmt.Printf("args: %v\n", os.Args)
		printHelp()
		return
	}

	serviceName := os.Args[1]

	if !isValidServiceName(serviceName) {
		fmt.Printf("Invalid service name '%s'\n", serviceName)
		fmt.Println("Service name must be lowercase and contain only letters and dashes")
		return
	}

	_, err := os.Stat(fmt.Sprintf("./build/%s", serviceName))
	if err != nil {
		fmt.Println("Service not found")
		fmt.Printf("Service %s not found in build directory: ./build/%s\n", serviceName, serviceName)
	}

	apiDefStart := fmt.Sprintf(`{
	new(gatewayContextName, apiHosts): {
	  name: '%s',
	  hosts: apiHosts,
	  defaultBackendURL: 'https://%s.europe-west2.%%s.gateway-contexts.thirdfort.internal' %% gatewayContextName,
	  routes: {
	`, serviceName, serviceName)
	apiDefEnd := fmt.Sprintf("  },\n}\n")

	//nolint:gosec,G204 // We trust the user to provide a valid service name
	cmd := exec.Command(fmt.Sprintf("./build/%s", serviceName))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	defer cmd.Process.Kill()

	// Wait for the service to start
	counter := 0
	var res *http.Response
	fmt.Printf("Waiting for service to start")
	for counter < 15 {
		res, err = http.Get("http://localhost:6060/openapi.json")
		if err == nil {
			defer res.Body.Close()
			break
		}
		fmt.Printf(".")
		counter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n\n")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var (
		oas  oas
		defs []Methods
	)
	err = json.Unmarshal(body, &oas)
	if err != nil {
		panic(err)
	}

	for path, details := range oas.Paths {
		d, ok := details.(map[string]any)
		if !ok {
			panic("Failed to get path details")
		}
		var allowed []string
		for method := range d {
			allowed = append(allowed, strings.ToUpper(method))
		}

		def := Methods{path: allowed}
		defs = append(defs, def)
	}

	var apiDef string
	indent := "  "
	for _, def := range defs {
		for path, methods := range def {
			apiDef += fmt.Sprintf("%s%s\"%s\"\n", indent, indent, path)
			apiDef += fmt.Sprintf("%s%s%s\"%s\": [\n", indent, indent, indent, "allowedMethods")
			for _, method := range methods {
				apiDef += fmt.Sprintf("%s%s%s%s\"%s\",\n", indent, indent, indent, indent, method)
			}
		}
	}
	fmt.Printf("%s%s%s", apiDefStart, apiDef, apiDefEnd)
}
