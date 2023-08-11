package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	common "fourclover.org/internal/common"
	report "fourclover.org/internal/report"
)

func LoadReport(filename string) (report.SnapshotReport, error) {
	var report report.SnapshotReport

	file, err := os.Open(filename)
	if err != nil {
		return report, err
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		return report, err
	}

	err = json.Unmarshal(jsonData, &report)
	if err != nil {
		return report, err
	}

	return report, nil
}

// Check if the report_checksum  matches the checksum of both reports matches the checksum of the files: [] in the report
// If not, the report is corrupted
func checkReportChecksum(reportType string, report report.SnapshotReport) (bool, error) {
	var genReportChecksum string

	log.Default().Println("INFO: Checking if the", reportType, "report is altered.")
	// Calculate the checksum of the report using the file names, paths, permissions, last modified dates, SHA256 hashes, date and fourclover version
	reportDate := report.Date
	reportfourcloverVersion := report.FourCloverVersion
	for _, file := range report.Files {
		genReportChecksum += file.Name + file.Path + file.Permission + file.LastModified + file.SHA256 + reportDate + reportfourcloverVersion
	}
	calculatedReportHash := sha256.Sum256([]byte(genReportChecksum))
	calculatedReportChecksum := hex.EncodeToString(calculatedReportHash[:])

	// Decrypt the checksum of the report using the cipher key
	DecryptedReportChecksum, err := common.Decrypt(common.CIPHER_KEY, report.ReportChecksum)
	if err != nil {
		return false, err
	}

	// Compare the calculated checksum of the report with the decrypted checksum of the report. If they are not equal, the report is altered
	if calculatedReportChecksum != DecryptedReportChecksum {
		return false, fmt.Errorf(reportType, "report is altered")
	}

	return true, nil
}

// Make custom relative path.
// In every path remove first directory name unless if the path has only one directory.
func makeRelativePath(path string) string {
	// Split the path into directories
	dirs := strings.Split(path, "/")
	// Remove first directory name unless it is "root"
	if len(dirs) > 1 {
		dirs = dirs[1:]
	}
	// Join the directories back into a path
	return strings.Join(dirs, "/")
}

