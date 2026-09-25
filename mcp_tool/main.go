package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/tool/geminitool"
	"google.golang.org/adk/v2/tool/mcptoolset"

	// "google.golang.org/adk/v2/cmd/launcher"
	// "google.golang.org/adk/v2/cmd/launcher/full"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"google.golang.org/genai"
)

var (
	projectID = flag.String("projectID", "redacted", "ProjectID with Gemini enabled")
	apiKey    = flag.String("apiKey", "redacted", "API KEy")
)

type MyArgs struct {
	Query string `json:"query" jsonschema:"The search query."`
}

type MyResult struct {
	Output string `json:"output"`
}

func GetWeather(ctx agent.Context, args MyArgs) (MyResult, error) {
	return MyResult{Output: fmt.Sprintf("the weather in %s is at least 55 degrees", args.Query)}, nil
}

func main() {

	flag.Parse()
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-3.8-flash", &genai.ClientConfig{
		//APIKey:  *apiKey", //os.Getenv("GOOGLE_API_KEY"),
		//Backend:  genai.BackendGeminiAPI, // requires api_key
		Project:  *projectID,
		Location: "us",
		// //Backend:  genai.BackendVertexAI,
		Backend: genai.BackendEnterprise,
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "weather_client", Version: "v1.0.0"}, nil)
	transport := &mcp.CommandTransport{Command: exec.Command("mcp_server/my-mcp-server-binary")}

	mcpConfig := mcptoolset.Config{
		Client:    client,
		Transport: transport,
	}

	ts, err := mcptoolset.New(mcpConfig)
	if err != nil {
		log.Fatalf("Failed to initialize MCP toolset: %v", err)
	}

	weatherTool, err := functiontool.New(functiontool.Config{
		Name:        "get_weather",
		Description: "Provides the current temperature in a city",
	}, GetWeather)
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	weatherAgent, err := llmagent.New(llmagent.Config{
		Name:        "weather_agent",
		Model:       model,
		Instruction: "You are a helpful assistant. Use your tools to provide accurate weather updates. Expect only the city name as input.",
		Description: "An agent capable of checking global weather conditions.",

		Tools: []tool.Tool{
			weatherTool,
			geminitool.GoogleSearch{},
		},
		Toolsets: []tool.Toolset{ts},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	r, err := runner.NewInMemory("my_app", weatherAgent)
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	msg := genai.NewContentFromText("san francisco", genai.RoleUser)

	events := r.Run(ctx, "user-123", "session-456", msg, agent.RunConfig{})
	for event, err := range events {
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}

		// Look for the final text block response from the agent execution
		if event.IsFinalResponse() && event.LLMResponse.Content != nil {
			for _, part := range event.LLMResponse.Content.Parts {
				fmt.Print(part.Text)
			}
			fmt.Println()
		}
	}

	// config := &launcher.Config{
	// 	AgentLoader: agent.NewSingleLoader(weatherAgent),
	// }

	// l := full.NewLauncher()
	// if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
	// 	log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	// }
}
