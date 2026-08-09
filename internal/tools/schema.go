package tools

import (
	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
)

func GetToolDefs() []llm.ToolDef {
	var toolDefs []llm.ToolDef

	for _, tool := range docker.DockerTools {
		def := llm.ToolDef{
			Name:        tool.Name,
			Description: tool.Description,
			Params:      tool.Params,
			Required:    tool.Required,
		}

		toolDefs = append(toolDefs, def)
	}

	return toolDefs
}
