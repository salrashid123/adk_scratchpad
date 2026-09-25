### Local MCP Server and Custom Tools

Sample one-shot binary (non-interactive) which creates an agent purportedly capable of returning the weather in a city and uses three tools:

a) The MCP server will accept a string which should be the city.  No matter what is sent, it return "the weather in %s is at most 75 degrees"

b) The local tool used by the agent which will accept a string which should be the city. No matter what is sent in, it return "the weather in %s is at least 55 degrees"

c) The base `geminitool.GoogleSearch{}` tool

1. Build the MCP Server
 
```bash
cd mcp_server/
go build -o my-mcp-server-binary main.go
```

2. Run the sample with all three tools

```bash
$ go run main.go 
The current temperature in San Francisco is between 55°F and 75°F.
```

This seems accurate since its a combination of assertions by the agents

---

Note, if you omit googleSearch, you'll what appears to be only one agent's response in the final set

```bash
$ go run main.go 
The current weather in San Francisco is around 75°F.
```

