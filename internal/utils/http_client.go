package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	client     *http.Client
	clientOnce sync.Once
)

func getHttpClient() *http.Client {
	clientOnce.Do(func() {
		jar, _ := cookiejar.New(nil)
		client = &http.Client{Jar: jar}
	})
	return client
}

func internalGet(url, authorization string) (map[string]any, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := getHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response: %s", body)
	}

	var result map[string]any
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return result, nil
}

func internalPost(url, authorization string, payload map[string]any) (map[string]any, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := getHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response: %s", body)
	}

	var result map[string]any
	if resp.StatusCode != http.StatusNoContent {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		}
	}
	return result, nil
}

func internalPutFile(url, authorization, file string) error {

	body, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open backupset '%s' due to: %w", file, err)
	}

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")

	resp, err := getHttpClient().Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response: %s", body)
	}

	return nil
}

func internalDelete(url, authorization string) error {

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := getHttpClient().Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response: %s", body)
	}

	return nil
}

func buildRootAuthorizationHeader(host string, config map[string]any) (string, error) {

	// Grab root client id and secret to compose the value for the authorization header
	rootClientID, err := GetRootClientIDFromHostConfig(host, config)
	var rootClientSecret string
	if err == nil {
		rootClientSecret, err = GetRootClientSecretFromHostConfig(host, config)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(rootClientID+":"+rootClientSecret))), nil
}

func buildUserAuthorizationHeader(user, password string) (string, error) {

	// Grab the user name and password for the user to compose the value for the authorization header
	// Note that user, if specified, is presumed to identify a registered/defined user. If such user
	// doesn't exist it is used as the user name instead. The password overwrites any password that
	// might be specified in the configuration. If user references a specified user that password is
	// used, if it didn't than a password needs to be specified as well if that user has one.

	var userName, userPassword string

	// Get the user name
	user, err := GetUserName(user)
	if err != nil {
		return "", err
	}

	// Lookup the user in list of configured users
	users := viper.GetStringMap("users")
	raw := users[user]
	if raw == nil {
		// Use the user name and password as provided
		userName = user
		userPassword = password
	} else {
		// Use the data for the specified user, using the password if one provided
		userMap, ok := raw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("invalid configuration for user '%s', format invalid", user)
		}
		// Name is not needed if it's the same as the user we set it on/for
		raw = userMap["name"]
		if raw != nil && raw.(string) != "" {
			userName = raw.(string)
		} else {
			userName = user
		}
		// Password if provided is used over the one specified for the user
		if password == "" {
			raw = userMap["password"]
			if raw != nil {
				userPassword = raw.(string)
			}
		} else {
			userPassword = password
		}
	}
	return fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(userName+":"+userPassword))), nil
}

func ManageAPIGet(host, path string) (map[string]any, error) {

	// Lookup the host's configuration
	config, err := GetHostConfiguration(host)
	if err != nil {
		return nil, err
	}

	// Grab the service root url
	serviceRootURL, err := GetServiceRootURLFromHostConfig(host, config)
	if err != nil {
		return nil, err
	}

	// Build URL and authorization header (root)
	url := fmt.Sprintf("%s/manage/v1/%s", serviceRootURL, path)
	authorization, err := buildRootAuthorizationHeader(host, config)
	if err != nil {
		return nil, err
	}

	return internalGet(url, authorization)
}

func ManageAPIPost(host, path string, payload map[string]any) (map[string]any, error) {

	// Lookup the host's configuration
	config, err := GetHostConfiguration(host)
	if err != nil {
		return nil, err
	}

	// Grab the service root url
	serviceRootURL, err := GetServiceRootURLFromHostConfig(host, config)
	if err != nil {
		return nil, err
	}

	// Build URL and authorization header (root)
	url := fmt.Sprintf("%s/manage/v1/%s", serviceRootURL, path)
	authorization, err := buildRootAuthorizationHeader(host, config)
	if err != nil {
		return nil, err
	}

	return internalPost(url, authorization, payload)
}

