package commands

import (
	"fmt"
	"log"

	common "fourclover.org/internal/common"
)

func Demo() bool {
	log.Default().Println("INFO: FourClover: demo page loaded")
	common.INFO_BOLD(common.APP_NAME, " - ", common.APP_DESCRIPTION)
	common.WHITE("INFO:", common.APP_NAME, "Version:", common.APP_VERSION)

	// Sample usage of scan
	common.INFO_BOLD("Sample usage of scan operation:")
	common.WHITE(`$ ` + common.APP_NAME + ` snapshot myproject -out snapshot_report.json -hashes sha256,sha1,md5 -exclude .git,.idea,node_modules -focus .py,.js,.txt -name "My Project"`)
	common.YELLOW("	↳ Note: The -hash supported algorithms: blake2b-256, blake2b-512, blake3, crc32, md5, sha1, sha3-224, sha3-256, sha3-384, sha3-512, sha256, sha512.")
	common.YELLOW("	↳ Note: The -exclude option is used to exclude files and directories from the scan. It takes a comma separated list of files and directory names.")
	common.YELLOW("	↳ Note: The -focus option is used to scan only files with the specified extensions. It takes a comma separated list of file extensions.")
	common.YELLOW("	↳ Note: The -name option is used to give a name to the report. This is useful when you want to compare reports with different names.")

	// Sample scan report
	common.INFO_BOLD("\n" + "Sample scan report:")
	sampleScanReport := `  {
	  "fourclover_version": "0.1.0",
	  "name": "test report",
	  "date": "2023-08-06T16:58:45+05:30",
	  "working_dir": "D:/myproject",
	  "total_files": 38,
	  "total_files_size": 70657,
	  "report_checksum": "VaAwRVR039V2ujn-urAH7CQh....PItufG-3GsmrTfU89weA=",
	  "files": [
            {
                "name": "db.sqlite3",
                "file_path": "myproject/dpm/1db.sqlite3",
                "size": 10244,
                "permission": "-rwxrwxrwx",
                "last_modified": "2023-02-03T17:15:58+05:30",
                "sha256": "c023acba588067e4c3b470f2f043528f4e63a47d2605a9d8837019cb4eed7564"
            }
            ...
       ]
    }`
	common.WHITE(sampleScanReport)

	// Sample usage of compare
	common.INFO_BOLD("\n" + "Sample usage of compare operation:")
	common.WHITE("$", common.APP_NAME, "compare -new snapshot_report_1.json -old snapshot_report_2.json -out comparison_report_name.txt")
	common.YELLOW("	↳ Note: The -new option is used to specify the new report.")
	common.YELLOW("	↳ Note: The -old option is used to specify the old report.")
	common.YELLOW("	↳ Note: The -out option is used to specify the output file name.")

	// Sample compare report
	common.INFO_BOLD("\n" + "Sample comparison report:")
	sampleCompareReport := `  -------------------- CHANGE REPORT --------------------
  Date:                   	2023 Feb 02 17:22 UTC
  Scanned directory:      	Old report: /home/loc/myproject
 				New report: /home/loc/myproject
  Total number of changes:	4
  Number of files added:  	1
  Number of files altered:	1
  Number of files deleted:	1
  Number of files renamed:	1
  -------------------------------------------------------
  [ADDED] admin.cpython-38.pyc
  [ALTERED] .gitignore
  [RENAMED] 1db.sqlite3 -> db.sqlite3
  [DELETED] LICENSE
  --------------------- REPORT END ----------------------`
	fmt.Printf("%s", sampleCompareReport)

	// Sample usage of policy
	common.INFO_BOLD("\n" + "Sample usage of policy operation:")
	common.WHITE("$", common.APP_NAME, "policy -targetdir mydirectory -policydir policyfolder -out /path/to/output/directory/myreport.json")
	common.YELLOW("	↳ Note: The -targetdir option is used to specify the directory to be scanned.")
	common.YELLOW("	↳ Note: The -policydir option is used to specify the directory containing the policy files.")
	common.YELLOW("	↳ Note: The -out option is used to specify the output directory and filename.")

	// Sample policy report
	common.INFO_BOLD("\n" + "Sample policy report:")
	samplePolicyReport := `  {
	  "title": "Policy Report",
	  "date": "2023-08-06T06:58:23+05:30",
	  "working_dir": "mydirectory",
	  "summary": "Findings reported in the policy report are based on the policy rules defined in the policy file.",
	  "findings": [
	  	{
		    "rule_id": "word-hunter",
		    "name": "Word Hunter",
			"type": "simple",
			"description": "The code contains a word that is not allowed.",
			"file_hash": "de5bc28dde5feb543d3992d0fac49ba2fc28f2b028f76a315f9e0e9d9985fdb8",
			"file_path": "ibm/Dockerfile",
			"references": [
				"https://www.example.com"
			],
			"classification": [
				{
					"cvss_score": 0,
					"cvss_vector": "NA"
				}
			],
			"tags": [
				"word-hunter",
				"simple"
			]
		  },
		  ...
	  ]
  }`

	common.WHITE(samplePolicyReport)

	// End of demo
	common.INFO_BOLD("\n" + "End of demo")

	return true
}
