package observability

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SendErrorToTeams(c *fiber.Ctx, err error) {
	webhookUri := os.Getenv("TEAMS_WEBHOOK_URI")
	if webhookUri == "" {
		return
	}

	// Pre Processing
	requestHeader, _ := json.Marshal(c.GetReqHeaders())
	queryParams, _ := url.QueryUnescape(string(c.Request().URI().QueryString()))

	// Payload
	var payload TeamsMessage
	payload.Summary = "Internal Server Error"
	payload.Title = "Internal Server Error"
	payload.ThemeColor = "FF0000"
	payload.Sections = []Sections{
		{
			Text: "Error:",
			Facts: []Content{
				{Name: "Message", Value: fmt.Sprintf("```json\n%s\n```", err.Error())},
				{Name: "Stack Trace", Value: fmt.Sprintf("```json\n%s\n```", c.Locals("STACKTRACE"))},
			},
		},
		{
			Text: "<br>",
		},
		{
			Text: "Basic Info:",
			Facts: []Content{
				{Name: "App", Value: "Order API"},
				{Name: "Env", Value: os.Getenv("APP_ENV")},
				{Name: "Endpoint", Value: fmt.Sprintf("[%s] %s", c.Route().Method, c.Route().Path)},
			},
		},
		{
			Text: "<br>",
		},
		{
			Text: "Request Detail:",
			Facts: []Content{
				{Name: "IP", Value: c.IP()},
				{Name: "User Agent", Value: c.Get("User-Agent")},
				{Name: "URL", Value: getRequestURI(c)},
				{Name: "Header", Value: fmt.Sprintf("```json\n%s\n```", requestHeader)},
				{Name: "Body", Value: fmt.Sprintf("```json\n%s\n```", c.Request().Body())},
				{Name: "Query Params", Value: fmt.Sprintf("```json\n%s\n```", strings.ReplaceAll(queryParams, "&", "\n"))},
			},
		},
		{
			Text: "<br>",
		},
		{
			Text: "Replicate Action:<br>",
			Facts: []Content{
				{
					Name:  "CURL",
					Value: fmt.Sprintf("```bash\n%s\n```", getCurl(c)),
				},
			},
			Markdown: true,
		},
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, _ := http.NewRequest("POST", webhookUri, bytes.NewBuffer(payloadJson))

	req.Header.Add("Content-Type", "application/json")

	go http.DefaultClient.Do(req)
}

func getCurl(c *fiber.Ctx) string {
	// Retrieve request information from Fiber context
	method := c.Method()

	// Construct curl command
	curlCommand := fmt.Sprintf("curl -X %s --url '%s'", method, getRequestURI(c))

	// Add headers
	headers := make([]string, 0)

	reqHeaders := c.GetReqHeaders()
	for key, value := range reqHeaders {
		if key != "Host" {
			headers = append(headers, fmt.Sprintf("-H '%s: %s'", key, value))
		}
	}
	if len(headers) > 0 {
		curlCommand += " " + strings.Join(headers, " ")
	}

	// Add request body
	body := string(c.Request().Body())
	if body != "" {
		escapedBody := strings.ReplaceAll(body, `'`, `'\''`)
		curlCommand += fmt.Sprintf(" -d '%s'", escapedBody)
	}

	return curlCommand
}

func getRequestURI(c *fiber.Ctx) string {
	path := strings.Replace(string(c.Request().RequestURI()), "/external", "", 1)

	host := c.BaseURL()
	if os.Getenv("APP_BASE_URL") != "" {
		host = os.Getenv("APP_BASE_URL")
	}

	return host + path
}

type TeamsMessage struct {
	Summary    string     `json:"summary,omitempty"`
	Title      string     `json:"title,omitempty"`
	ThemeColor string     `json:"themeColor,omitempty"`
	Sections   []Sections `json:"sections,omitempty"`
}

type Sections struct {
	Text     string    `json:"text,omitempty"`
	Facts    []Content `json:"facts,omitempty"`
	Markdown bool      `json:"markdown,omitempty"`
}

type Content struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
