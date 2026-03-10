package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mythingies/plugin-webex/internal/webex"
)

func registerListRecordings(s *mcpserver.MCPServer, client *webex.Client) {
	tool := mcp.NewTool("list_recordings",
		mcp.WithDescription("List Webex recordings, including recordings shared with you by other hosts. Returns recording IDs needed for get_recording_transcript."),
		mcp.WithString("from",
			mcp.Description("Start date/time in ISO 8601 format (default: 7 days ago)."),
		),
		mcp.WithString("to",
			mcp.Description("End date/time in ISO 8601 format (default: now)."),
		),
		mcp.WithNumber("max",
			mcp.Description("Maximum number of recordings to return (default 10)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		from := req.GetString("from", "")
		to := req.GetString("to", "")
		max := req.GetInt("max", 10)

		if from == "" {
			from = time.Now().UTC().Add(-7 * 24 * time.Hour).Format(time.RFC3339)
		}
		if to == "" {
			to = time.Now().UTC().Format(time.RFC3339)
		}

		recordings, err := client.ListRecordings(from, to, max)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list recordings: %v", err)), nil
		}

		if len(recordings) == 0 {
			return mcp.NewToolResultText("No recordings found in the specified time range."), nil
		}

		text := fmt.Sprintf("%d recording(s):\n\n", len(recordings))
		for _, r := range recordings {
			shared := ""
			if r.ShareToMe {
				shared = " [SHARED WITH YOU]"
			}
			dur := time.Duration(r.DurationSeconds) * time.Second
			text += fmt.Sprintf("- **%s**%s\n  Recording ID: %s\n  Meeting ID: %s\n  Host: %s\n  Recorded: %s\n  Duration: %s\n  Format: %s | Status: %s\n\n",
				r.Topic, shared, r.ID, r.MeetingID, r.HostEmail, r.TimeRecorded, dur, r.Format, r.Status)
		}
		return mcp.NewToolResultText(text), nil
	})
}
