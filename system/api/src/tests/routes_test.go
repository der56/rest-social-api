package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
)

func TestSetupAuthenticationRoutes(t *testing.T) {
	baseURL := "http://localhost:3000"

	// Generating random user data for registration
	username := randomdata.SillyName()
	password := randomdata.SillyName()
	email := randomdata.Email()
	firstname := randomdata.FirstName(randomdata.Male)
	lastname := randomdata.LastName()
	username2 := randomdata.SillyName()
	passwd2 := randomdata.SillyName()

	// Register user
	reqData := map[string]string{
		"Username":  username,
		"Password":  password,
		"Email":     email,
		"Firstname": firstname,
		"Lastname":  lastname,
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		t.Fatalf("Error converting to JSON: %v", err)
	}

	// Sending POST request for registration
	resp, err := http.Post(baseURL+"/api/register", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Error in registration request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Fatalf("Registration failed. Status: %v. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	t.Log("Registration successful, Status Code:", resp.StatusCode)

	// Login user
	reqLogin := map[string]string{
		"Username": username,
		"Password": password,
	}

	loginBody, err := json.Marshal(reqLogin)
	if err != nil {
		t.Fatalf("Error converting to JSON: %v", err)
	}

	// Sending POST request for login
	resp, err = http.Post(baseURL+"/api/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Error in login request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ = io.ReadAll(resp.Body)
	var loginData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &loginData)
	if err != nil {
		t.Fatalf("Error deserializing login response: %v", err)
	}

	// Extracting token from response header
	token := resp.Header.Get("Session")
	if token == "" {
		t.Fatalf("Failed to get token from 'Session' header")
	}

	t.Log("Login successful, Token:", token)

	// Request profile info
	reqProfile, err := http.NewRequest("GET", baseURL+"/api/profile", nil)
	if err != nil {
		t.Fatalf("Error creating profile request: %v", err)
	}

	reqProfile.Header.Set("Session", token)

	client := &http.Client{}
	resp, err = client.Do(reqProfile)
	if err != nil {
		t.Fatalf("Error in profile request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ = io.ReadAll(resp.Body)
	t.Log("Profile Response:", string(bodyBytes))

	var profileData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &profileData)
	if err != nil {
		t.Fatalf("Error deserializing profile response: %v", err)
	}

	// Parsing profile response
	type ProfileResponse struct {
		Message string `json:"message"`
		User    struct {
			ID        string `json:"ID"`
			Email     string `json:"email"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
			Picture   struct {
				String string `json:"String"`
				Valid  bool   `json:"Valid"`
			} `json:"picture"`
			Username string `json:"username"`
		} `json:"user"`
	}

	var profileResponse ProfileResponse
	err = json.Unmarshal(bodyBytes, &profileResponse)
	if err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	t.Log("User ID in profile:", profileResponse.User.ID)

	userID := profileResponse.User.ID
	if userID == "" {
		t.Fatalf("The 'ID' field is empty in the profile response")
	}

	// Test API routes
	tests := []struct {
		method       string
		route        string
		body         string
		expectedCode int
	}{
		{"POST", "/api/logout", "", 200},
		{"POST", "/api/login", `{"Username":"` + username + `", "Password":"` + password + `"}`, 200},
		{"GET", "/api/profile", "", 200},
		{"GET", "/api/profile/6a689342-6b5f-4a0e-a641-0c0d8a06b8cc", "", 200},
		{"POST", "/api/update-username", `{"username": "` + username2 + `"}`, 200},
		{"POST", "/api/update-password", `{"password": "` + passwd2 + `"}`, 200},
		{"POST", "/api/followuser/6a689342-6b5f-4a0e-a641-0c0d8a06b8cc", "", 200},
		{"POST", "/api/unfollowuser/6a689342-6b5f-4a0e-a641-0c0d8a06b8cc", "", 200},
		{"GET", "/api/getfollowers/" + userID, "", 200},
	}

	// Iterating over each test case to verify API routes
	for _, tt := range tests {
		var req *http.Request
		var err error

		// Create request based on method
		if tt.method == "POST" {
			req, err = http.NewRequest(tt.method, baseURL+tt.route, strings.NewReader(tt.body))
		} else {
			req, err = http.NewRequest(tt.method, baseURL+tt.route, nil)
		}

		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		if tt.method == "POST" || tt.method == "GET" {
			req.Header.Set("Session", token)
		}

		// Execute request and verify response
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error in HTTP request: %v", err)
		}
		defer resp.Body.Close()

		t.Log("Testing route:", tt.route, "Status:", resp.StatusCode)
		assert.Equal(t, tt.expectedCode, resp.StatusCode, "Expected status code to be %d for route %s", tt.expectedCode, tt.route)
	}
}
