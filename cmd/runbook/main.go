package main

import (
	"os"

	"github.com/Financial-Times/cmdb-go/cmdb"
	"github.com/Sirupsen/logrus"
	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/jawher/mow.cli"
	_ "github.com/joho/godotenv/autoload"
	yaml "gopkg.in/yaml.v2"
)

func main() {

	app := cli.App("runbook", "Automatically update cmdb runbook")

	globalSystemFile := app.String(cli.StringOpt{
		Name:   "globalSystemFile",
		Value:  "cmdb/cmdb-global.yaml",
		Desc:   "Global default values for system information for CMDB.",
		EnvVar: "GLOBAL_SYSTEM_FILE",
	})

	systemFile := app.String(cli.StringOpt{
		Name:   "systemFile",
		Value:  "cmdb/cmdb.yaml",
		Desc:   "Place to look for system information for CMDB.",
		EnvVar: "SYSTEM_FILE",
	})

	cmdbEndpoint := app.String(cli.StringOpt{
		Name:   "cmdbEndpoint",
		Value:  "https://cmdb.ft.com",
		Desc:   "CMDB Endpoint (test/prod etc.)",
		EnvVar: "CMDB_ENDPOINT",
	})

	apiKey := app.String(cli.StringOpt{
		Name:   "apiKey",
		Desc:   "API Key to access CMDB",
		EnvVar: "CMDB_API_KEY",
	})

	app.Action = func() {

		logrus.Infof("Connecting to %s", *cmdbEndpoint)

		globalSystem := loadSystem(*globalSystemFile)
		system := loadSystem(*systemFile)

		finalSystem := combineSystems(globalSystem, system)
		logrus.Debugf("Final system: %v", finalSystem)

		cmdb, err := cmdb.NewClient(*cmdbEndpoint, *apiKey)
		if err != nil {
			panic(err)
		}

		cmdb.UpdateSystem(finalSystem)

		logrus.Infof("Updated System %s", system.SystemCode)

	}

	app.Run(os.Args)
}

func loadSystem(f string) cmdb.System {
	logrus.Infof("Loading file %s", f)
	file, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	system := cmdb.System{}
	err = yaml.Unmarshal(file, &system)
	if err != nil {
		panic(err)
	}

	system.Troubleshooting = system.ImportMarkdown(system.Troubleshooting, f)

	logrus.Debug(system)
	return system
}

func combineSystems(defaultSystem, overrideSystem cmdb.System) cmdb.System {
	updateStringConditionly(&defaultSystem.SystemCode, overrideSystem.SystemCode)
	updateStringConditionly(&defaultSystem.Name, overrideSystem.Name)
	updateStringConditionly(&defaultSystem.Description, overrideSystem.Description)
	updateStringConditionly(&defaultSystem.ServiceTier, overrideSystem.ServiceTier)
	updateStringConditionly(&defaultSystem.LifecycleStage, overrideSystem.LifecycleStage)
	updateStringConditionly(&defaultSystem.ArchitectureDiagram, overrideSystem.ArchitectureDiagram)
	updateStringConditionly(&defaultSystem.Troubleshooting, overrideSystem.Troubleshooting)
	updateLinksConditionally(&defaultSystem.MoreInformation, overrideSystem.MoreInformation)
	updateLinksConditionally(&defaultSystem.Monitoring, overrideSystem.Monitoring)
	updateStringConditionly(&defaultSystem.GitRepository, overrideSystem.GitRepository)
	updateStringConditionly(&defaultSystem.HostPlatform, overrideSystem.HostPlatform)
	updateContactConditionally(&defaultSystem.PrimaryContact, overrideSystem.PrimaryContact)
	updateContactConditionally(&defaultSystem.SecondaryContact, overrideSystem.SecondaryContact)
	updateContactConditionally(&defaultSystem.Programme, overrideSystem.Programme)
	updateContactConditionally(&defaultSystem.ProductOwner, overrideSystem.ProductOwner)
	updateContactConditionally(&defaultSystem.TechnicalLead, overrideSystem.TechnicalLead)

	return defaultSystem
}

func updateStringConditionly(defaultField *string, overrideField string) {
	if overrideField != "" {
		*defaultField = overrideField
	}
}

func updateContactConditionally(defaultContact *cmdb.Contact, overrideContact cmdb.Contact) {
	if len(overrideContact.Entries) > 0 {
		*defaultContact = overrideContact
	}
}

func updateLinksConditionally(defaultLink *[]cmdb.Link, overrideLink []cmdb.Link) {
	if len(overrideLink) > 0 {
		*defaultLink = overrideLink
	}
}
