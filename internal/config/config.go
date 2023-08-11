package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Check if cli argument has certain key. Return true if it has, false if it doesn't.
func CheckCliArg(key string) bool {
	for _, arg := range os.Args {
		if arg == key {
			return true
		}
	}
	return false
}

// Check if fourclover config values is provided in config file and then in system environment variable.
// Use the fourclover version provided in config file or system environment variable if it exists.

// ------------------------ Global ------------------------ //
// Get fourclover version
func GetFourCloverVersion() string {
	var fourcloverVersion string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverVersionEnv, _ := GetSystemEnvVar("FOURCLOVER_VERSION")

	if yamlConfig.Version != "" && fourcloverVersionEnv == "" {
		fourcloverVersion = yamlConfig.Version
	} else if fourcloverVersionEnv != "" {
		fourcloverVersion = fourcloverVersionEnv
	} else {
		log.Default().Println("ERROR: No fourclover version specified in the config file or system environment variable.")
		return ""
	}

	return fourcloverVersion
}

// Get fourclover help
func GetFourCloverHelp() bool {
	fourcloverIsHelpArg := CheckCliArg("help")
	return fourcloverIsHelpArg
}

// Get fourclover demo
func GetFourCloverDemo() bool {
	fourcloverIsDemoArg := CheckCliArg("demo")
	return fourcloverIsDemoArg
}

// Get fourclover is_snapshot status
func GetFourCloverIsSnapshot() bool {
	var fourcloverIsSnapshot bool
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverIsSnapshotEnv, _ := GetSystemEnvVar("FOURCLOVER_IS_SNAPSHOT")
	fourcloverIsSnapshotArg := CheckCliArg("snapshot")
	if fourcloverIsSnapshotArg {
		fourcloverIsSnapshot = true
	} else if yamlConfig.Global.IsSnapshot {
		fourcloverIsSnapshot = true
	} else if !yamlConfig.Global.IsSnapshot {
		fourcloverIsSnapshot = false
	} else if fourcloverIsSnapshotEnv == "true" {
		fourcloverIsSnapshot = true
	} else if fourcloverIsSnapshotEnv == "false" {
		fourcloverIsSnapshot = false
	} else {
		log.Fatalf("ERROR: No fourclover is_snapshot status specified in the config file or system environment variable.")
		return false
	}

	return fourcloverIsSnapshot
}

// Get fourclover is_compare status
func GetFourCloverIsCompare() bool {
	var fourcloverIsCompare bool
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverIsCompareEnv, _ := GetSystemEnvVar("FOURCLOVER_IS_COMPARE")
	fourcloverIsCompareArg := CheckCliArg("compare")
	if fourcloverIsCompareArg {
		fourcloverIsCompare = true
	} else if yamlConfig.Global.IsCompare {
		fourcloverIsCompare = true
	} else if !yamlConfig.Global.IsCompare {
		fourcloverIsCompare = false
	} else if fourcloverIsCompareEnv == "true" {
		fourcloverIsCompare = true
	} else if fourcloverIsCompareEnv == "false" {
		fourcloverIsCompare = false
	} else {
		log.Fatalf("ERROR: No fourclover is_compare status specified in the config file or system environment variable.")
		return false
	}

	return fourcloverIsCompare
}

// Get fourclover is_policy status
func GetFourCloverIsPolicy() bool {
	var fourcloverIsPolicy bool
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverIsPolicyEnv, _ := GetSystemEnvVar("FOURCLOVER_IS_POLICY")
	fourcloverIsPolicyArg := CheckCliArg("policy")
	if fourcloverIsPolicyArg {
		fourcloverIsPolicy = true
	} else if yamlConfig.Global.IsPolicy {
		fourcloverIsPolicy = true
	} else if !yamlConfig.Global.IsPolicy {
		fourcloverIsPolicy = false
	} else if fourcloverIsPolicyEnv == "true" {
		fourcloverIsPolicy = true
	} else if fourcloverIsPolicyEnv == "false" {
		fourcloverIsPolicy = false
	} else {
		log.Fatalf("ERROR: No fourclover is_policy status specified in the config file or system environment variable.")
		return false
	}

	return fourcloverIsPolicy
}

