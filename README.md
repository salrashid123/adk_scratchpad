## ADK Scratchpad

Just samples of using genai and gemini i'm using to learn

Note that for these examples, I used `Gemini Enterprise` and Application Default Credentials (ADC) (no api_key unless stated otherwise)

the intent of this is progressively add on more complex learning scenarios as i figure it out.  This is not a general tutorial and is pretty much just a log of my own progress

* `mcp_tool`:  very trivial demo of a custom tool, local mcp server and google search using 

---

To use any of the samples, you'll need GE enabled on the current project and ADC configured
 

```bash
gcloud auth application-default login
```

---

#### 1. mcp_tool

A trivial helloworld example of gemini enterprise ADK v2 which returns _static_ temperature in a city

    - instantiates a local stdio MCP server
    - instantiates a local tool
    - instantiates googlesearch tool 