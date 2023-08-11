package commands

import (
	"flag"
	"log"
	"os"
	"time"

	common "fourclover.org/internal/common"
	config "fourclover.org/internal/config"
	policyStatic "fourclover.org/internal/policy/static"
)

var (
	policyCommand  *flag.FlagSet
	inputPolicyDir common.PolicyDir
	inputTargetDir common.TargetDir

	policyExportReportPath string
)

func Policy(configData []byte) (bool, error) {
	timeNow := time.Now().Format("2006-01-02_15-04-05")
	if configData == nil {
		log.Default().Println("INFO: Reading policy CLI arguments.")

		policyCommand = flag.NewFlagSet("policy", flag.ExitOnError)

		// eg: fourclover policy -targetdir [target directory]
		policyCommand.Var(&inputTargetDir, "targetdir", "Path of the directory containing the source code to scan. (eg: -targetdir [target directory])")

		// eg: fourclover policy -i snap.json -p policydir
		policyCommand.Var(&inputPolicyDir, "policydir", "[Optional] Path of the directory containing the policy files. (eg: -policydir mypolicydir) (default: policydir)")

		// eg: fourclover policy -i snap.json -p policydir -o reportfile.txt
		policyCommand.StringVar(&policyExportReportPath, "out", "policy_report_"+timeNow+".json", "[Optional] Export report file")

		// eg: fourclover policy -i snap.json -p policydir -o outputdir -s true
		policyCommand.BoolVar(&supressScanOutput, "supress", false, "[Optional] Supress the output of the program. (default: false)")

		policyCommand.Parse(os.Args[2:])

		if len(os.Args) == 2 {
			log.Default().Println("Usage of policy check:")
			policyCommand.PrintDefaults()
			return false, nil
		}

		if inputTargetDir == "" {
			log.Fatalf("ERROR: snapshot: missing required flag: -input. The missing flag is required for the snapshot file to use for the policy scan.")
		}

		if inputPolicyDir == "" {
			log.Fatalf("ERROR: snapshot: missing required flag: -policydir. The missing flag is required for the directory containing the policy files.")
		}
	}

	// If the config file has data in it, then use that data to set the values for the variables. If not, then use the default values.
	if configData != nil {
		log.Default().Println("INFO: Reading policy configuration file.")

		fourcloverTargetDirPath := config.GetFourCloverPolicyTargetDirPath()
		if fourcloverTargetDirPath != "" {
			log.Default().Println("INFO: Target directory path found in the configuration file. -> " + fourcloverTargetDirPath)
			inputTargetDir = common.TargetDir(fourcloverTargetDirPath)
		} else {
			log.Fatalf("ERROR: Target directory path not found in the configuration file.")
		}

		fourcloverPolicyDirPath := config.GetFourCloverPolicyDirPath()
		if fourcloverPolicyDirPath != "" {
			log.Default().Println("INFO: Policy directory path found in the configuration file. -> " + fourcloverPolicyDirPath)
			inputPolicyDir = common.PolicyDir(fourcloverPolicyDirPath)
		} else {
			log.Fatalf("ERROR: Policy directory path not found in the configuration file.")
		}

		fourcloverPolicyReportFilePath := config.GetFourCloverPolicyExportReportPath()
		if fourcloverPolicyReportFilePath != "" {
			log.Default().Println("INFO: Export SnapshotReport file path found in the configuration file. -> " + fourcloverPolicyReportFilePath)
			policyExportReportPath = fourcloverPolicyReportFilePath
		} else {
			log.Fatalf("ERROR: Export SnapshotReport file path not found in the configuration file.")
		}
	}

	// Actual program startes here after all the checks
	log.Default().Println("INFO: Starting policy check.")

	result, err := policyStatic.StaticPolicyCheck(inputTargetDir, inputPolicyDir, policyExportReportPath)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	return result, nil
}