// Get fourclover save_report_to_file status
func GetFourCloverSaveReportToFile() bool {
	var fourcloverSaveToFile bool
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSaveToFileEnv, _ := GetSystemEnvVar("FOURCLOVER_SAVE_TO_FILE")
	fourcloverSaveToFileArg := CheckCliArg("-sr")
	if fourcloverSaveToFileArg {
		fourcloverSaveToFile = true
	} else if yamlConfig.Global.SaveReportToFile {
		fourcloverSaveToFile = true
	} else if !yamlConfig.Global.SaveReportToFile {
		fourcloverSaveToFile = false
	} else if fourcloverSaveToFileEnv == "true" {
		fourcloverSaveToFile = true
	} else if fourcloverSaveToFileEnv == "false" {
		fourcloverSaveToFile = false
	} else {
		log.Default().Println("WARN: No fourclover save_to_file status specified in the config file or system environment variable. Using default value of true.")
		fourcloverSaveToFile = true
	}

	return fourcloverSaveToFile
}

// Get fourclover save_report_to_database status
func GetFourCloverSaveReportToDatabase() bool {
	var fourcloverSaveToDatabase bool
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSaveToDatabaseEnv, _ := GetSystemEnvVar("FOURCLOVER_SAVE_TO_DATABASE")
	fourcloverSaveToDatabaseArg := CheckCliArg("-db")
	if fourcloverSaveToDatabaseArg {
		fourcloverSaveToDatabase = true
	} else if yamlConfig.Global.SaveReportToDatabase {
		fourcloverSaveToDatabase = true
	} else if !yamlConfig.Global.SaveReportToDatabase {
		fourcloverSaveToDatabase = false
	} else if fourcloverSaveToDatabaseEnv == "true" {
		fourcloverSaveToDatabase = true
	} else if fourcloverSaveToDatabaseEnv == "false" {
		fourcloverSaveToDatabase = false
	} else {
		log.Default().Println("WARN: No fourclover save_to_database status specified in the config file or system environment variable. Using default value of false.")
		fourcloverSaveToDatabase = false
	}

	return fourcloverSaveToDatabase
}

// Get fourclover start_server status
func GetFourCloverStartServer() bool {
	var fourcloverStartServer bool
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverStartServerEnv, _ := GetSystemEnvVar("FOURCLOVER_START_SERVER")
	if yamlConfig.Global.StartServer {
		fourcloverStartServer = true
	} else if !yamlConfig.Global.StartServer {
		fourcloverStartServer = false
	} else if fourcloverStartServerEnv == "true" {
		fourcloverStartServer = true
	} else if fourcloverStartServerEnv == "false" {
		fourcloverStartServer = false
	} else {
		log.Default().Println("WARN: No fourclover start_server status specified in the config file or system environment variable. Using default value of false.")
		fourcloverStartServer = false
	}

	return fourcloverStartServer
}

// Get fourclover suppress_logs status
func GetFourCloverSuppressLogs() bool {
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSuppressLogsEnv, _ := GetSystemEnvVar("FOURCLOVER_SUPPRESS_LOGS")
	fourcloverSuppressLogsArg := CheckCliArg("-supress")
	if fourcloverSuppressLogsArg {
		return true
	}
	if yamlConfig.Global.SuppressLogs {
		return true
	}
	if !yamlConfig.Global.SuppressLogs {
		return false
	}
	if fourcloverSuppressLogsEnv == "true" {
		return true
	}
	if fourcloverSuppressLogsEnv == "false" {
		return false
	} else {
		log.Default().Println("WARN: No fourclover suppress_logs status specified in the config file or system environment variable. Using default value of false.")
		return false
	}
}

// ------------------------ Snapshot Configurations ------------------------ //

// Get fourclover snapshot name
func GetFourCloverSnapshotName() string {
	var fourcloverSnapshotName string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSnapshotNameEnv, _ := GetSystemEnvVar("FOURCLOVER_SNAPSHOT_NAME")
	if yamlConfig.Snapshot.Name != "" && fourcloverSnapshotNameEnv == "" {
		fourcloverSnapshotName = yamlConfig.Snapshot.Name
	} else if fourcloverSnapshotNameEnv != "" {
		fourcloverSnapshotName = fourcloverSnapshotNameEnv
	} else {
		log.Default().Println("WARN: No fourclover snapshot name specified in the config file or system environment variable.")
		return ""
	}

	return fourcloverSnapshotName
}

