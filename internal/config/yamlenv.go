package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	Version string `yaml:"version"` // 0.5.0

	Global struct {
		IsSnapshot           bool `yaml:"is_snapshot"`             // true
		IsCompare            bool `yaml:"is_compare"`              // false
		IsPolicy             bool `yaml:"is_policy"`               // false
		SaveReportToDatabase bool `yaml:"save_report_to_database"` // false
		SaveReportToFile     bool `yaml:"save_report_to_file"`     // true
		StartServer          bool `yaml:"start_server"`            // false
		SuppressLogs         bool `yaml:"suppress_logs"`           // false
	} `yaml:"global"`

	Snapshot struct {
		Name             string   `yaml:"name"`               // scan_2021-01-01_00-00-00
		Directory        string   `yaml:"directory"`          // .\example\
		HashAlgorithms   []string `yaml:"hash_algorithms"`    // [sha256, sha512]
		ExcludeThem      []string `yaml:"exclude_them"`       // [node_modules, .git]
		FocusExtensions  []string `yaml:"focus_extensions"`   // [.js, .css, .html]
		ExportReportPath string   `yaml:"export_report_path"` // zzz/snapshot_report_2021-01-01_00-00-00.json
	} `yaml:"snapshot"`

	Compare struct {
		OldReportPath    string `yaml:"old_report_path"`    // zzz/snapshot_report_2021-01-01_00-00-00.json
		NewReportPath    string `yaml:"new_report_path"`    // zzz/snapshot_report_2021-01-01_00-00-00.json
		ExportReportPath string `yaml:"export_report_path"` // zzz/compare_report_2021-01-01_00-00-00.json
	} `yaml:"compare"`

	Policy struct {
		SnapshotReportPath string `yaml:"snapshot_report_path"`  // zzz/snapshot_report_2021-01-01_00-00-00.json
		TargetDirPath      string `yaml:"target_directory_path"` // .\example\
		DirPath            string `yaml:"policy_directory_path"` // .\example\policies\
		ExportReportPath   string `yaml:"export_report_path"`    // zzz/policy_report_2021-01-01_00-00-00.json
	} `yaml:"policy"`

	Database struct {
		Sqlite struct {
			Path string `yaml:"path"` // fourclover.sqlite
		} `yaml:"sqlite"`
	} `yaml:"database"`

	Server struct {
		Host     string `yaml:"host"`     // localhost
		Port     int    `yaml:"port"`     // 8080
		Username string `yaml:"username"` // admin
		Password string `yaml:"password"` // admin
	} `yaml:"server"`

	Logger struct {
		Path string `yaml:"path"` // fourclover.log
	} `yaml:"logger"`
}

// Read YAML config file, if it exists return true, else return false
func YamlConfigurationStatus() bool {
	// Read YAML file
	_, err := ioutil.ReadFile("fourclover-config.yaml")
	if err != nil {
		return false
	}

	return true
}

// GetYamlConfigData returns the YAML config data
func GetYamlConfigData() (YamlConfig, error) {
	// Read YAML file
	configData, err := ioutil.ReadFile("fourclover-config.yaml")
	if err != nil {
		log.Default().Println("ERROR: error reading fourclover-config YAML file. Skipping...")
		return YamlConfig{}, err
	}

	// Parse YAML
	var yamlConfig YamlConfig
	err = yaml.Unmarshal(configData, &yamlConfig)
	if err != nil {
		log.Default().Println("ERROR: error parsing fourclover-config YAML:", err)
		return YamlConfig{}, err
	}

	// Return YAML config data
	return yamlConfig, nil
}
