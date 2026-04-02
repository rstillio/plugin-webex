package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mythingies/plugin-webex/internal/webex"
)

func registerShareFile(s *mcpserver.MCPServer, client *webex.Client) {
	tool := mcp.NewTool("share_file",
		mcp.WithDescription("Upload and share a local file as an attachment in a Webex space."),
		mcp.WithString("room_id",
			mcp.Required(),
			mcp.Description("The ID of the space to share the file in."),
		),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("The local file path to upload."),
		),
		mcp.WithString("message",
			mcp.Description("Optional text or markdown message to include with the file."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		roomID, err := req.RequireString("room_id")
		if err != nil {
			return mcp.NewToolResultError("room_id is required"), nil
		}
		filePath, err := req.RequireString("file_path")
		if err != nil {
			return mcp.NewToolResultError("file_path is required"), nil
		}

		message := req.GetString("message", "")

		msg, err := client.ShareFile(roomID, filePath, message)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to share file: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("File shared successfully (message ID: %s).", msg.ID)), nil
	})
}
