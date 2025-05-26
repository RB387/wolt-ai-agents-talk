# Demystifying AI Agents: From Vanilla LLMs to Multi-Agent Systems

This repository contains Go code examples from a Wolt Golang Meetup talk, demonstrating the progression of building LLM-based agents. It showcases different approaches, from using vanilla LLM models to implementing complex multi-agent systems.

## Repository Structure and Components

This repository is structured to illustrate a clear progression in developing AI agents. Each directory represents a step in this journey:

### 1. `baseline/`
This directory demonstrates the initial approach of using vanilla Large Language Models (LLMs) directly via their APIs. It highlights the inherent issues and limitations of this method when no additional agentic logic is applied.
   - **Purpose**: To establish a starting point and showcase the challenges that lead to the development of more sophisticated agent structures.

### 2. `basic_agent/`
Here, you'll find a raw implementation of an AI agent written entirely in Go. This version is built from scratch without relying on external frameworks or libraries.
   - **Purpose**: To illustrate the fundamental components and logic required to build a functional AI agent.

### 3. `llm_chain/`
This directory presents an improved version of the AI agent. It utilizes a Go-based framework (specifics of the framework would ideally be mentioned if known, otherwise a generic statement is fine) to simplify the code structure and agent behavior.
   - **Purpose**: To show how frameworks can streamline agent development, making the code more manageable and robust.

### 4. `multi_agent/`
This section demonstrates a more advanced concept: composing multiple cooperating agents. This approach addresses challenges such as context window limitations, enabling agent specialization, and promoting modular design.
   - **Purpose**: To explore how complex problems can be tackled by breaking them down and assigning specialized tasks to different agents.

### 5. `mcp/`
The `MCP` (Model Context Protocol) directory illustrates how to build servers that expose additional tools and functionalities to LLM models. A key example provided is integration with an IDE (cursor for example).
   - **Purpose**: To showcase how agents can be equipped with external tools, significantly expanding their capabilities beyond text generation.

## Understanding the Progression

The examples in this repository are designed to provide a step-by-step understanding of building AI agents:

- **Starting Simple**: The `baseline` shows the raw power and limitations of LLMs.
- **Building Blocks**: `basic_agent` introduces fundamental agentic control logic.
- **Leveraging Frameworks**: `llm_chain` demonstrates how frameworks can abstract complexity.
- **Advanced Architectures**: `multi_agent` explores collaborative intelligence.
- **Extending Capabilities**: `mcp` shows how to empower agents with external tools.

## References

- **[LangGraph: Agent Architectures](https://langchain-ai.github.io/langgraph/concepts/agentic_concepts/)** - Overview of different agent architectures and concepts.
- **[LangGraph: Multi-Agent Systems](https://langchain-ai.github.io/langgraph/concepts/multi_agent)** - Explores breaking AI Agent into multiple smaller, independent agents to address complexity scaling issues.
- **[ReAct: Synergizing Reasoning and Acting in Language Models](https://arxiv.org/abs/2210.03629)** - This seminal paper explores using LLMs to generate both reasoning traces and task-specific actions in an interleaved manner. ReAct demonstrates how reasoning traces help models track and update action plans while handling exceptions, and how actions enable interfaces with external knowledge sources.
- **[Microsoft AI Agents for Beginners](https://github.com/microsoft/ai-agents-for-beginners)** - Comprehensive curriculum for learning about AI agents
- **[Anthropic: Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents)** - This guide shares practical insights from Anthropic's work with dozens of teams building LLM agents across industries. It emphasizes using simple, composable patterns rather than complex frameworks for building successful agent implementations.
- **[Anthropic: Model Context Protocol](https://www.anthropic.com/news/model-context-protocol)** - Anthropic's open-source protocol for connecting AI assistants to external data sources, tools, and environments. MCP provides a universal standard to replace fragmented integrations, allowing models to access relevant information across content repositories, business tools, and development environments.
- **[Coding the Coding Agents](https://github.com/zencoderai/coding-the-coding-agents)** - Python examples of coding multi-agents.
