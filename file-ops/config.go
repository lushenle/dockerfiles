package config

import (
	"io/ioutil"
	"os"

	"github.com/lunny/log"
	yaml "gopkg.in/yaml.v2"
)

// ConfigYaml resources metadata config
type ConfigYaml struct {
	Region    string `yaml:"region,omitempty" json:"region,omitempty"`
	EC2Config struct {
		TagName                  string `yaml:"tag_name" json:"tag_name"`
		TagValue                 string `yaml:"tag_value" json:"tag_value"`
		AmiID                    string `yaml:"ami_id" json:"ami_id"`
		InstanceType             string `yaml:"instance_type" json:"instance_type"`
		KeyPairName              string `yaml:"key_pair_name" json:"key_pair_name"`
		SecurityGroupName        string `yaml:"security_group_name" json:"security_group_name"`
		SecurityGroupDescription string `yaml:"security_group_description" json:"security_group_description"`
	}
	S3Config struct {
		BucketName string `yaml:"bucket_name" json:"bucket_name"`
		FileName   string `yaml:"file_name" json:"file_name"`
	}
	SQSConfig struct {
		SQSName string `yaml:"sqs_name" json:"sqs_name"`
	}
}

var yamlConfig *ConfigYaml = nil

// GetYamlConfig return config object
func GetYamlConfig() *ConfigYaml {
	return yamlConfig
}

// ParseYamlConfig read config from file and parse
func ParseYamlConfig(config string) *ConfigYaml {
	f, err := os.Open(config)
	if err != nil {
		log.Errorf("open config file err: %v\n", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)

	err = yaml.Unmarshal([]byte(data), &yamlConfig)
	if err != nil {
		log.Errorf("unmarshal config err: %v\n", err)
	}

	return yamlConfig
}
