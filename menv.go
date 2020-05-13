package menv

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

// LoadEnvironment load all environments
func LoadEnvironment(profiles ...interface{}) error {
	for _, profile := range profiles {
		err := NewEnvironmentLoader().Parse(profile).Load()
		if err != nil {
			return err
		}
	}
	return nil
}

// EnvironmentLoader interface to manager environment variables
type EnvironmentLoader interface {
	Parse(profile interface{}) EnvironmentLoader
	Load() error
}

// LoaderContext is struct to save loader data
type LoaderContext struct {
	EnvironmentFile string
	Environments    []Environment
}

// Environment represent environment variables
type Environment struct {
	Name       string
	Required   bool
	FieldValue reflect.Value
}

// Parse load all profile data
func (l *LoaderContext) Parse(profile interface{}) EnvironmentLoader {
	var profileType = reflect.TypeOf(profile)
	if profileType.Kind() != reflect.Ptr {
		return l
	}
	if profileType.Kind() == reflect.Ptr {
		profileType = profileType.Elem()
	}
	var profileValue = reflect.ValueOf(profile)
	if profileValue.Kind() == reflect.Ptr {
		profileValue = profileValue.Elem()
	}
	for i := 0; i < profileType.NumField(); i++ {
		var field = profileType.Field(i)
		if field.Type.Kind() != reflect.String {
			continue
		}
		if !profileValue.Field(i).CanSet() {
			continue
		}
		var name = field.Tag.Get("name")
		if name == "" {
			continue
		}
		var localEnv = Environment{}
		localEnv.Name = name
		if envFile, ok := field.Tag.Lookup("file"); ok && envFile != "" {
			l.EnvironmentFile = envFile
		}
		if required, ok := field.Tag.Lookup("required"); ok && required == "true" {
			localEnv.Required = true
		}
		localEnv.FieldValue = profileValue.Field(i)
		l.Environments = append(l.Environments, localEnv)
	}
	return l
}

func (l *LoaderContext) loadEnvironmentFile() error {
	content, err := ioutil.ReadFile(l.EnvironmentFile)
	if err != nil {
		return nil
	}
	var contentStr = string(content)
	contentStr = strings.ReplaceAll(contentStr, "\t", "")
	contentStr = strings.ReplaceAll(contentStr, "\r", "")
	var contents = strings.Split(contentStr, "\n")
	for _, item := range contents {
		var splitedItem = strings.Split(item, "=")
		if len(splitedItem) != 2 {
			return errors.New("invalid content")
		}
		var value = splitedItem[1]
		value = strings.ReplaceAll(value, "\"", "")
		value = strings.ReplaceAll(value, "'", "")
		os.Setenv(splitedItem[0], value)
	}
	return nil
}

// Load load profile
func (l *LoaderContext) Load() error {
	err := l.loadEnvironmentFile()
	if err != nil {
		return err
	}
	for _, environment := range l.Environments {
		var environmentValue = os.Getenv(environment.Name)
		if environmentValue == "" && environment.Required {
			return fmt.Errorf("required variable %s not found", environment.Name)
		}
		environment.FieldValue.SetString(environmentValue)
	}
	return nil
}

// NewEnvironmentLoader constructor of environmentLoader
func NewEnvironmentLoader() EnvironmentLoader {
	return &LoaderContext{EnvironmentFile: ".env"}

}
