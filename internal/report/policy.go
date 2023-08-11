package report

import (
	"encoding/json"
	"os"
	"time"
)

// Classification struct represents the CVSS classification of a finding
type Classification struct {
	CvssScore  float64 `json:"cvss_score"`
	CvssVector string  `json:"cvss_vector"`
}

// Findings struct represents a policy finding
type Finding struct {
	RuleID         string           `json:"rule_id"`
	Name           string           `json:"name"`
	Type           string           `json:"type"`
	Description    string           `json:"description"`
	FileHash       string           `json:"file_hash"`
	FilePath       string           `json:"file_path"`
	References     []string         `json:"references"`
	Classification []Classification `json:"classification"`
	Tags           []string         `json:"tags"`
}

// Report struct represents a policy report
type PolicyReport struct {
	Title      string    `json:"title"`
	Date       string    `json:"date"`
	WorkingDir string    `json:"working_dir"`
	Summary    string    `json:"summary"`
	Findings   []Finding `json:"findings"`
}

func PreparePolicyFinding(ruleid string, name string, ruleType string, description string, fileHash string, filePath string, references []string, classificationCvssScore float64, classificationMetrics string, tags []string) Finding {

	classification := []Classification{}
	classification = append(classification, Classification{CvssScore: classificationCvssScore, CvssVector: classificationMetrics})

	finding := Finding{
		RuleID:         ruleid,
		Name:           name,
		Type:           ruleType,
		Description:    description,
		FileHash:       fileHash,
		FilePath:       filePath,
		References:     references,
		Classification: classification,
		Tags:           tags,
	}

	return finding
}

// preparePolicyReport prepares a policy report
func preparePolicyReport(WorkingDir string, reportTitle string, reportSummary string, findings []Finding) PolicyReport {
	report := PolicyReport{
		Title:      reportTitle,
		Date:       time.Now().Format(time.RFC3339),
		WorkingDir: WorkingDir,
		Summary:    reportSummary,
		Findings:   findings,
	}

	return report
}

func SavePolicyReport(allFindings []Finding, WorkingDir string, policyExportReportPath string) error {
	reportName := "Policy Report"
	reportSummary := "Findings reported in the policy report are based on the policy rules defined in the policy file."
	report := preparePolicyReport(WorkingDir, reportName, reportSummary, allFindings)

	jsonData, err := json.Marshal(report)
	if err != nil {
		return err
	}

	file, err := os.Create(policyExportReportPath)
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
