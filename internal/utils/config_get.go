package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func SaveConfiguration() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to update configuration: %v", err)
	}
	return nil
}

func GetHostName(name string) (string, error) {

	// No host specified then use active host
	if name == "" {
		name = viper.GetString("host")
	}

	// No (active) host specified then return an error
	if name == "" {
		return "", fmt.Errorf("no host specified")
	}

	return name, nil
}

func GetHostConfiguration(name string) (map[string]any, error) {

	// Get the host name
	name, err := GetHostName(name)
	if err != nil {
		return nil, err
	}

	// Lookup the host in list of configured hosts
	hosts := viper.GetStringMap("hosts")
	raw := hosts[name]
	if raw == nil {
		return nil, fmt.Errorf("no configuration specified for host '%s'", name)
	}
	host, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for host '%s', format invalid", name)
	}
	return host, nil
}

func GetUserName(name string) (string, error) {

	// No user specified then use active user
	if name == "" {
		name = viper.GetString("user")
	}

	// No (active) user specified then return an error
	if name == "" {
		return "", fmt.Errorf("no user specified")
	}

	return name, nil
}

func GetUserConfiguration(name string) (map[string]any, error) {

	// Get the user name
	name, err := GetUserName(name)
	if err != nil {
		return nil, err
	}

	// Lookup the user in list of configured users
	users := viper.GetStringMap("users")
	raw := users[name]
	if raw == nil {
		return nil, fmt.Errorf("no configuration specified for user '%s'", name)
	}
	user, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for user '%s', format invalid", name)
	}
	return user, nil
}

func getStringFromHostConfig(name string, config map[string]any, prop_name, prop_desc string) (string, error) {

	// Lookup the property in the configuration of the host
	raw := config[prop_name]
	var value string
	if raw != nil {
		cast, ok := raw.(string)
		if !ok {
			return "", fmt.Errorf("invalid %s format for host '%s'", prop_name, name)
		}
		value = cast
	}
	if raw == nil || value == "" {
		return "", fmt.Errorf("invalid configuration, no %s specified for host '%s'", prop_name, name)
	}
	return value, nil
}

func GetServiceRootURLFromHostConfig(name string, config map[string]any) (string, error) {

	// Lookup the service root url property in the configuration of the host
	return getStringFromHostConfig(name, config, "service_root_url", "service root URL")
}

func GetRootClientIDFromHostConfig(name string, config map[string]any) (string, error) {

	// Lookup the root client id property in the configuration of the host
	return getStringFromHostConfig(name, config, "root_client_id", "root client id")
}

func GetRootClientSecretFromHostConfig(name string, config map[string]any) (string, error) {

	// Lookup the root client secret property in the configuration of the host
	return getStringFromHostConfig(name, config, "root_client_secret", "root client secret")
}

func GetServiceRootURL(host string) (string, error) {

	// Lookup the host's configuration
	config, err := GetHostConfiguration(host)
	if err != nil {
		return "", err
	}

	// Return the service root URL from the host's configuration
	return GetServiceRootURLFromHostConfig(host, config)
}

func GetInstanceNameFromHostConfig(instance string, config map[string]any) (string, error) {

	// No instance specified then use active instance
	if instance == "" {
		raw := config["instance"]
		if raw != nil {
			cast, ok := raw.(string)
			if !ok {
				return "", fmt.Errorf("invalid configuration, 'instance' property is not a string")
			}
			instance = cast
		}
		if raw == nil || instance == "" {
			return "", fmt.Errorf("no instance specified")
		}
	}
	return instance, nil
}

func GetInstanceName(host, instance string) (string, error) {

	// Lookup the host's configuration
	config, err := GetHostConfiguration(host)
	if err != nil {
		return "", err
	}

	// No instance specified then use active instance
	return GetInstanceNameFromHostConfig(instance, config)
}

func GetInstanceRootURL(host, instance string) (string, error) {

	// Lookup the host's configuration
	config, err := GetHostConfiguration(host)
	if err != nil {
		return "", err
	}

	// No instance specified then use active instance
	instance, err = GetInstanceNameFromHostConfig(instance, config)
	if err != nil {
		return "", err
	}

	// Return the service root URL from the host's configuration
	serviceRootURL, err := GetServiceRootURLFromHostConfig(host, config)
	if err != nil {
		return "", err
	}

	// Return the instance service root URL
	return fmt.Sprintf("%s/%s/api/v1", serviceRootURL, instance), nil
}

func GetDatabaseRootURL(host, instance, database string) (string, error) {

	// Grab the instance root URL
	instanceRootURL, err := GetInstanceRootURL(host, instance)
	if err != nil {
		return "", err
	}

	// Make sure a database name is specified
	if database == "" {
		return "", fmt.Errorf("no database specified")
	}

	// Return the database root URL
	return fmt.Sprintf("%s/Databases('%s')", instanceRootURL, database), nil
}
