package utils

import (
	"archive/tar"
	"github.com/hardboiledalex/go-tools/lib/logging"
	"bufio"
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"time"
)

var (
	CurrentHostname, _ = os.Hostname()
	currentFQDN        string
	DevMode            bool
	EulaAccepted       bool
)

func GetCurrentFQDN() string {
	if currentFQDN == "" {
		out, err := exec.Command("sh", "-c", "hostname --fqdn").CombinedOutput()
		if err != nil {
			logging.LogError("Cannot obtain FQDN", err.Error())
			os.Exit(1)
		}
		currentFQDN = strings.TrimSuffix(string(out[:]), "\n")
	}

	return currentFQDN
}

func GetDirectories(dir string) []string {
	file, err := os.Open(dir)
	if err != nil {
		logging.LogErrorf("Failed to open directory: %s", dir)
		os.Exit(1)
	}
	defer CloseFile(file)

	var directories []string
	list, _ := file.Readdir(0)
	for _, file := range list {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}
	return directories
}

func GetTarFilesInDir(dir string) []string {
	file, err := os.Open(dir)
	if err != nil {
		logging.LogErrorf("Failed to open directory: %s", dir)
		os.Exit(1)
	}
	defer CloseFile(file)

	var tarFiles []string
	list, _ := file.Readdir(0)
	for _, file := range list {
		if ContainsOneOfStrings(file.Name(), ".tar", ".tar.gz", ".tgz") {
			tarFiles = append(tarFiles, file.Name())
		}
	}
	return tarFiles
}

func GetFilesInDirByPrefix(dir string, prefix string) ([]string, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer CloseFile(file)

	var filesFound []string
	list, _ := file.Readdir(0)
	for _, file := range list {
		if ContainsOneOfStrings(file.Name(), prefix) {
			filesFound = append(filesFound, file.Name())
		}
	}
	return filesFound, nil
}

func GetFilesInDirByPrefixes(dir string, prefixes []string) []string {
	file, err := os.Open(dir)
	if err != nil {
		logging.LogErrorf("Failed to open directory: %s", dir)
		os.Exit(1)
	}
	defer CloseFile(file)

	var filesFound []string
	list, _ := file.Readdir(0)
	for _, file := range list {
		if !file.IsDir() {
			for _, prefix := range prefixes {
				if ContainsOneOfStrings(file.Name(), prefix) {
					filesFound = append(filesFound, file.Name())
				}
			}
		}
	}
	return filesFound
}

func GetDirContentAsStrings(dir string) []string {
	var filesFound []string
	fileInfoList := GetDirContent(dir)
	for _, file := range fileInfoList {
		filesFound = append(filesFound, file.Name())
	}

	return filesFound
}

func GetDirContent(dir string) []os.FileInfo {
	file, err := os.Open(dir)
	if err != nil {
		logging.LogErrorf("Failed to open directory: %s", dir)
		os.Exit(1)
	}
	defer CloseFile(file)

	list, err := file.Readdir(0)
	if err != nil {
		logging.LogErrorf("Cannot retrieve content of %s\n", dir)
		os.Exit(1)
	}
	return list
}

func FindStringInArray(str string, arr []string) int {
	for i, s := range arr {
		if s == str {
			return i
		}
	}
	return -1
}

func FindStringInArrayByPrefix(prefix string, arr []string) int {
	for i, str := range arr {
		if strings.HasPrefix(str, prefix) {
			return i
		}
	}
	return -1
}

func FindStringDuplicates(arr []string) []string {
	uniqueSet := make(map[string]struct{})
	var duplicates []string
	for _, item := range arr {
		_, found := uniqueSet[item]
		if found {
			duplicates = append(duplicates, item)
		} else {
			uniqueSet[item] = struct{}{}
		}
	}
	return duplicates
}

func RemoveStringDuplicates(slice []string) []string {
	uniqueSet := make(map[string]bool)
	var sliceWithoutDuplicates []string
	for _, entry := range slice {
		if _, value := uniqueSet[entry]; !value {
			uniqueSet[entry] = true
			sliceWithoutDuplicates = append(sliceWithoutDuplicates, entry)
		}
	}
	return sliceWithoutDuplicates
}

func ContainsOneOfStrings(sourceString string, strs ...string) bool {
	for _, s := range strs {
		if strings.Contains(sourceString, s) {
			return true
		}
	}
	return false
}

func IsAnyOfStringsInArray(args []string, strs ...string) bool {
	for _, arg := range args {
		for _, s := range strs {
			if strings.EqualFold(s, arg) {
				return true
			}
		}
	}
	return false
}

func ContainsOneOfSuffixes(sourceString string, strs ...string) bool {
	for _, s := range strs {
		if strings.HasSuffix(sourceString, s) {
			return true
		}
	}
	return false
}

func ContainsOneOfPrefixes(sourceString string, strs ...string) bool {
	for _, s := range strs {
		if strings.HasPrefix(sourceString, s) {
			return true
		}
	}
	return false
}