// Get fourclover snapshot directory
func GetFourCloverSnapshotDirectory() string {
	var fourcloverSnapshotDirectory string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSnapshotDirectoryEnv, _ := GetSystemEnvVar("FOURCLOVER_SNAPSHOT_DIRECTORY")
	if yamlConfig.Snapshot.Directory != "" && fourcloverSnapshotDirectoryEnv == "" {
		fourcloverSnapshotDirectory = yamlConfig.Snapshot.Directory
	} else if fourcloverSnapshotDirectoryEnv != "" {
		fourcloverSnapshotDirectory = fourcloverSnapshotDirectoryEnv
	} else {
		log.Default().Println("ERROR: No fourclover snapshot directory specified in the config file or system environment variable.")
		return ""
	}

	return fourcloverSnapshotDirectory
}

// Get fourclover snapshot hash_algorithms
func GetFourCloverSnapshotHashAlgorithm() []string {
	var fourcloverSnapshotHashAlgorithm []string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSnapshotHashAlgorithmEnv, _ := GetSystemEnvVar("FOURCLOVER_SNAPSHOT_HASH_ALGORITHM")
	if yamlConfig.Snapshot.HashAlgorithms != nil && fourcloverSnapshotHashAlgorithmEnv == "" {
		fourcloverSnapshotHashAlgorithm = yamlConfig.Snapshot.HashAlgorithms
	} else if fourcloverSnapshotHashAlgorithmEnv != "" {
		fourcloverSnapshotHashAlgorithm = []string{fourcloverSnapshotHashAlgorithmEnv}
	} else {
		log.Default().Println("ERROR: No fourclover snapshot hash_algorithm specified in the config file or system environment variable.")
		return []string{}
	}

	return fourcloverSnapshotHashAlgorithm
}

// Get fourclover snapshot exclude_them
func GetFourCloverSnapshotExcludeThem() []string {
	var fourcloverSnapshotExcludeThem []string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSnapshotExcludeThemEnv, _ := GetSystemEnvVar("FOURCLOVER_SNAPSHOT_EXCLUDE_THEM")
	if yamlConfig.Snapshot.ExcludeThem != nil && fourcloverSnapshotExcludeThemEnv == "" {
		fourcloverSnapshotExcludeThem = yamlConfig.Snapshot.ExcludeThem
	} else if fourcloverSnapshotExcludeThemEnv != "" {
		fourcloverSnapshotExcludeThem = []string{fourcloverSnapshotExcludeThemEnv}
	} else {
		log.Default().Println("ERROR: No fourclover snapshot exclude_them specified in the config file or system environment variable.")
		return []string{}
	}

	return fourcloverSnapshotExcludeThem
}

// Get fourclover snapshot focus_extensions
func GetFourCloverSnapshotFocusExtensions() []string {
	var fourcloverSnapshotFocusExtensions []string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSnapshotFocusExtentionsEnv, _ := GetSystemEnvVar("FOURCLOVER_SNAPSHOT_FOCUS_EXTENSIONS")
	if yamlConfig.Snapshot.FocusExtensions != nil && fourcloverSnapshotFocusExtentionsEnv == "" {
		fourcloverSnapshotFocusExtensions = yamlConfig.Snapshot.FocusExtensions
	} else if fourcloverSnapshotFocusExtentionsEnv != "" {
		fourcloverSnapshotFocusExtensions = []string{fourcloverSnapshotFocusExtentionsEnv}
	} else {
		log.Default().Println("ERROR: No fourclover snapshot focus_extensions specified in the config file or system environment variable.")
		return []string{}
	}

	return fourcloverSnapshotFocusExtensions
}

