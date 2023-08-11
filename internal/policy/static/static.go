package policy

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"fourclover.org/internal/common"
	"fourclover.org/internal/report"
	"fourclover.org/internal/snapshot"
)

func loadPolicies(dir string) ([]string, error) {
	var policies []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			policies = append(policies, path)
		}
		return nil
	})

	return policies, err
}

func readInputFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func StaticPolicyCheck(inputTargetDir common.TargetDir, policiesDir common.PolicyDir, policyExportReportPath string) (bool, error) {
	inputTargetDirString := inputTargetDir.String()
	policiesDirString := policiesDir.String()

	policies, err := loadPolicies(policiesDirString)
	if err != nil || len(policies) == 0 {
		log.Fatalf("Failed to load policies: %v", err)
	}

	// Get list of all files from the input snapshot
	// For each file, evaluate the policy
	// If any policy fails, then return false
	// If all policies pass, then return true
	targetFilesPathList := []string{}

	hashAlgorithms := append([]string{}, "sha256")
	excludeThem := []string{}
	focusExtensions := []string{}
	snapshotName := "policycheck"
	snapshotReport, err := snapshot.SnapshotDirectory(inputTargetDirString, hashAlgorithms, excludeThem, focusExtensions, snapshotName)
	if err != nil {
		log.Fatalln("Failed to create snapshot report: ", err)
	}

	for _, file := range snapshotReport.Files {
		fileMeta := file.SHA256 + ":" + file.Path
		targetFilesPathList = append(targetFilesPathList, fileMeta)
	}

	allFindings := []report.Finding{}
	for _, targetFile := range targetFilesPathList {
		for _, policyFile := range policies {
			finding, err := SimpleStaticPolicyCheck(policyFile, targetFile)
			if err != nil {
				log.Fatalf("Failed to evaluate policy: %v", err)
			}

			allFindings = append(allFindings, finding...)
		}
	}

	if len(allFindings) > 0 {
		// Save the report to a file
		err = report.SavePolicyReport(allFindings, inputTargetDirString, policyExportReportPath)
		if err != nil {
			log.Fatalf("Failed to save report: %v", err)
		}

		log.Println("INFO: Saving report to: ", policyExportReportPath)

		return false, nil
	}

	return true, nil
}
