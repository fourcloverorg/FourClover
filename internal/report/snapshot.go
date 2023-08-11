package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	common "fourclover.org/internal/common"
)

type SnapshotReport struct {
	FourCloverVersion string        `json:"fourclover_version"` // fourclover version
	Name              string        `json:"name"`               // SnapshotReport name
	Date              string        `json:"date"`
	WorkingDir        string        `json:"working_dir"`      // Current working directory
	TotalFiles        int           `json:"total_files"`      // Total files in the directory
	TotalFilesSize    int64         `json:"total_files_size"` // Total file sizes in bytes
	ReportChecksum    string        `json:"report_checksum"`  // SHA256 checksum of the report
	Files             []common.File `json:"files"`
}

func SaveSnapshotReport(report SnapshotReport, filename string) error {
	jsonData, err := json.Marshal(report)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func reportCompareOperation(changes map[string]string) (workingDirs []string, addedFiles []string, modifiedFiles []string, deletedFiles []string, duplicatedFiles []string, renamedFiles []string, unchangedFiles []string, unique_changes map[string]string) {
	workingDirs = []string{}
	addedFiles = []string{}
	modifiedFiles = []string{}
	deletedFiles = []string{}
	duplicatedFiles = []string{}
	renamedFiles = []string{}
	unchangedFiles = []string{}
	for file, change := range changes {
		switch change {
		case "workingDirs":
			workingDirs = append(workingDirs, file)
		case "added":
			addedFiles = append(addedFiles, file)
		case "modified":
			modifiedFiles = append(modifiedFiles, file)
		case "deleted":
			deletedFiles = append(deletedFiles, file)
		case "duplicated":
			duplicatedFiles = append(duplicatedFiles, file)
		case "renamed":
			renamedFiles = append(renamedFiles, file)
		case "unchanged":
			unchangedFiles = append(unchangedFiles, file)
		default:
			common.INFO("Unknown change type", change)
		}
	}

	return workingDirs, addedFiles, modifiedFiles, deletedFiles, duplicatedFiles, renamedFiles, unchangedFiles, changes
}

// Save the comparison report to a txt file
func writeCompareReport(compareReportFileName string, workingDirs []string, addedFiles []string, modifiedFiles []string, deletedFiles []string, duplicatedFiles []string, renamedFiles []string, unchangedFiles []string, changes map[string]string) error {
	// Create the report file with the current date and time
	reportFile, err := os.Create(compareReportFileName)
	if err != nil {
		return fmt.Errorf("error creating report file: %w", err)
	}
	defer reportFile.Close()

	// Writing the report in the following format
	/*
		-------------------- CHANGE REPORT --------------------
		Date:                           2023 Feb 02 14:41 UTC
		Scanned directory:              Old report: .
		                                New report: .
		Total number of changes:        26
		Number of files added:          21
		Number of files altered:        1
		Number of files deleted:        1
		Number of files renamed:        2
		-------------------------------------------------------
		[ADDED] 0002_auto_20210802_1243.py
		[ADDED] urls.cpython-38.pyc
		[ALTERED] .gitignore
		[RENAMED] README.md to PlZREADME.md
		[DELETED] LICENSE
		--------------------- REPORT END ----------------------
	*/

	// Remove "unchanged" from the changes map
	for _, file := range unchangedFiles {
		delete(changes, file)
	}

	// Remove "duplicate" from the changes map
	for _, file := range duplicatedFiles {
		delete(changes, file)
	}

	// Remove "workingDirs" from the changes map
	for _, file := range workingDirs {
		delete(changes, file)
	}

	// Write the report header
	reportFile.WriteString("-------------------- CHANGE REPORT --------------------\n")
	reportFile.WriteString("Date:                   \t" + time.Now().Format("2006 Feb 01 15:04 UTC") + "\n")
	reportFile.WriteString("Scanned directory:      \t" + workingDirs[0] + "\n")
	reportFile.WriteString("Total number of changes:\t" + strconv.Itoa(len(changes)) + "\n")
	reportFile.WriteString("Number of files added:  \t" + strconv.Itoa(len(addedFiles)) + "\n")
	reportFile.WriteString("Number of files altered:\t" + strconv.Itoa(len(modifiedFiles)) + "\n")
	reportFile.WriteString("Number of files deleted:\t" + strconv.Itoa(len(deletedFiles)) + "\n")
	reportFile.WriteString("Number of files renamed:\t" + strconv.Itoa(len(renamedFiles)) + "\n")
	reportFile.WriteString("Number of files unchanged:\t" + strconv.Itoa(len(unchangedFiles)) + "\n")
	reportFile.WriteString("-------------------------------------------------------\n")

	// Write the report body
	if len(addedFiles) > 0 {
		for _, file := range addedFiles {
			reportFile.WriteString("[ADDED] " + file + "\n")
		}
	}

	if len(modifiedFiles) > 0 {
		for _, file := range modifiedFiles {
			reportFile.WriteString("[ALTERED] " + file + "\n")
		}
	}

	if len(renamedFiles) > 0 {
		for _, file := range renamedFiles {
			reportFile.WriteString("[RENAMED] " + file + "\n")
		}
	}

	if len(deletedFiles) > 0 {
		for _, file := range deletedFiles {
			reportFile.WriteString("[DELETED] " + file + "\n")
		}
	}

	if len(addedFiles) == 0 && len(modifiedFiles) == 0 && len(deletedFiles) == 0 && len(renamedFiles) == 0 {
		reportFile.WriteString("No changes observed!\n")
	}

	// Write the report footer
	reportFile.WriteString("--------------------- REPORT END ----------------------")
	// Tool Credits for the report.
	credits := "\n\nReport generated using the fourclover v" + common.APP_VERSION
	reportFile.WriteString(credits)
	// Inform the user that the report has been saved.
	log.Default().Println("INFO: This report has been saved to: " + compareReportFileName)

	return nil
}

func SaveCompareReport(compareReportFileName string, changes map[string]string) (int, error) {
	workingDirs, addedFiles, modifiedFiles, deletedFiles, duplicatedFiles, renamedFiles, unchangedFiles, changes := reportCompareOperation(changes)
	totalFilesObserved := len(addedFiles) + len(modifiedFiles) + len(deletedFiles) + len(renamedFiles) + len(unchangedFiles)

	// printReport(workingDirs, addedFiles, modifiedFiles, deletedFiles, duplicatedFiles, renamedFiles, unchangedFiles, changes)
	err := writeCompareReport(compareReportFileName, workingDirs, addedFiles, modifiedFiles, deletedFiles, duplicatedFiles, renamedFiles, unchangedFiles, changes)
	if err != nil {
		return 0, fmt.Errorf("error saving report: %w", err)
	}
	return totalFilesObserved, nil
}
