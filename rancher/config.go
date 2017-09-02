// Copyright (c) Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package rancher

// Config for the rancher REST client
type Config struct {
	// Stack name where the faas functions get deployed
	FunctionsStackName string
	// cattle API url
	CattleURL string
	// cattle access key
	CattleAccessKey string
	// cattle secret key
	CattleSecretKey string
}

// NewClientConfig creates a new config for rancher REST client
func NewClientConfig(fn string, url string, aKey string, sKey string) (*Config, error) {
	config := Config{
		FunctionsStackName: fn,
		CattleURL:          url,
		CattleAccessKey:    aKey,
		CattleSecretKey:    sKey,
	}
	return &config, nil
}
