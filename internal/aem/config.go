package aem

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/packr/v2"
	"github.com/jlentink/aem/internal/aem/objects"
	"github.com/jlentink/aem/internal/cli/project"
	"github.com/jlentink/aem/internal/output"
	"os"
)

// Cnf active configuration
var Cnf *objects.Config

const (
	instanceMainDefault = `local-author`
	instanceEnv         = `AEM_ME`
)

// ConfigExists is there a config available
func ConfigExists() bool {
	path, _ := project.GetConfigFileLocation()
	return project.Exists(path)
}

// Render 's the template to a string
func Render() string {
	box := packr.New("templates", "../../_templates/")
	b, err:= box.Find("aem.toml")

	if err != nil {
		return ""
	}

	return string(b)
}

// WriteConfigFile to disk
func WriteConfigFile() (int, error) {
	p, err := project.GetConfigFileLocation()
	if err != nil {
		return 0, err
	}

	p = p[0:len(p)-5] + ".example.toml"
	return project.WriteTextFile(p, Render())
}

// GetConfig Read config page
func GetConfig() (*objects.Config, error) {
	p, err := project.GetConfigFileLocation()
	if err != nil {
		return nil, fmt.Errorf("could not find config file")
	}

	cnf := objects.Config{}
	_, err = toml.DecodeFile(p, &cnf)
	if err != nil {
		return nil, fmt.Errorf("could not decode config file: %s", err.Error())
	}

	Cnf = &cnf
	objects.Cnf = &cnf

	if cnf.Schema != objects.SchemaVersion {
		fmt.Println("Your toml schema does not match the this version.")
		fmt.Printf("It should be \"%s\" but it is \"%s\". Please update.\n", objects.SchemaVersion, cnf.Schema)
		os.Exit(2)
	}

	return &cnf, nil
}

// GetDefaultInstanceName Instance based on resolution order
func GetDefaultInstanceName() string {
	envName := os.Getenv(instanceEnv)
	if len(envName) > 0 {
		return envName
	}

	c, err := GetConfig()
	if err != nil {
		output.Printf(output.VERBOSE, "Error in  config returning default author")
		return instanceMainDefault
	}

	if len(c.DefaultInstance) > 0 {
		return c.DefaultInstance
	}

	return instanceMainDefault
}
