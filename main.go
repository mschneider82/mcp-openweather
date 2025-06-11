package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"text/template"

	owm "github.com/briandowns/openweathermap"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// template used for output
const weatherTemplateTxt = `Current weather for {{.Name}}:
    Conditions: {{range .Weather}} {{.Description}} {{end}}
    Now:         {{.Main.Temp}} {{.Unit}}
    High:        {{.Main.TempMax}} {{.Unit}}
    Low:         {{.Main.TempMin}} {{.Unit}}
	Pressure:    {{.Main.Pressure}}
	Humidity:    {{.Main.Humidity}}
	FeelsLike:   {{.Main.FeelsLike}}
	Wind Speed:  {{.Wind.Speed}}
	Wind Degree:    {{.Wind.Deg}}
	Sunrise:     {{.Sys.Sunrise}} Unixtime
	Sunset:      {{.Sys.Sunset}} Unixtime
`

const forecastTemplateTxt = `Weather Forecast for {{.City.Name}}:
{{range .List}}Date & Time: {{.DtTxt}}
Conditions:  {{range .Weather}}{{.Main}} {{.Description}}{{end}}
Temp:        {{.Main.Temp}} 
High:        {{.Main.TempMax}} 
Low:         {{.Main.TempMin}}

{{end}}
`

// Pre-parse templates for better performance
var (
	weatherTemplate  = template.Must(template.New("weather").Parse(weatherTemplateTxt))
	forecastTemplate = template.Must(template.New("forecast").Parse(forecastTemplateTxt))
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"weather",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("weather",
		mcp.WithDescription("Get current and forecast weather information for a specific City"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("Location to get weather. If location has a space, wrap the location in double quotes."),
		),
		mcp.WithString("units",
			mcp.DefaultString("c"),
			mcp.Description("Temperature units (celsius|fahrenheit|kelvin) - default: c"),
		),
		mcp.WithString("lang",
			mcp.DefaultString("en"),
			mcp.Description("Language for weather descriptions - default: en"),
		),
	)

	// Add tool handler
	s.AddTool(tool, weatherHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// getCurrent gets the current weather for the provided
// location in the units provided.
func getCurrent(location, units, lang string) (*owm.CurrentWeatherData, error) {
	w, err := owm.NewCurrent(units, lang, os.Getenv("OWM_API_KEY"))
	if err != nil {
		return nil, err
	}
	w.CurrentByName(location)

	return w, nil
}

func getForecast5(location, units, lang string) (*owm.Forecast5WeatherData, error) {
	w, err := owm.NewForecast("5", units, lang, os.Getenv("OWM_API_KEY"))
	if err != nil {
		return nil, err
	}
	w.DailyByName(location, 99)
	forecast := w.ForecastWeatherJson.(*owm.Forecast5WeatherData)

	return forecast, err
}

func weatherHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	city, ok := request.Params.Arguments["city"].(string)
	if !ok {
		return nil, errors.New("city must be a string")
	}

	units, _ := request.Params.Arguments["units"].(string)
	if units == "" {
		units = os.Getenv("OWM_UNITS")
		if units == "" {
			units = "c" 
		}
	}
	
	lang, _ := request.Params.Arguments["lang"].(string)
	if lang == "" {
		lang = os.Getenv("OWM_LANG")
		if lang == "" {
			lang = "en" 
		}
	}

	current, err := getCurrent(city, units, lang)
	if err != nil {
		return nil, err
	}
	forecast, err := getForecast5(city, units, lang)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	if err := weatherTemplate.Execute(&buf, current); err != nil {
		return nil, err
	}
	if err := forecastTemplate.Execute(&buf, forecast); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s", buf.String())), err
}

// {"method":"tools/call","params":{"name":"weather","arguments":{"city":"Munich","lang":"de","units":"c"}},"jsonrpc":"2.0","id":9}