// Get fourclover snapshot export_report_path
func GetFourCloverSnapshotExportReportPath() string {
	if !GetFourCloverSaveReportToFile() {
		log.Default().Println("WARN: fourclover snapshot export_report is not needed because save_to_file is set to false.")
		return ""
	}

	var fourcloverSnapshotExportReportPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverSnapshotExportReportPathEnv, _ := GetSystemEnvVar("FOURCLOVER_SNAPSHOT_EXPORT_REPORT")
	if yamlConfig.Snapshot.ExportReportPath != "" && fourcloverSnapshotExportReportPathEnv == "" {
		fourcloverSnapshotExportReportPath = yamlConfig.Snapshot.ExportReportPath
	} else if fourcloverSnapshotExportReportPathEnv != "" {
		fourcloverSnapshotExportReportPath = fourcloverSnapshotExportReportPathEnv
	} else {
		log.Default().Println("ERROR: No fourclover snapshot export_report specified in the config file or system environment variable. SnapshotReport will not be exported.")
		return ""
	}

	return fourcloverSnapshotExportReportPath
}

// ----------------------- Compare Snapshot Configurations ----------------------- //

// Get fourclover compare old_report_path
func GetFourCloverCompareOldReportPath() string {
	var fourcloverCompareOldReportPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverCompareOldReportPathEnv, _ := GetSystemEnvVar("FOURCLOVER_COMPARE_OLD_REPORT_PATH")
	if yamlConfig.Compare.OldReportPath != "" && fourcloverCompareOldReportPathEnv == "" {
		fourcloverCompareOldReportPath = yamlConfig.Compare.OldReportPath
	} else if fourcloverCompareOldReportPathEnv != "" {
		fourcloverCompareOldReportPath = fourcloverCompareOldReportPathEnv
	} else {
		log.Fatalf("ERROR: No fourclover compare old_report path specified in the config file or system environment variable.")
		return ""
	}

	return fourcloverCompareOldReportPath
}

// Get fourclover compare new_report_path
func GetFourCloverCompareNewReportPath() string {
	var fourcloverCompareNewReportPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverCompareNewReportPathEnv, _ := GetSystemEnvVar("FOURCLOVER_COMPARE_NEW_REPORT_PATH")
	if yamlConfig.Compare.NewReportPath != "" && fourcloverCompareNewReportPathEnv == "" {
		fourcloverCompareNewReportPath = yamlConfig.Compare.NewReportPath
	} else if fourcloverCompareNewReportPathEnv != "" {
		fourcloverCompareNewReportPath = fourcloverCompareNewReportPathEnv
	} else {
		log.Fatalf("ERROR: No fourclover compare new_report path specified in the config file or system environment variable.")
		return ""
	}

	return fourcloverCompareNewReportPath
}

// Get fourclover compare export_report_path
func GetFourCloverCompareExportReportPath() string {
	var fourcloverCompareExportReportPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverCompareExportReportPathEnv, _ := GetSystemEnvVar("FOURCLOVER_COMPARE_EXPORT_REPORT_PATH")
	if yamlConfig.Compare.ExportReportPath != "" && fourcloverCompareExportReportPathEnv == "" {
		fourcloverCompareExportReportPath = yamlConfig.Compare.ExportReportPath
	} else if fourcloverCompareExportReportPathEnv != "" {
		fourcloverCompareExportReportPath = fourcloverCompareExportReportPathEnv
	} else {
		log.Fatalf("ERROR: No fourclover compare export_report path specified in the config file or system environment variable.")
		return ""
	}

	return fourcloverCompareExportReportPath
}

// ----------------------- Policy Configurations ----------------------- //

// Get fourclover snapshot_report_path
func GetFourCloverPolicySnapshotReportPath() string {
	var GetFourCloverPolicySnapshotReportPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	GetFourCloverPolicySnapshotReportPathEnv, _ := GetSystemEnvVar("FOURCLOVER_POLICY_SNAPSHOT_REPORT_PATH")
	if yamlConfig.Policy.SnapshotReportPath != "" && GetFourCloverPolicySnapshotReportPathEnv == "" {
		GetFourCloverPolicySnapshotReportPath = yamlConfig.Policy.SnapshotReportPath
	} else if GetFourCloverPolicySnapshotReportPathEnv != "" {
		GetFourCloverPolicySnapshotReportPath = GetFourCloverPolicySnapshotReportPathEnv
	} else {
		log.Fatalf("ERROR: No fourclover policy snapshot_report_path specified in the config file or system environment variable.")
		return ""
	}

	return GetFourCloverPolicySnapshotReportPath
}

