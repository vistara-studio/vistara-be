package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client handles communication with vistara-ai service
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new AI service client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second, // AI requests may take longer
		},
	}
}

// SmartPlanRequest represents the travel planning request
type SmartPlanRequest struct {
	Destination         string    `json:"destination" validate:"required"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required"`
	Budget             *float64  `json:"budget,omitempty"`
	ActivityPreferences []string  `json:"activity_preferences,omitempty"`
	TravelStyle        *string   `json:"travel_style,omitempty"`
	ActivityIntensity  *string   `json:"activity_intensity,omitempty"`
	UserID             *string   `json:"user_id,omitempty"`
	IsPremium          *bool     `json:"is_premium,omitempty"`
	LocalBusinessIDs   []string  `json:"local_business_ids,omitempty"`
	AttractionIDs      []string  `json:"attraction_ids,omitempty"`
}

// SmartPlanResponse represents the AI-generated travel plan
type SmartPlanResponse struct {
	Plan              interface{} `json:"plan"`
	Destination       string      `json:"destination"`
	StartDate         time.Time   `json:"start_date"`
	EndDate           time.Time   `json:"end_date"`
	Budget            *float64    `json:"budget,omitempty"`
	TravelStyle       *string     `json:"travel_style,omitempty"`
	ActivityIntensity *string     `json:"activity_intensity,omitempty"`
	GeneratedAt       time.Time   `json:"generated_at"`
	UserID            *string     `json:"user_id,omitempty"`
}

// GenerateSmartPlan requests AI service to generate a smart travel plan
func (c *Client) GenerateSmartPlan(req *SmartPlanRequest) (*SmartPlanResponse, error) {
	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/api/v1/smart-planner", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for service-to-service communication
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Service", "vistara-be")
	httpReq.Header.Set("X-API-Key", c.APIKey)

	// Execute request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Success bool               `json:"success"`
		Message string             `json:"message"`
		Data    *SmartPlanResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("AI service error: %s", response.Message)
	}

	return response.Data, nil
}
