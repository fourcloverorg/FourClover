package commands

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	common "fourclover.org/internal/common"
	config "fourclover.org/internal/config"
	report "fourclover.org/internal/report"
	snapshot "fourclover.org/internal/snapshot"
)

// Snapshot command flags
var (
	dir                  string
	reportFileExportPath string
	snapshotName         string
	supressScanOutput    bool
	saveReportToFile     bool
	saveReportToDatabase bool
)
var (
	snapCommand         *flag.FlagSet
	directoryToSnapshot common.DirectoryToSnapshot
	hashAlgorithms      common.HashAlgorithms
	excludeThem         common.ExcludeThem
	focusExtensions     common.FocusExtensions
)

var (
	err                   error
	compareCommand        *flag.FlagSet
	report1File           string
	report2File           string
	compareReportFileName string
	supressCompareOutput  bool
	compareStartTime      time.Time
	CompareResult         map[string]string
	totalFilesObserved    int
	compareEndTime        time.Time
	compareDuration       time.Duration
)

func Snapshot(configData []byte) (bool, error) {
	timeNow := time.Now().Format("2006-01-02_15-04-05")
	if configData == nil {
		log.Default().Println("INFO: Reading snapshot CLI arguments.")

		snapCommand = flag.NewFlagSet("snapshot", flag.ExitOnError)
		snapCommand.Var(&directoryToSnapshot, "dir", "Directory for the snapshot.")
		// eg: fourclover snapshot -dir [target directory] -hashs sha256,sha512
		snapCommand.Var(&hashAlgorithms, "hashs", "[Optional] Additional hash algorithms to include in the report (comma-separated) (default sha256)")
		// eg: fourclover snapshot -dir . -exclude .git,.idea,node_modules,venv
		snapCommand.Var(&excludeThem, "exclude", "[Optional] Exclude files and directories from the snapshot (comma-separated)")
		// eg: fourclover snapshot -dir . -focus .go,.py,.js,.html,.css,.json,.txt,.md,.yml,.gitkeep,.gitlab-ci.yml
		snapCommand.Var(&focusExtensions, "focus", "[Optional] Focus on certain file extensions only (comma-separated eg. .py,.js) (All files will be added snapshot if not specified)")
		snapCommand.BoolVar(&supressScanOutput, "supress", false, "[Optional] Supress the output of the program. (default: false)")
		snapCommand.StringVar(&snapshotName, "name", "snapshot_"+timeNow, "[Optional] Name of the snapshot (default: snapshot_"+timeNow+")")
		// Save snapshot report to file decision
		snapCommand.BoolVar(&saveReportToFile, "save", true, "[Optional] Save snapshot report to file. (default: true)")
		// Export snapshot report to file path
		snapCommand.StringVar(&reportFileExportPath, "out", "snapshot_report_"+timeNow+".json", "[Optional] Export report file to path. (default: snapshot_report_"+timeNow+".json)")
		// Save snapshot report to database decision
		snapCommand.BoolVar(&saveReportToDatabase, "db", false, "[Optional] Save snapshot report to database. (Note: Default SQLite path will be used) (default: false)")

		snapCommand.Parse(os.Args[2:])

		if len(os.Args) == 2 {
			log.Default().Println("Usage of snapshot:")
			snapCommand.PrintDefaults()
			return false, nil
		}
	}

	// If the config file has data in it, then use that data to set the values for the variables. If not, then use the default values.
	if configData != nil {
		// Get directory for the snapshot
		fourcloverSnapshotDirectory := config.GetFourCloverSnapshotDirectory()
		if fourcloverSnapshotDirectory != "" {
			log.Default().Println("INFO: Using directory for the snapshot specified in the config file or the environment variables.")
			directoryToSnapshot = common.DirectoryToSnapshot(fourcloverSnapshotDirectory)
		} else {
			log.Fatalf("ERROR: No directory for the snapshot specified in the config file or the environment variables.")
		}

		// Get Hash Algorithms
		fourcloverSnapshotHashAlgorithms := config.GetFourCloverSnapshotHashAlgorithm()
		if fourcloverSnapshotHashAlgorithms != nil {
			log.Default().Println("INFO: Using hash algorithms specified in the config file or the environment variables.")
			hashAlgorithms = common.HashAlgorithms(fourcloverSnapshotHashAlgorithms)
		} else {
			log.Default().Println("INFO: No hash algorithms specified in the config file or the environment variables. Using default hash algorithm: sha256")
		}

		// Get Exclude Them
		fourcloverSnapshotExcludeThem := config.GetFourCloverSnapshotExcludeThem()
		if fourcloverSnapshotExcludeThem != nil {
			log.Default().Println("INFO: Using files or directories to exclude specified in the config file or the environment variables.")
			excludeThem = common.ExcludeThem(fourcloverSnapshotExcludeThem)
		} else {
			log.Default().Println("INFO: No files or directories to exclude specified in the config file or the environment variables.")
		}

		// Get Focus Extensions
		fourcloverSnapshotFocusExtensions := config.GetFourCloverSnapshotFocusExtensions()
		if fourcloverSnapshotFocusExtensions != nil {
			log.Default().Println("INFO: Using file extensions to focus on specified in the config file or the environment variables.")
			focusExtensions = common.FocusExtensions(fourcloverSnapshotFocusExtensions)
		} else {
			log.Default().Println("INFO: No file extensions to focus on specified in the config file or the environment variables.")
		}

		// Get Save SnapshotReport To File
		fourcloverSaveReportToFile := config.GetFourCloverSaveReportToFile()
		if !fourcloverSaveReportToFile {
			log.Default().Println("INFO: Using save report to file specified in the config file or the environment variables.")
			saveReportToFile = fourcloverSaveReportToFile
		} else {
			log.Default().Println("INFO: No save report to file specified in the config file or the environment variables.")
		}

		// Get Save SnapshotReport To Database
		fourcloverSaveReportToDatabase := config.GetFourCloverSaveReportToDatabase()
		if !fourcloverSaveReportToDatabase {
			log.Default().Println("INFO: Using save report to database specified in the config file or the environment variables.")
			saveReportToDatabase = fourcloverSaveReportToDatabase
		} else {
			log.Default().Println("INFO: No save report to database specified in the config file or the environment variables.")
		}

		// Get Snapshot Export SnapshotReport Path
		fourcloverSnapshotExportReportPath := config.GetFourCloverSnapshotExportReportPath()
		if fourcloverSnapshotExportReportPath != "" {
			log.Default().Println("INFO: Using export report file specified in the config file or the environment variables.")
			reportFileExportPath = fourcloverSnapshotExportReportPath
		} else {
			log.Default().Println("INFO: No export report file specified in the config file or the environment variables. Using default export report file: snapshot_report_" + timeNow + ".json")
			reportFileExportPath = "snapshot_report_" + timeNow + ".json"
		}

		// Get Snapshot Name
		fourcloverSnapshotName := config.GetFourCloverSnapshotName()
		if fourcloverSnapshotName != "" {
			log.Default().Println("INFO: Using snapshot name specified in the config file or the environment variables.")
			snapshotName = fourcloverSnapshotName
		} else {
			log.Default().Println("INFO: No snapshot name specified in the config file or the environment variables. Using default snapshot name: snapshot_" + timeNow)
			snapshotName = "snapshot_" + timeNow
		}
	}

	fmt.Sprintln("INFO: Snapshot has been started. Please wait...")
	snapshotStartTime := time.Now()

	// Add the default hash algorithm to the list of hash algorithms "sha256"
	hashAlgorithms = append(hashAlgorithms, "sha256")

	// Add the default file extensions to include in the snapshot
	if len(focusExtensions) == 0 {
		focusExtensions = append(focusExtensions, ".*")
	}

	dir = directoryToSnapshot.String()
	dirSnapshotReport, err := snapshot.SnapshotDirectory(dir, hashAlgorithms, excludeThem, focusExtensions, snapshotName)
	if err != nil {
		err = fmt.Errorf("failed to snapshot directory %s: %s", dir, err)
		return false, err
	}

	// Save the report to a file
	if saveReportToFile {
		if reportFileExportPath == "" {
			reportFileExportPath = "snapshot_report_" + time.Now().Format("2006-01-02_15-04-05") + ".json"
		}
		err = report.SaveSnapshotReport(dirSnapshotReport, reportFileExportPath)
		if err != nil {
			err = fmt.Errorf("failed to save report to %s: %s", reportFileExportPath, err)
			return false, err
		}
	}

	snapshotEndTime := time.Now()
	snapshotDuration := snapshotEndTime.Sub(snapshotStartTime)
	snapshotDuration = snapshotDuration.Round(time.Second) // Round the duration to the nearest second
	totalFiles := len(dirSnapshotReport.Files)
	totalSize := dirSnapshotReport.TotalFilesSize
	log.Default().Println("INFO: Snapshot completed in", snapshotDuration, "for", totalFiles, "files", "(", totalSize, "bytes", ")")
	log.Default().Println("INFO: SnapshotReport saved to", reportFileExportPath)
	return true, nil
}