// Get fourclover target_directory_path
func GetFourCloverPolicyTargetDirPath() string {
	var GetFourCloverPolicyTargetDirPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	GetFourCloverPolicyTargetDirPathEnv, _ := GetSystemEnvVar("FOURCLOVER_POLICY_TARGET_DIRECTORY_PATH")
	if yamlConfig.Policy.TargetDirPath != "" && GetFourCloverPolicyTargetDirPathEnv == "" {
		GetFourCloverPolicyTargetDirPath = yamlConfig.Policy.TargetDirPath
	} else if GetFourCloverPolicyTargetDirPathEnv != "" {
		GetFourCloverPolicyTargetDirPath = GetFourCloverPolicyTargetDirPathEnv
	} else {
		log.Fatalf("ERROR: No fourclover policy target_directory_path specified in the config file or system environment variable.")
		return ""
	}

	return GetFourCloverPolicyTargetDirPath
}

// Get fourclover policy directory_path
func GetFourCloverPolicyDirPath() string {
	var GetFourCloverPolicyDirPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	GetFourCloverPolicyDirPathEnv, _ := GetSystemEnvVar("FOURCLOVER_POLICY_DIRECTORY_PATH")
	if yamlConfig.Policy.DirPath != "" && GetFourCloverPolicyDirPathEnv == "" {
		GetFourCloverPolicyDirPath = yamlConfig.Policy.DirPath
	} else if GetFourCloverPolicyDirPathEnv != "" {
		GetFourCloverPolicyDirPath = GetFourCloverPolicyDirPathEnv
	} else {
		log.Fatalf("ERROR: No fourclover policy directory_path specified in the config file or system environment variable.")
		return ""
	}

	return GetFourCloverPolicyDirPath
}

// Get fourclover policy export_report_path
func GetFourCloverPolicyExportReportPath() string {
	var GetFourCloverPolicyExportReportPath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	GetFourCloverPolicyExportReportPathEnv, _ := GetSystemEnvVar("FOURCLOVER_POLICY_EXPORT_REPORT_PATH")
	if yamlConfig.Policy.ExportReportPath != "" && GetFourCloverPolicyExportReportPathEnv == "" {
		GetFourCloverPolicyExportReportPath = yamlConfig.Policy.ExportReportPath
	} else if GetFourCloverPolicyExportReportPathEnv != "" {
		GetFourCloverPolicyExportReportPath = GetFourCloverPolicyExportReportPathEnv
	} else {
		log.Fatalf("ERROR: No fourclover policy export_report_path specified in the config file or system environment variable.")
		return ""
	}

	return GetFourCloverPolicyExportReportPath
}

// ----------------------- Database Configurations ----------------------- //

// Get fourclover database sqlite path
func GetFourCloverDatabaseSqlitePath() string {
	if !GetFourCloverSaveReportToDatabase() {
		log.Default().Println("WARN: fourclover database sqlite path is not needed because save_to_database is set to false.")
	}

	defaultDBPath := "fourclover.sqlite"
	isDefaultDBPathExists := func() bool {
		_, err := os.Stat(defaultDBPath)
		return err == nil
	}

	var fourcloverDatabaseSqlitePath string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverDatabaseSqlitePathEnv, _ := GetSystemEnvVar("FOURCLOVER_DATABASE_SQLITE_PATH")
	if yamlConfig.Database.Sqlite.Path != "" && fourcloverDatabaseSqlitePathEnv == "" {
		fourcloverDatabaseSqlitePath = yamlConfig.Database.Sqlite.Path
	} else if fourcloverDatabaseSqlitePathEnv != "" {
		fourcloverDatabaseSqlitePath = fourcloverDatabaseSqlitePathEnv
	} else if isDefaultDBPathExists() {
		fourcloverDatabaseSqlitePath = defaultDBPath
	} else {
		log.Default().Println("ERROR: No fourclover database sqlite path specified in the config file or system environment variable. Creating a new database file in the current directory.")
		fourcloverDatabaseSqlitePath = defaultDBPath
	}

	return fourcloverDatabaseSqlitePath
}

// ----------------------- Server Configurations ----------------------- //