func IfArgumentsAllowed(args []string, allowedArgs []string) {
	for _, arg := range args {
		if !ifArgumentAllowed(arg, allowedArgs) {
			logging.LogErrorf("Not allowed argument %s \n ", arg)
			os.Exit(1)
		}
	}
}

func ifArgumentAllowed(arg string, allowedArgs []string) bool {
	for _, aarg := range allowedArgs {
		if aarg == arg {
			return true
		}
	}
	return false
}

func ReadFromFile(filePath string) string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logging.LogError(err)
		os.Exit(1)
	}
	return string(bytes)
}

func GetLines(content string) []string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)
	var arr []string
	for scanner.Scan() {
		line := scanner.Text()
		arr = append(arr, line)
	}
	return arr
}

func ExtractSubstring(original string, prefix string, suffix string) string {
	substring := strings.Replace(original, prefix, "", 1)
	substring = strings.Replace(substring, suffix, "", 1)
	return substring
}

func OpenFile(path string) *os.File {
	logging.LogInfo("Opening file ", path)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logging.LogErrorf("Could not open file %s", path)
		os.Exit(1)
	}
	return f
}

func AddLineToFile(file *os.File, line string) {
	_, err := file.WriteString(line)
	if err != nil {
		logging.LogErrorf("Unable to add line to %s", file.Name())
		os.Exit(1)
	}
}

func CloseFile(file *os.File) {
	logging.LogTracef("Closing file %s", file.Name())
	if file != nil {
		err := file.Close()
		if err != nil {
			logging.LogErrorf("Cannot close file: %s", file.Name())
		}
	}
}

/*
Finding a file by filename inside tar or tar.gz archives and return its content as string
*/
func ReadFileContentFromTar(pathToTar string, regexpPattern string) (string, error) {
	file, err := os.Open(pathToTar)
	archive, err := gzip.NewReader(file)
	if err != nil {
		if IsProgressRunning() {
			StopProgress()
		}
		logging.LogErrorf("Cannot open tar file %s", pathToTar)
		os.Exit(1)
	}
	var content string
	tarReader := tar.NewReader(archive)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			logging.LogError(err)
			os.Exit(1)
		}
		matched, err := regexp.MatchString(regexpPattern, strings.ToLower(header.Name))
		if matched {
			bytes, _ := ioutil.ReadAll(tarReader)
			content = string(bytes)
		}
	}
	if content != "" {
		return content, nil
	} else {
		return "", fmt.Errorf("The file matching pattern \"%s\" was not found in the %s.", regexpPattern, pathToTar)
	}
}

// GetHomeDirectory obtains user home directory regardless of the host system
func GetHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		logging.LogError(err)
		os.Exit(1)
	}
	return currentUser.HomeDir
}

func ResolveIPByHostname(hostname string) (string, error) {
	addr, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}
	return addr[0].String(), nil
}

func ResolveHostnameByIp(ipAddress string) string {
	addr, err := net.LookupAddr(ipAddress)
	if err != nil {
		logging.LogErrorf("Invalid IP address %s", ipAddress)
		os.Exit(1)
	}
	hostname := strings.TrimRight(addr[0], ".")
	return hostname
}

// ReadFileToString reads all file content to string
func ReadFileToString(filePath string) string {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logging.LogErrorf("Failed to read file: %s", filePath)
		os.Exit(1)
	}

	return string(fileBytes[:])
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func BytesToString(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGT"[exp])
}

// Converts binary data to the array of strings
func BytesToStrings(binaryData []byte) []string {
	return strings.Split(string(binaryData[:]), "\n")
}

// Converts array of strings to the binary representation
func StringsToBytes(stringArray []string) []byte {
	return []byte(strings.Join(stringArray, "\n"))
}

// Encloses string in double quotes
func QuoteString(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func CreateMapFromSlice(input []string, keyValueSeparator string) map[string]string {
	result := make(map[string]string)
	for _, line := range input {
		items := strings.Split(line, keyValueSeparator)
		result[strings.Trim(items[0], " ")] = strings.Join(items[1:], keyValueSeparator)
	}
	return result
}

var progressBar *spinner.Spinner

func InitSpinner() {
	progressBar = spinner.New(spinner.CharSets[36], 100*time.Millisecond, spinner.WithWriter(os.Stdout))
	_ = progressBar.Color("green", "bold")
	// Please keep following snippet. It is useful for spinner debugging.
	//progressBar.PostUpdate = func(s *spinner.Spinner) {
	//	logging.LogMessageDirect(fmt.Sprintf("Spinner post-update %v", time.Now()))
	//}
}

func IsProgressRunning() bool {
	return progressBar.Active()
}

const progressDelay = 50

func StartProgress() {
	progressBar.Start()
	time.Sleep(progressDelay * time.Millisecond)
}

func StopProgress() {
	progressBar.Stop()
	time.Sleep(progressDelay * time.Millisecond)
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
