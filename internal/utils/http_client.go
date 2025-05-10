package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func internalGet(path, authorization string) (map[string]any, error) {
	fullURL := fmt.Sprintf("%s/%s", viper.GetString("service-root-url"), path)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func internalPost(path, authorization string, payload map[string]any) (map[string]any, error) {
	fullURL := fmt.Sprintf("%s/%s", viper.GetString("service-root-url"), path)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func internalPutFile(path, authorization, file string) error {
	fullURL := fmt.Sprintf("%s/%s", viper.GetString("service-root-url"), path)

	body, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open backupset '%s' due to: %w", file, err)
	}

	req, err := http.NewRequest(http.MethodPut, fullURL, body)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func internalDelete(path, authorization string) error {
	fullURL := fmt.Sprintf("%s/%s", viper.GetString("service-root-url"), path)

	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func ManageAPIGet(path string) (map[string]any, error) {
	manageRelativePath := fmt.Sprintf("manage/v1/%s", path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("root-client-id")+":"+viper.GetString("root-client-secret"))))

	return internalGet(manageRelativePath, authorization)
}

func ManageAPIPost(path string, payload map[string]any) (map[string]any, error) {
	manageRelativePath := fmt.Sprintf("manage/v1/%s", path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("root-client-id")+":"+viper.GetString("root-client-secret"))))

	return internalPost(manageRelativePath, authorization, payload)
}

func ManageAPIDelete(path string) error {
	manageRelativePath := fmt.Sprintf("manage/v1/%s", path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("root-client-id")+":"+viper.GetString("root-client-secret"))))

	return internalDelete(manageRelativePath, authorization)
}

func InstanceAPIGet(path string) (map[string]any, error) {
	instranceRelativePath := fmt.Sprintf("%s/api/v1/%s", viper.GetString("service-instance"), path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("user")+":"+viper.GetString("password"))))

	return internalGet(instranceRelativePath, authorization)
}

func InstanceAPIPost(path string, payload map[string]any) (map[string]any, error) {
	instranceRelativePath := fmt.Sprintf("%s/api/v1/%s", viper.GetString("service-instance"), path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("user")+":"+viper.GetString("password"))))

	return internalPost(instranceRelativePath, authorization, payload)
}

func InstanceAPIDelete(path string) error {
	instranceRelativePath := fmt.Sprintf("%s/api/v1/%s", viper.GetString("service-instance"), path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("user")+":"+viper.GetString("password"))))

	return internalDelete(instranceRelativePath, authorization)
}

func DatabaseAPIGet(path string) (map[string]any, error) {
	instranceRelativePath := fmt.Sprintf("%s/api/v1/Databases('%s')/%s", viper.GetString("service-instance"), viper.GetString("database"), path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("user")+":"+viper.GetString("password"))))

	return internalGet(instranceRelativePath, authorization)
}

func DatabaseAPIPost(path string, payload map[string]any) (map[string]any, error) {
	instranceRelativePath := fmt.Sprintf("%s/api/v1/Databases('%s')/%s", viper.GetString("service-instance"), viper.GetString("database"), path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("user")+":"+viper.GetString("password"))))

	return internalPost(instranceRelativePath, authorization, payload)
}

func DatabaseAPIDelete(path string) error {
	instranceRelativePath := fmt.Sprintf("%s/api/v1/Databases('%s')/%s", viper.GetString("service-instance"), viper.GetString("database"), path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("user")+":"+viper.GetString("password"))))

	return internalDelete(instranceRelativePath, authorization)
}

func DatabaseAPIPutFile(path, file string) error {
	instranceRelativePath := fmt.Sprintf("%s/api/v1/Databases('%s')/%s", viper.GetString("service-instance"), viper.GetString("database"), path)
	authorization := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(viper.GetString("user")+":"+viper.GetString("password"))))

	return internalPutFile(instranceRelativePath, authorization, file)
}
