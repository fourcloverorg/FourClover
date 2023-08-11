package common

import (
	"log"
	"strings"
)

// Custom types for flags in snapshot
type DirectoryToSnapshot string
type HashAlgorithms []string
type ExcludeThem []string
type FocusExtensions []string

// Custom types for flags in static policy check
type SnapshotFile string
type PolicyDir string
type TargetDir string

// DirectoryToSnapshot is a custom type that implements the flag.Value interface.
// It is used to parse a single directory to scan. if the user passes multiple directories, they will receive an error message.
func (d *DirectoryToSnapshot) Set(value string) error {
	err := ValidateDirectory(value)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
		return nil
	}
	*d = DirectoryToSnapshot(value)
	return nil
}

func (d *DirectoryToSnapshot) String() string {
	return string(*d)
}

// HashAlgorithms is a custom type that implements the flag.Value interface.
// It is used to parse the comma-separated list of hash algorithms to include in the scan.
func (h *HashAlgorithms) Set(value string) error {
	*h = append(*h, strings.Split(value, ",")...)
	return nil
}

func (h *HashAlgorithms) String() string {
	return strings.Join(*h, ",")
}

// excludeThem is a custom type that implements the flag.Value interface.
// It is used to parse the comma-separated list of directories and files to exclude from the scan.
func (e *ExcludeThem) Set(value string) error {
	*e = append(*e, strings.Split(value, ",")...)
	return nil
}

func (e *ExcludeThem) String() string {
	return strings.Join(*e, ",")
}

// FileExtension is a custom type that implements the flag.Value interface.
// It is used to parse the comma-separated list of file extensions to include in the scan.
func (f *FocusExtensions) Set(value string) error {
	*f = append(*f, strings.Split(value, ",")...)
	return nil
}

func (f *FocusExtensions) String() string {
	return strings.Join(*f, ",")
}

// SnapshotFile is a custom type that implements the flag.Value interface.
// It is used to parse the snapshot file to use for the scan.
func (s *SnapshotFile) Set(value string) error {
	err := ValidateFile(value)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
		return nil
	}
	*s = SnapshotFile(value)
	return nil
}

func (s *SnapshotFile) String() string {
	return string(*s)
}

// PolicyDir is a custom type that implements the flag.Value interface.
// It is used to parse the directory containing the policy files.
func (p *PolicyDir) Set(value string) error {
	err := ValidateDirectory(value)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
		return nil
	}
	*p = PolicyDir(value)
	return nil
}

func (p *PolicyDir) String() string {
	return string(*p)
}

// TargetDir is a custom type that implements the flag.Value interface.
// It is used to parse the directory to scan.
func (t *TargetDir) Set(value string) error {
	err := ValidateDirectory(value)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
		return nil
	}
	*t = TargetDir(value)
	return nil
}

func (t *TargetDir) String() string {
	return string(*t)
}
