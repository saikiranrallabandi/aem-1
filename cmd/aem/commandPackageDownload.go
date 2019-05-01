package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
)

type commandPackageDownload struct {
	From             string
	To               string
	Type             string
	Role             string
	Name             string
	Package          string
	utility          *utility
	projectStructure projectStructure
	http             *httpRequests
	forceDownload    bool
}

func (c *commandPackageDownload) Init() {
	c.From = configDefaultInstance
	c.utility = new(utility)
	c.projectStructure = newProjectStructure()
	c.forceDownload = false
	c.Package = ""
	c.http = new(httpRequests)
}

func (c *commandPackageDownload) readConfig() bool {
	return true
}

func (c *commandPackageDownload) GetCommand() []string {
	return []string{"package-download"}
}

func (c *commandPackageDownload) GetHelp() string {
	return "Download package from AEM instance."
}

func (c *commandPackageDownload) Execute(args []string) {
	c.getOpt(args)

	fromInstance := c.utility.getInstanceByName(c.From)

	if len(c.Package) > 0 {
		packages := c.utility.pkgsFromString(fromInstance, c.Package)
		c.downloadPackages(fromInstance, packages)
	} else {
		pkgPicker := newPackagePicker()
		packages := pkgPicker.picker(fromInstance)
		c.downloadPackages(fromInstance, packages)
	}
}

func (c *commandPackageDownload) downloadPackages(instance aemInstanceConfig, pkgs []packageDescription) {
	for _, pkg := range pkgs {
		fmt.Printf("Download: %s\n", pkg.Name)
		c.http.downloadPackage(instance, pkg, c.forceDownload)
	}
}

func (c *commandPackageDownload) getOpt(args []string) {
	getopt.FlagLong(&c.From, "from", 'f', "Pull content from (default: "+configDefaultInstance+")")
	getopt.FlagLong(&c.Package, "package", 'p', "Define package package:version (no interactive mode)")
	getopt.FlagLong(&c.forceDownload, "force-download", 'd', "Force new download")
	getopt.CommandLine.Parse(args)
}
