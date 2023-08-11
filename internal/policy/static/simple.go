package policy

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"fourclover.org/internal/report"
	"github.com/antonmedv/expr"
	"github.com/go-yaml/yaml"
)

type Vulnerability struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type SimpleRule struct {
	RuleID string `yaml:"rule_id"`
	Meta   struct {
		Type           string   `yaml:"type"`
		Name           string   `yaml:"name"`
		Author         string   `yaml:"author"`
		Severity       string   `yaml:"severity"`
		Description    string   `yaml:"description"`
		Reference      []string `yaml:"reference"`
		Classification struct {
			CvssScore   float64 `yaml:"cvss-score"`
			CvssMetrics string  `yaml:"cvss-metrics"`
		} `yaml:"classification"`
		Tags []string `yaml:"tags"`
	} `yaml:"meta"`
	Simple struct {
		Patterns  map[string]string `yaml:"patterns"`
		Condition string            `yaml:"condition"`
	} `yaml:"simple"`
}

type Rules struct {
	Rules []SimpleRule `yaml:"rules"`
}

type CustomEnv struct {
	Patterns        map[string]string
	FilePath        string
	Condition       string
	PatternsMatched int
}

func replaceWholeWord(expression, target, replacement string) string {
	regex := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(target) + `\b`)
	return regex.ReplaceAllString(expression, replacement)
}

func SimpleStaticPolicyCheck(policyFile string, targetFile string) ([]report.Finding, error) {

	targetFileHash := strings.Split(targetFile, ":")[0]
	targetFilePath := strings.Split(targetFile, ":")[1]
	// Load the YAML rules file
	yamlFile, err := ioutil.ReadFile(policyFile)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// Parse the YAML rules into a struct
	var rules Rules
	err = yaml.Unmarshal(yamlFile, &rules)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Get the code from the file
	codeFile, err := ioutil.ReadFile(targetFilePath)
	if err != nil {
		log.Fatalf("Failed to read code file: %v", err)
	}
	// Convert the code to a string
	code := string(codeFile)

	for _, rule := range rules.Rules {
		// Replace the pattern references in the condition with array indexing
		currentCondition := strings.ReplaceAll(rule.Simple.Condition, "patterns", "Patterns")

		// Check if the code matches any of the patterns
		patternsMatched := 0
		for index, pattern := range rule.Simple.Patterns {
			// Ignore / at the start and end of the pattern
			if pattern[0] == '/' {
				pattern = pattern[1:]
			}
			if pattern[len(pattern)-1] == '/' {
				pattern = pattern[:len(pattern)-1]
			}

			// Match the pattern against the code
			matched, err := regexp.MatchString(pattern, code)
			if err != nil {
				log.Fatalf("Error matching pattern for rule '%s': %v", rule.RuleID, err)
			}
			if matched {
				patternsMatched++
				// Replace the index in the condition with true
				currentCondition = replaceWholeWord(currentCondition, index, "true")
			} else {
				// Replace the index in the condition with false
				currentCondition = replaceWholeWord(currentCondition, index, "false")
			}
		}

		// Prepare the environment
		env := CustomEnv{
			Patterns:        rule.Simple.Patterns,
			Condition:       currentCondition,
			FilePath:        targetFilePath,
			PatternsMatched: patternsMatched,
		}

		// Compile the condition into an expression program
		program, err := expr.Compile(currentCondition, expr.Env(env))
		if err != nil {
			log.Fatalf("Error compiling condition for rule '%s': %v", rule.RuleID, err)
		}

		// Evaluate the condition using the compiled program and environment
		output, err := expr.Run(program, env)
		if err != nil {
			log.Fatalf("Error evaluating condition for rule '%s': %v", rule.RuleID, err)
		}

		result, ok := output.(bool)
		if !ok {
			log.Fatalf("Condition evaluation for rule '%s' did not produce a boolean result", rule.RuleID)
		}

		if result {
			log.Printf("INFO: Rule '%s' matched file '%s'", rule.RuleID, targetFilePath)

			finding := report.PreparePolicyFinding(rule.RuleID, rule.Meta.Name, rule.Meta.Type, rule.Meta.Description, targetFileHash, targetFilePath, rule.Meta.Reference, rule.Meta.Classification.CvssScore, rule.Meta.Classification.CvssMetrics, rule.Meta.Tags)

			return []report.Finding{finding}, nil
		}
	}

	return nil, nil
}
