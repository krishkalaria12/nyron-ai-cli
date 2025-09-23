package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/krishkalaria12/nyron-ai-cli/config"
)

type WebSearchParams struct {
	Query string `json:"query"`
	GL    string `json:"gl,omitempty"`
	HL    string `json:"hl,omitempty"`
	Type  string `json:"type,omitempty"`
	Page  int    `json:"page,omitempty"`
}

type SearchParameters struct {
	Q           string `json:"q"`
	GL          string `json:"gl"`
	HL          string `json:"hl"`
	Autocorrect bool   `json:"autocorrect"`
	Page        int    `json:"page"`
	Type        string `json:"type"`
}

type KnowledgeGraph struct {
	Title             string            `json:"title"`
	Type              string            `json:"type"`
	Website           string            `json:"website"`
	ImageUrl          string            `json:"imageUrl"`
	Description       string            `json:"description"`
	DescriptionSource string            `json:"descriptionSource"`
	DescriptionLink   string            `json:"descriptionLink"`
	Attributes        map[string]string `json:"attributes"`
}

type Sitelink struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type OrganicResult struct {
	Title      string            `json:"title"`
	Link       string            `json:"link"`
	Snippet    string            `json:"snippet"`
	Date       string            `json:"date,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Sitelinks  []Sitelink        `json:"sitelinks,omitempty"`
	Position   int               `json:"position"`
}

type PeopleAlsoAsk struct {
	Question string `json:"question"`
	Snippet  string `json:"snippet"`
	Title    string `json:"title"`
	Link     string `json:"link"`
}

type RelatedSearch struct {
	Query string `json:"query"`
}

type WebSearchResult struct {
	Success          bool             `json:"success"`
	Message          string           `json:"message"`
	SearchParameters SearchParameters `json:"searchParameters"`
	KnowledgeGraph   *KnowledgeGraph  `json:"knowledgeGraph,omitempty"`
	Organic          []OrganicResult  `json:"organic"`
	PeopleAlsoAsk    []PeopleAlsoAsk  `json:"peopleAlsoAsk,omitempty"`
	RelatedSearches  []RelatedSearch  `json:"relatedSearches,omitempty"`
}

func WebSearch(params WebSearchParams) (WebSearchResult, ToolError) {
	apiKey := config.Config("SERPER_API_KEY")
	if apiKey == "" {
		return WebSearchResult{}, ToolError{
			Success: false,
			Message: "SERPER_API_KEY environment variable not set",
			Err:     nil,
		}
	}

	// Set defaults
	if params.GL == "" {
		params.GL = "us"
	}
	if params.HL == "" {
		params.HL = "en"
	}
	if params.Type == "" {
		params.Type = "search"
	}
	if params.Page == 0 {
		params.Page = 1
	}

	// Build URL
	baseURL := "https://google.serper.dev/search"
	u, err := url.Parse(baseURL)
	if err != nil {
		return WebSearchResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error parsing URL: %s", err.Error()),
			Err:     err,
		}
	}

	q := u.Query()
	q.Set("q", params.Query)
	q.Set("gl", params.GL)
	q.Set("hl", params.HL)
	q.Set("type", params.Type)
	if params.Page > 1 {
		q.Set("page", fmt.Sprintf("%d", params.Page))
	}
	q.Set("apiKey", apiKey)
	u.RawQuery = q.Encode()

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return WebSearchResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error creating request: %s", err.Error()),
			Err:     err,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return WebSearchResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error executing request: %s", err.Error()),
			Err:     err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WebSearchResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("API request failed with status: %d", resp.StatusCode),
			Err:     nil,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WebSearchResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error reading response: %s", err.Error()),
			Err:     err,
		}
	}

	var apiResponse struct {
		SearchParameters SearchParameters `json:"searchParameters"`
		KnowledgeGraph   *KnowledgeGraph  `json:"knowledgeGraph,omitempty"`
		Organic          []OrganicResult  `json:"organic"`
		PeopleAlsoAsk    []PeopleAlsoAsk  `json:"peopleAlsoAsk,omitempty"`
		RelatedSearches  []RelatedSearch  `json:"relatedSearches,omitempty"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return WebSearchResult{}, ToolError{
			Success: false,
			Message: fmt.Sprintf("Error parsing JSON response: %s", err.Error()),
			Err:     err,
		}
	}

	// Build result
	result := WebSearchResult{
		Success:          true,
		Message:          fmt.Sprintf("Search completed for query: %s", params.Query),
		SearchParameters: apiResponse.SearchParameters,
		KnowledgeGraph:   apiResponse.KnowledgeGraph,
		Organic:          apiResponse.Organic,
		PeopleAlsoAsk:    apiResponse.PeopleAlsoAsk,
		RelatedSearches:  apiResponse.RelatedSearches,
	}

	return result, ToolError{}
}
