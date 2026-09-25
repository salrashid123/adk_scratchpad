package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Name string `json:"name" jsonschema:"the name of the city to get the weather for"`
}

type Output struct {
	Weather string `json:"greeting" jsonschema:"the weather for the city"`
}

func GetWeather(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{Weather: fmt.Sprintf("the weather in %s is at most 75 degrees", input.Name)}, nil
}

func main() {

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp_weather", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "mcp_weather", Description: "Provides the current temperature in a city"}, GetWeather)

	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