// Get fourclover server host and port
func GetFourCloverServerHostAndPort() (string, int) {
	if !GetFourCloverStartServer() {
		log.Default().Println("WARN: fourclover server host and port is not needed because save_to_server is set to false.")
		return "", 0
	}

	var fourcloverServerHost string
	var fourcloverServerPort int
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}

	StringToInt := func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
			return 0
		}
		return i
	}

	fourcloverServerHostEnv, _ := GetSystemEnvVar("FOURCLOVER_SERVER_HOST")
	fourcloverServerPortEnv, _ := GetSystemEnvVar("FOURCLOVER_SERVER_PORT")

	if yamlConfig.Server.Host != "" && fourcloverServerHostEnv == "" {
		fourcloverServerHost = yamlConfig.Server.Host
	} else if fourcloverServerHostEnv != "" {
		fourcloverServerHost = fourcloverServerHostEnv
	} else {
		log.Fatalf("ERROR: No fourclover server host specified in the config file or system environment variable.")
		return "", 0
	}

	if yamlConfig.Server.Port != 0 && fourcloverServerPortEnv == "" {
		fourcloverServerPort = yamlConfig.Server.Port
	} else if fourcloverServerPortEnv != "" {
		fourcloverServerPort = StringToInt(fourcloverServerPortEnv)
	} else {
		log.Fatalf("ERROR: No fourclover server port specified in the config file or system environment variable.")
		return "", 0
	}

	return fourcloverServerHost, fourcloverServerPort
}

// Get fourclover server username and password
func GetFourCloverServerUsernameAndPassword() (string, string) {
	if !GetFourCloverStartServer() {
		log.Default().Println("WARN: fourclover server username and password is not needed because save_to_server is set to false.")
		return "", ""
	}

	var fourcloverServerUsername string
	var fourcloverServerPassword string
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}

	fourcloverServerUsernameEnv, _ := GetSystemEnvVar("FOURCLOVER_SERVER_USERNAME")
	fourcloverServerPasswordEnv, _ := GetSystemEnvVar("FOURCLOVER_SERVER_PASSWORD")

	if yamlConfig.Server.Username != "" && fourcloverServerUsernameEnv == "" {
		fourcloverServerUsername = yamlConfig.Server.Username
	} else if fourcloverServerUsernameEnv != "" {
		fourcloverServerUsername = fourcloverServerUsernameEnv
	} else {
		log.Fatalf("ERROR: No fourclover server username specified in the config file or system environment variable.")
		return "", ""
	}

	if yamlConfig.Server.Password != "" && fourcloverServerPasswordEnv == "" {
		fourcloverServerPassword = yamlConfig.Server.Password
	} else if fourcloverServerPasswordEnv != "" {
		fourcloverServerPassword = fourcloverServerPasswordEnv
	} else {
		log.Fatalf("ERROR: No fourclover server password specified in the config file or system environment variable.")
		return "", ""
	}

	return fourcloverServerUsername, fourcloverServerPassword
}

// ----------------------- Logger Configurations ----------------------- //

// Get fourclover logger path
func GetFourCloverLoggerPath() string {
	if GetFourCloverSuppressLogs() {
		log.Default().Println("WARN: fourclover logger path is not needed because logs are suppressed as requested.")
		return ""
	}

	var fourcloverLoggerPath string
	var defaultLoggerPath = "logs/fourclover_" + time.Now().Format("2006-01-02_15-04-05") + ".log"
	var yamlConfig YamlConfig
	if YamlConfigurationStatus() {
		yamlConfig, _ = GetYamlConfigData()
	}
	fourcloverLoggerPathEnv, _ := GetSystemEnvVar("FOURCLOVER_LOGGER_PATH")
	if yamlConfig.Logger.Path != "" && fourcloverLoggerPathEnv == "" {
		fourcloverLoggerPath = yamlConfig.Logger.Path
	} else if fourcloverLoggerPathEnv != "" {
		fourcloverLoggerPath = fourcloverLoggerPathEnv
	} else {
		log.Default().Println("INFO: No fourclover logger path specified in the config file or system environment variable. Using default path:", defaultLoggerPath)
		fourcloverLoggerPath = defaultLoggerPath
	}

	return fourcloverLoggerPath
}
