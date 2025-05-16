# MCP Weather Service

A simple MCP service that provides current weather and 5-day forecast using the free OpenWeatherMap API.

## Features
- Current weather conditions
- 5-day weather forecast
- Configurable units (Celsius/Fahrenheit/Kelvin)
- Multi-language support
- Simple MCP integration

## Prerequisites
- Go 1.20+
- OpenWeatherMap API key

## Installation
1. Get OpenWeatherMap API key from [https://openweathermap.org/api](https://openweathermap.org/api)
2. Clone repository:
   ```bash
   git clone https://github.com/mschneider82/mcp-openweather.git
   cd mcp-openweather
   ```

3. Build:

   ```bash
   go build -o mcp-weather
   ```

### Installing via Smithery

To install MCP Weather Service for Claude Desktop automatically via [Smithery](https://smithery.ai/server/@mschneider82/mcp-openweather):

```bash
npx -y @smithery/cli install @mschneider82/mcp-openweather --client claude
```

## Configuration

Set required environment variable:

   ```bash
   export OWM_API_KEY="your_api_key_here"
   ```

## Usage

   ```json

    "mcpServers": {
        "mcp-openweather": {
            "command": "/home/YOURUSER/git/mcp-openweather/mcp-openweather",
            "env": {
                "OWM_API_KEY": "PUT_API_KEY_HERE"
            }
        }
    }

   ```

## MCP Parameters

   ```json
    {
        "tool": "weather",
        "params": {
            "city": "Berlin",          // Required
            "units": "c",             // Optional (c|f|k)
            "lang": "en"              // Optional (en|de|fr|...)
        }
    }
   ```

## Example Response

   ```txt
    Current weather for München:
        Conditions:  Klarer Himmel 
        Now:         17.78 metric
        High:        18.66 metric
        Low:         17.36 metric
        Pressure:    1017
        Humidity:    40
        FeelsLike:   16.66
        Wind Speed:  1.79
        Wind Degree:    61
        Sunrise:     1747107429 Unixtime
        Sunset:      1747161784 Unixtime
    Weather Forecast for München:
    Date & Time: 2025-05-13 12:00:00 +0000 UTC
    Conditions:  Clear Klarer Himmel
    Temp:        17.78 
    High:        18.93 
    Low:         17.78

    Date & Time: 2025-05-13 15:00:00 +0000 UTC
    Conditions:  Clear Klarer Himmel
    Temp:        18.33 
    High:        19.44 
    Low:         18.33

    Date & Time: 2025-05-13 18:00:00 +0000 UTC
    Conditions:  Clear Klarer Himmel
    Temp:        16.32 
    High:        16.32 
    Low:         15.59

    Date & Time: 2025-05-13 21:00:00 +0000 UTC
    Conditions:  Clear Klarer Himmel
    Temp:        9.18 
    High:        9.18 
    Low:         9.18
...
   ```

## Error Handling

# Common error scenarios:

    Missing API key: FATAL: OWM_API_KEY environment variable not set

    Invalid city: current weather error: failed to fetch current weather: city not found

    Network issues: failed to fetch forecast: context deadline exceeded

## Development

Contributions welcome! Please follow these steps:

    Fork the repository

    Create a feature branch

    Commit changes

    Push to the branch

    Create a Pull Request

## License

MIT License - see LICENSE for details
