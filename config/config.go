package config

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"io/ioutil"
)

// Struct for the global config
type Config struct {
	Interval int            `hcl:"interval"`
	Prefix   string         `hcl:"prefix"`
	Statsd   string         `hcl:"statsd"`
	Groups   []TargetGroups `hcl:"target_group"`
}

// nested struct for the groups
type TargetGroups struct {
	Name     string    `hcl:",key"`
	Prefix   string    `hcl:"prefix"`
	Interval int       `hcl:"interval"`
	Targets  []Targets `hcl:"target"`
}

// final tested struct for the targets in each group
type Targets struct {
	Address string `hcl:"address"`
	Label   string `hcl:",key"`
}

// parse the config
func Parse(ConfigFile string) (*Config, error) {

	result := &Config{}
	var errors *multierror.Error

	// read the config file
	config, err := ioutil.ReadFile(ConfigFile)

	// if there's an issue, bomb
	if err != nil {
		return nil, err
	}

	// parse the config tree
	hclParseTree, err := hcl.Parse(string(config))
	if err != nil {
		return nil, err
	}

	// decode each object in the tree
	if err := hcl.DecodeObject(&result, hclParseTree); err != nil {
		return nil, err
	}

	return result, errors.ErrorOrNil()

}
