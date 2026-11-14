package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type ParseOptions struct {
	IncludeDefaultSection bool
}

const (
	SSHKeyConfigKey       = "around_ssh_key"
	INIDefaultSectionName = "DEFAULT"
)

func CollectConfigurations(configPath string) []string {
	configurationPaths := []string{}
	info, err := os.Stat(configPath)
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		directory, err := os.ReadDir(configPath)
		if err != nil {
			panic(err)
		}
		for _, file := range directory {
			if file.Type().IsRegular() {
				configurationPaths = append(configurationPaths, filepath.Join(configPath, file.Name()))
			}
		}
	} else {
		configurationPaths = append(configurationPaths, configPath)
	}
	return configurationPaths
}

func ParseFile(fileName string, opt ParseOptions) ([]Configuration, error) {
	result := []Configuration{}
	iniFile, err := ini.ShadowLoad(fileName)
	if err != nil {
		return nil, err
	}
	for _, section := range iniFile.Sections() {
		if section.Name() == INIDefaultSectionName && !opt.IncludeDefaultSection {
			continue
		}
		conf := Configuration{
			Source: fileName,
		}
		conf.Name = section.Name()
		conf.SSHKeys = section.Key(SSHKeyConfigKey).ValueWithShadows()
		section.DeleteKey(SSHKeyConfigKey)
		for _, field := range section.Keys() {
			conf.Variables = append(conf.Variables, Variable{
				Name:  field.Name(),
				Value: field.Value(),
			})
		}
		result = append(result, conf)
	}
	return result, nil
}

func Validate(configuration []Configuration) error {
	for _, conf := range configuration {
		for _, reserved := range []string{"unset", "start_ssh_agent", "stop_ssh_agent"} {
			if conf.Name == reserved {
				return fmt.Errorf("invalid section name %s in file %s", conf.Name, conf.Source)
			}
		}
	}
	return nil
}
