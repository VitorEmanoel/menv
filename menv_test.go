package menv

import (
	"testing"
)

type Environments struct {
	Token 	string 		`name:"TOKEN" required:"true"`
	Test	string		`name:"TEST"`
}

func TestSimplesEnvironmentVariable(t *testing.T) {
	var variables = Environments{}
	err := LoadEnvironment(&variables)
	if err != nil {
		t.Error("Error in load environment", err)
	}
	if variables.Token != "TOKEN_TEST" {
		t.Errorf("Invalid values, want TOKEN_TEST variables.Token is %s", variables.Token)
	}

	if variables.Test != "test" {
		t.Errorf("Invalid values, want test variables.Token is %s", variables.Test)
	}
}