func CompareReports(oldReportFile string, newReportFile string) (map[string]string, error) {
	oldReport, err := LoadReport(oldReportFile)
	if err != nil {
		return nil, fmt.Errorf("error loading old report: %v", err)
	}

	newReport, err := LoadReport(newReportFile)
	if err != nil {
		return nil, fmt.Errorf("error loading new report: %v", err)
	}

	// Check if the report_checksum  matches the checksum of both reports matches the checksum of the files: [] in the report
	// If not, the report is corrupted
	if ok, err := checkReportChecksum("old", oldReport); !ok {
		return nil, fmt.Errorf("error checking old report checksum: %v", err)
	}

	if ok, err := checkReportChecksum("new", newReport); !ok {
		return nil, fmt.Errorf("error checking new report checksum: %v", err)
	}

	// Compare the reports
	changes := make(map[string]string)

	// Sort the files in the report by name
	sort.Slice(oldReport.Files, func(i, j int) bool {
		return oldReport.Files[i].Name < oldReport.Files[j].Name
	})
	sort.Slice(newReport.Files, func(i, j int) bool {
		return newReport.Files[i].Name < newReport.Files[j].Name
	})

	// make all the paths relative
	for i := range oldReport.Files {
		oldReport.Files[i].Path = makeRelativePath(oldReport.Files[i].Path)
	}
	for i := range newReport.Files {
		newReport.Files[i].Path = makeRelativePath(newReport.Files[i].Path)
	}

	// Maps List
	addedFiles := make(map[string]string)
	deletedFiles := make(map[string]string)
	modifiedFiles := make(map[string]string)
	toRenamedFiles := make(map[string]string)
	fromRenamedFiles := make(map[string]string)

	// Segregate the deleted files, if the checksums are not present in the new report
	// Save them in the format: "file_name"
	log.Default().Println("INFO: Segregating the deleted files.")
	for _, oldFile := range oldReport.Files {
		var found bool
		for _, newFile := range newReport.Files {
			if newFile.SHA256 == oldFile.SHA256 || newFile.Path == oldFile.Path {
				found = true
				break
			}
		}
		if !found {
			// changes[oldFile.Path] = "deleted"
			deletedFiles[oldFile.Path] = oldFile.SHA256
		}
	}

	// Deleted
	for _, oldFile := range oldReport.Files {
		var found bool
		for _, newFile := range newReport.Files {
			if oldFile.Path == newFile.Path {
				found = true
				break
			}
		}
		if !found {
			// changes[oldFile.Path] = "deleted"
			deletedFiles[oldFile.Path] = oldFile.SHA256
		}
	}

	// Segregate the renamed files, if the checksums are same but the file names are different.
	// Save them in the format: "old_file_name -> new_file_name"
	// If file found in global_duplicates, then skip it
	// If the same file name is found in the new report, then skip it
	log.Default().Println("INFO: Segregating the renamed files.")
	for _, oldFile := range oldReport.Files {
		// Check if the file is present in the new report and has a different name
		for _, newFile := range newReport.Files {
			if oldFile.SHA256 == newFile.SHA256 && oldFile.Name != newFile.Name {
				// Check if the SHA256 of the new file is unique in the new report
				count := 0
				for _, otherFile := range newReport.Files {
					if otherFile.SHA256 == newFile.SHA256 {
						count++
					}
				}
				if count == 1 {
					// Add the renamed file to the changes map with label "renamed"
					toRenamedFiles[oldFile.Path] = newFile.Path
					fromRenamedFiles[newFile.Path] = oldFile.Path
				}
			}
		}
	}

	// Segregate the added files, if the checksums are not present in the old report
	// Save them in the format: "file_name"
	log.Default().Println("INFO: Segregating the added files.")
	for _, newFile := range newReport.Files {
		var found bool
		for _, oldFile := range oldReport.Files {
			if oldFile.SHA256 == newFile.SHA256 {
				if newFile.Path == oldFile.Path {
					found = true
					break
				}
			}
		}
		if !found {
			addedFiles[newFile.Path] = newFile.SHA256
		}
	}

	// Segregate the modified files, if the checksums are different
	// Save them in the format: "file_name"
	log.Default().Println("INFO: Segregating the altered files.")
	// Check for modified files
	for _, oldFile := range oldReport.Files {
		var found bool
		for _, newFile := range newReport.Files {
			if oldFile.Path != newFile.Path && oldFile.SHA256 == newFile.SHA256 {
				found = true
				break
			}
			if _, ok := toRenamedFiles[oldFile.Path]; ok {
				found = true
				break
			}
			if _, ok := deletedFiles[oldFile.Path]; ok {
				found = true
				break
			}
		}
		if !found {
			changes[oldFile.Path] = "modified"
			modifiedFiles[oldFile.Path] = oldFile.SHA256
		}
	}

	// Segregate the duplicated files, if multiple files have the same checksum with different file names
	// Create a map to store the sha256 values and their corresponding file names for both reports
	oldSha256Map := make(map[string][]string)
	newSha256Map := make(map[string][]string)

	// Loop through the files in the old report and add their sha256 values to the map
	for _, file := range oldReport.Files {
		sha256 := file.SHA256
		oldSha256Map[sha256] = append(oldSha256Map[sha256], file.Name)
	}

	// Loop through the files in the new report and add their sha256 values to the map
	for _, file := range newReport.Files {
		sha256 := file.SHA256
		newSha256Map[sha256] = append(newSha256Map[sha256], file.Name)
	}

	// Loop through the sha256 values in the old report sha256 map and check if they are present in the new report sha256 map
	for sha256, oldFileNames := range oldSha256Map {
		newFileNames := newSha256Map[sha256]
		if len(newFileNames) > 0 {
			// Duplicate file(s) found
			if len(oldFileNames) > 1 || len(newFileNames) > 1 {
				// Only consider the duplicate files that are present in both reports
				var commonFileNames []string
				for _, oldFileName := range oldFileNames {
					for _, newFileName := range newFileNames {
						if oldFileName == newFileName {
							commonFileNames = append(commonFileNames, oldFileName)
							break
						}
					}
				}
				if len(commonFileNames) > 0 {
					changes[sha256] = "duplicated"
				}
			}
		}
	}

	// Segregate the unchanged files, if the checksums are same and the file names are same
	// Save them in the format: "file_name"
	for _, oldFile := range oldReport.Files {
		for _, newFile := range newReport.Files {
			if oldFile.SHA256 == newFile.SHA256 && oldFile.Name == newFile.Name {
				changes[oldFile.Path] = "unchanged"
			}
		}
	}

	// Add file to the changes map with label "added" if the file is present in addedFiles map, and not present in renamedFiles map, deletedFiles map, and modifiedFiles map
	for _, newFile := range newReport.Files {
		if _, ok := addedFiles[newFile.Path]; ok {
			if _, ok := fromRenamedFiles[newFile.Path]; !ok {
				if _, ok := deletedFiles[newFile.Path]; !ok {
					if _, ok := modifiedFiles[newFile.Path]; !ok {
						changes[newFile.Path] = "added"
					}
				}
			}
		}
	}

	// Add the deleted files to the changes map with label "deleted" if the file is not present in addedFiles map, renamedFiles map, and modifiedFiles map
	for _, oldFile := range oldReport.Files {
		if _, ok := deletedFiles[oldFile.Path]; ok {
			if _, ok := addedFiles[oldFile.Path]; !ok {
				if _, ok := toRenamedFiles[oldFile.Path]; !ok {
					if _, ok := modifiedFiles[oldFile.Path]; !ok {
						changes[oldFile.Path] = "deleted"
					}
				}
			}
		}
	}

	// Add the renamed files to the changes map with label "renamed" if the file is not present in addedFiles map, deletedFiles map, and modifiedFiles map
	for _, oldFile := range oldReport.Files {
		if _, ok := toRenamedFiles[oldFile.Path]; ok {
			for oldPath, newPath := range toRenamedFiles {
				changes[oldPath+" -> "+newPath] = "renamed"
			}
		}
	}

	// Get working directories from the reports "working_dir" field.
	// Save them in the format: "Old report: old_working_dir\n\t\t\t\tNew report: new_working_dir"
	oldReportDir := oldReport.WorkingDir
	newReportDir := newReport.WorkingDir
	workingDirMessage := "Old report: " + oldReportDir + "\n\t\t\t\tNew report: " + newReportDir
	changes[workingDirMessage] = "workingDirs"

	return changes, nil
}
