package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const weatherURL = "https://api.open-meteo.com/v1/forecast?latitude=10.8231&longitude=106.6297&current=temperature_2m"

func main() {
	weatherTool := mcp.NewTool(
		"get_hcmc_temperature",
		mcp.WithDescription("Get the current temperature in Ho Chi Minh City in Celsius"),
	)

	mcpServer := server.NewMCPServer("Weather Example", "0.1.0", server.WithToolCapabilities(true))
	mcpServer.AddTool(weatherTool, getTemperature)

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintln(os.Stderr, "server error:", err)
	}
}

func getTemperature(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, nil)
	if err != nil {
		return mcp.NewToolResultError("create weather request: " + err.Error()), nil
	}

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return mcp.NewToolResultError("fetch weather: " + err.Error()), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError("weather service returned " + resp.Status), nil
	}

	var data struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return mcp.NewToolResultError("decode weather response: " + err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Ho Chi Minh City: %.1f°C", data.Current.Temperature)), nil
}