func ManageAPIDelete(host, path string) error {
	// Lookup the host's configuration
	config, err := GetHostConfiguration(host)
	if err != nil {
		return err
	}

	// Grab the service root url
	serviceRootURL, err := GetServiceRootURLFromHostConfig(host, config)
	if err != nil {
		return err
	}

	// Build URL and authorization header (root)
	url := fmt.Sprintf("%s/manage/v1/%s", serviceRootURL, path)
	authorization, err := buildRootAuthorizationHeader(host, config)
	if err != nil {
		return err
	}

	return internalDelete(url, authorization)
}

func InstanceAPIGet(host, instance, user, password, path string) (map[string]any, error) {
	// Grab the instance root url
	instanceRootURL, err := GetInstanceRootURL(host, instance)
	if err != nil {
		return nil, err
	}

	// Build URL and authorization header (user)
	url := fmt.Sprintf("%s/%s", instanceRootURL, path)
	authorization, err := buildUserAuthorizationHeader(user, password)
	if err != nil {
		return nil, err
	}

	return internalGet(url, authorization)
}

func InstanceAPIPost(host, instance, user, password, path string, payload map[string]any) (map[string]any, error) {
	// Grab the instance root url
	instanceRootURL, err := GetInstanceRootURL(host, instance)
	if err != nil {
		return nil, err
	}

	// Build URL and authorization header (user)
	url := fmt.Sprintf("%s/%s", instanceRootURL, path)
	authorization, err := buildUserAuthorizationHeader(user, password)
	if err != nil {
		return nil, err
	}

	return internalPost(url, authorization, payload)
}

func InstanceAPIDelete(host, instance, user, password, path string) error {
	// Grab the instance root url
	instanceRootURL, err := GetInstanceRootURL(host, instance)
	if err != nil {
		return err
	}

	// Build URL and authorization header (user)
	url := fmt.Sprintf("%s/%s", instanceRootURL, path)
	authorization, err := buildUserAuthorizationHeader(user, password)
	if err != nil {
		return err
	}

	return internalDelete(url, authorization)
}

func DatabaseAPIGet(host, instance, database, user, password, path string) (map[string]any, error) {
	// Grab the database root url
	databaseRootURL, err := GetDatabaseRootURL(host, instance, database)
	if err != nil {
		return nil, err
	}

	// Build URL and authorization header (user)
	url := fmt.Sprintf("%s/%s", databaseRootURL, path)
	authorization, err := buildUserAuthorizationHeader(user, password)
	if err != nil {
		return nil, err
	}

	return internalGet(url, authorization)
}

func DatabaseAPIPost(host, instance, database, user, password, path string, payload map[string]any) (map[string]any, error) {
	// Grab the database root url
	databaseRootURL, err := GetDatabaseRootURL(host, instance, database)
	if err != nil {
		return nil, err
	}

	// Build URL and authorization header (user)
	url := fmt.Sprintf("%s/%s", databaseRootURL, path)
	authorization, err := buildUserAuthorizationHeader(user, password)
	if err != nil {
		return nil, err
	}

	return internalPost(url, authorization, payload)
}

func DatabaseAPIDelete(host, instance, database, user, password, path string) error {
	// Grab the database root url
	databaseRootURL, err := GetDatabaseRootURL(host, instance, database)
	if err != nil {
		return err
	}

	// Build URL and authorization header (user)
	url := fmt.Sprintf("%s/%s", databaseRootURL, path)
	authorization, err := buildUserAuthorizationHeader(user, password)
	if err != nil {
		return err
	}

	return internalDelete(url, authorization)
}

func DatabaseAPIPutFile(host, instance, database, user, password, path, file string) error {
	// Grab the database root url
	databaseRootURL, err := GetDatabaseRootURL(host, instance, database)
	if err != nil {
		return err
	}

	// Build URL and authorization header (user)
	url := fmt.Sprintf("%s/%s", databaseRootURL, path)
	authorization, err := buildUserAuthorizationHeader(user, password)
	if err != nil {
		return err
	}

	return internalPutFile(url, authorization, file)
}
