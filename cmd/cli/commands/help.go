package commands

import (
	"log"

	common "fourclover.org/internal/common"
)

func Help() bool {
	common.INFO_BOLD(common.APP_NAME, " - ", common.APP_DESCRIPTION)
	help := `
	Version: ` + common.APP_VERSION + `
	Website: ` + common.APP_URL + `
	Usage: ` + common.APP_NAME + ` <subcommand> [options]

	Subcommands:
		snapshot	- Create a snapshot of a directory
		  ↳ ` + common.APP_NAME + ` snapshot -dir myproject -out /path/to/output/directory

		compare		- Compare two reports and show the difference
		  ↳ ` + common.APP_NAME + ` compare -old /path/to/old/report -new /path/to/new/report -out /path/to/output/directory/myreport.txt

		policy		- Check a directory against list of user defined policies
		  ↳ ` + common.APP_NAME + ` policy -targetdir /path/to/directory -policydir /path/to/policy/file -out /path/to/output/directory/myreport.txt

		help		- Show this help message
		demo		- Show a demonstration of the tool
	`
	log.Default().Println(help)
	return true
}
