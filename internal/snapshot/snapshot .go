package snapshot

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/zeebo/blake3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"

	common "fourclover.org/internal/common"
	report "fourclover.org/internal/report"
)

// type File struct {
// 	Name         string `json:"name"`
// 	Path         string `json:"path"`
// 	Size         int64  `json:"size"`
// 	Permission   string `json:"permission"`
// 	LastModified string `json:"last_modified"`

// 	Blake2b_256 string `json:"blake2b-256,omitempty"`
// 	Blake2b_512 string `json:"blake2b-512,omitempty"`
// 	Blake3      string `json:"blake3,omitempty"`
// 	CRC32       string `json:"crc32,omitempty"`
// 	MD5         string `json:"md5,omitempty"`
// 	SHA1        string `json:"sha1,omitempty"`
// 	SHA3_224    string `json:"sha3-224,omitempty"`
// 	SHA3_256    string `json:"sha3-256,omitempty"`
// 	SHA3_384    string `json:"sha3-384,omitempty"`
// 	SHA3_512    string `json:"sha3-512,omitempty"`
// 	SHA256      string `json:"sha256"`
// 	SHA512      string `json:"sha512,omitempty"`
// }

var nilReport report.SnapshotReport

func SnapshotDirectory(dir string, hashAlgorithms common.HashAlgorithms, excludeThem common.ExcludeThem, focusExtensions common.FocusExtensions, scanName string) (report.SnapshotReport, error) {
	log.Default().Println("INFO: Scanning directory", dir)
	var report report.SnapshotReport
	report.FourCloverVersion = common.APP_VERSION // fourclover version
	report.Date = time.Now().Format(time.RFC3339)
	Working_Dir, _ := filepath.Abs(dir)
	report.WorkingDir = filepath.ToSlash(Working_Dir)
	report.Name = scanName

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		var pathInfo os.FileInfo

		pathInfo, _ = os.Stat(path)

		if err != nil {
			return err
		}

		// Exclude directories and files from the scan if they are in the excludeThem list
		if len(excludeThem) > 0 {
			for _, excludeThem := range excludeThem {
				if (pathInfo.Name() == excludeThem) || (filepath.Ext(path) == excludeThem) {
					if pathInfo.IsDir() { // Exclude directories
						log.Default().Println("	↳ WARN: Excluding directory", pathInfo.Name())
						return filepath.SkipDir
					}
					if !pathInfo.IsDir() { // Exclude files or files with specific extensions
						log.Default().Println("	↳ WARN: Excluding file", pathInfo.Name())
						return nil
					}
				}
			}
		}

		// Print some verbose output
		if pathInfo.IsDir() {
			log.Default().Println("INFO: Scanning", path)
		}

		if pathInfo.IsDir() { // TODO: Fix Known to identify some files as directories
			return nil
		}

		// Check if the file extension is in the focusExtensions list, if not skip the file
		if len(focusExtensions) > 0 {
			var found bool
			for _, focusExtension := range focusExtensions {
				if focusExtension == ".*" {
					found = true
					break
				}
				if filepath.Ext(path) == focusExtension {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Additional file metadata
		fileSize := info.Size()                            // File size in bytes eg: 1234
		filePermissions := info.Mode()                     // File mode bits eg: -rw-r--r--
		fileModtime := info.ModTime().Format(time.RFC3339) // Last modification time eg: 2021-03-01T12:00:00+01:00

		// Calculate checksums of the file based on the hash algorithms
		// blake2b-256,blake2b-512,blake3,crc32,md5,sha1,sha3-224,sha3-256,sha3-384,sha3-512,sha256,sha512

		file := common.File{
			Name:         info.Name(),
			Path:         filepath.ToSlash(path),
			Size:         fileSize,
			Permission:   filePermissions.String(),
			LastModified: fileModtime,
		}
		var hash []byte
		for _, HashAlgorithms := range hashAlgorithms {
			switch HashAlgorithms {
			case "md5":
				tempHash := md5.Sum(data)
				hash = tempHash[:]
				file.MD5 = hex.EncodeToString(hash)
			case "sha1":
				tempHash := sha1.Sum(data)
				hash = tempHash[:]
				file.SHA1 = hex.EncodeToString(hash)
			case "sha256":
				tempHash := sha256.Sum256(data)
				hash = tempHash[:]
				file.SHA256 = hex.EncodeToString(hash)
			case "sha512":
				tempHash := sha512.Sum512(data)
				hash = tempHash[:]
				file.SHA512 = hex.EncodeToString(hash)
			case "blake2b-256":
				tempHash, _ := blake2b.New256(nil)
				tempHash.Write(data)
				hash = tempHash.Sum(nil)
				file.Blake2b_256 = hex.EncodeToString(hash)
			case "blake2b-512":
				tempHash, _ := blake2b.New512(nil)
				tempHash.Write(data)
				hash = tempHash.Sum(nil)
				file.Blake2b_512 = hex.EncodeToString(hash)
			case "blake3":
				tempHash := blake3.Sum256(data)
				hash = tempHash[:]
				file.Blake3 = hex.EncodeToString(hash)
			case "sha3-224":
				tempHash := sha3.Sum224(data)
				hash = tempHash[:]
				file.SHA3_224 = hex.EncodeToString(hash)
			case "sha3-256":
				tempHash := sha3.Sum256(data)
				hash = tempHash[:]
				file.SHA3_256 = hex.EncodeToString(hash)
			case "sha3-384":
				tempHash := sha3.Sum384(data)
				hash = tempHash[:]
				file.SHA3_384 = hex.EncodeToString(hash)
			case "sha3-512":
				tempHash := sha3.Sum512(data)
				hash = tempHash[:]
				file.SHA3_512 = hex.EncodeToString(hash)
			case "crc32":
				tempHash := crc32.ChecksumIEEE(data)
				hash = []byte(fmt.Sprintf("%x", tempHash))
				file.CRC32 = string(hash)
			default:
				allAlgorithms := "List of supported algorithms: blake2b-256, blake2b-512, blake3, crc32, md5, sha1, sha3-224, sha3-256, sha3-384, sha3-512, sha256, sha512"
				return fmt.Errorf("unknown hash algorithm: %s \n%s", HashAlgorithms, allAlgorithms)
			}
		}
		report.Files = append(report.Files, file)

		// Calculate the total number of files in the report
		report.TotalFiles += 1
		// Calculate the total size of the files in the report
		report.TotalFilesSize += fileSize
		return nil
	})

	if err != nil {
		return nilReport, err
	}

	sort.Slice(report.Files, func(i, j int) bool {
		return report.Files[i].Name < report.Files[j].Name
	})

	// Calculate checksum of the fields in "files": [] in the report
	var checksum string
	for _, file := range report.Files {
		checksum += file.Name + file.Path + file.Permission + file.LastModified + file.SHA256 + report.Date + report.FourCloverVersion
	}
	reportHash := sha256.Sum256([]byte(checksum))
	// Convert the checksum to a string
	reportHashString := hex.EncodeToString(reportHash[:])

	// Encrypt the checksum with AES 256 and store it in the report
	if encrypted, err := common.Encrypt(common.CIPHER_KEY, reportHashString); err != nil {
		return nilReport, fmt.Errorf("error encrypting report checksum: %s", err)
	} else {
		report.ReportChecksum = encrypted
	}

	return report, nil
}