// Compare compares two reports and returns a comparison report
func Compare(configData []byte) (bool, error) {
	timeNow := time.Now().Format("2006-01-02_15-04-05")
	if configData == nil {
		log.Default().Println("INFO: Reading compare CLI arguments.")
		compareCommand = flag.NewFlagSet("compare", flag.ExitOnError)
		compareCommand.StringVar(&report1File, "old", "", "Old report file")
		compareCommand.StringVar(&report2File, "new", "", "New report file")
		compareCommand.StringVar(&compareReportFileName, "out", "comparison_report_"+time.Now().Format("2006-01-02_15-04-05")+".txt", "Export comparison report file")
		compareCommand.BoolVar(&supressCompareOutput, "supress", false, "[Optional] Supress the output of the program")

		compareCommand.Parse(os.Args[2:])

		if len(os.Args) == 2 {
			log.Default().Println("Usage of compare:")
			compareCommand.PrintDefaults()
			return false, nil
		}
	}

	if configData != nil {

		fourcloverCompareOldReportPath := config.GetFourCloverCompareOldReportPath()
		if fourcloverCompareOldReportPath != "" {
			log.Default().Println("INFO: Reading old report file from config file or environment variable.")
			report1File = fourcloverCompareOldReportPath
		} else {
			log.Fatal("ERROR: error parsing config file or environment variable: old_report or FOURCLOVER_COMPARE_OLD_REPORT_PATH is missing.")
		}

		fourcloverCompareNewReportPath := config.GetFourCloverCompareNewReportPath()
		if fourcloverCompareNewReportPath != "" {
			log.Default().Println("INFO: Reading new report file from config file or environment variable.")
			report2File = fourcloverCompareNewReportPath
		} else {
			log.Fatal("ERROR: error parsing config file or environment variable: new_report or FOURCLOVER_COMPARE_NEW_REPORT_PATH is missing.")
		}

		fourcloverCompareExportReportPath := config.GetFourCloverCompareExportReportPath()
		if fourcloverCompareExportReportPath != "" {
			log.Default().Println("INFO: Reading export report file from config file or environment variable.")
			compareReportFileName = fourcloverCompareExportReportPath
		} else {
			log.Default().Println("ERROR: error parsing config file or environment variable: export_report or FOURCLOVER_COMPARE_EXPORT_REPORT_PATH is missing.")
			log.Default().Println("INFO: Using default value for export_report.", "Default value:", "comparison_report_"+timeNow+".txt")
			compareReportFileName = "comparison_report_" + timeNow + ".txt"
		}
	}

	compareStartTime = time.Now()
	log.Default().Println("INFO: Starting report comparison... Please wait...")
	CompareResult, err = snapshot.CompareReports(report1File, report2File)
	if err != nil {
		return false, err
	}

	// Print and Save the comparison report to a file
	if compareReportFileName == "" {
		compareReportFileName = "comparison_report_" + timeNow + ".txt"
	}
	totalFilesObserved, err = report.SaveCompareReport(compareReportFileName, CompareResult)
	if err != nil {
		return false, err
	}

	compareEndTime = time.Now()
	compareDuration = compareEndTime.Sub(compareStartTime)
	compareDuration = compareDuration.Round(time.Second)

	log.Default().Println("INFO: Compare completed in", compareDuration, "for", totalFilesObserved, "files.")
	return true, nil
}
