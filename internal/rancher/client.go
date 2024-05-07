package rancher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// AuthenticatedTransport is a RoundTripper that adds Authorization header to each request
type AuthenticatedTransport struct {
	Token string
}

// ClusterResponse represents the JSON response structure for cluster data
type clusterResponse struct {
	Data []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

// ClusterTokenResponse represents the JSON response structure for cluster registration token data
type clusterTokenResponse struct {
	NodeCommand         string `json:"nodeCommand"`
	InsecureNodeCommand string `json:"insecureNodeCommand"`
}

// httpClient along with configuration for Rancher API
type rancherHttpClient struct {
	apiURL     string
	apiToken   string
	secure     bool
	httpClient *http.Client
}

// Constants
const rancherHttpAPIVersion string = "v3"

// Singleton resource used to send requests to the Rancher API
var httpClient *rancherHttpClient

// HttpClient initiates or returns the existing httpClient that can be used to communicate with Rancher REST API
func HttpClient(apiURL, apiToken string, secure bool) *rancherHttpClient {
	if httpClient == nil {
		httpClient = &rancherHttpClient{
			apiURL:   apiURL,
			apiToken: apiToken,
			secure:   secure,
			httpClient: &http.Client{
				Transport: &AuthenticatedTransport{
					Token: apiToken,
				},
			},
		}
	}
	return httpClient
}

// RetrieveClusterIDByName retrieves the ID of the cluster given its name
func (r *rancherHttpClient) RetrieveClusterIDByName(clusterName string) (string, error) {
	url := fmt.Sprintf("%s/%s/clusters", r.apiURL, rancherHttpAPIVersion)
	log.Info().Msgf("GET %s", url)
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var clusterResp clusterResponse
	err = json.NewDecoder(resp.Body).Decode(&clusterResp)
	if err != nil {
		return "", err
	}

	for _, cluster := range clusterResp.Data {
		if cluster.Name == clusterName {
			return cluster.ID, nil
		}
	}
	return "", fmt.Errorf("cluster with name '%s' not found", clusterName)
}

// RetrieveNodeCommand retrieves the shell registration command for the given cluster name
func (r *rancherHttpClient) RetrieveNodeCommand(clusterName string) (string, error) {
	clusterID, err := r.RetrieveClusterIDByName(clusterName)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/%s/clusters/%[3]s/clusterregistrationtokens/%[3]s:default-token", r.apiURL, rancherHttpAPIVersion, clusterID)
	log.Info().Msgf("GET %s", url)
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp clusterTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	if r.secure {
		return strings.TrimSpace(tokenResp.NodeCommand), nil
	}

	return strings.TrimSpace(tokenResp.InsecureNodeCommand), nil
}

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
func (t *AuthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	return http.DefaultTransport.RoundTrip(req)
}
