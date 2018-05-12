package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	// Major version when you make incompatible API changes
	Major int = iota
	// Minor version when you add functionality in a backwards-compatible manner
	Minor
	// Patch version when you make backwards-compatible bug fixes
	Patch
)

// Version represents a SemVer version i.e. 1.0.0
type Version struct {
	Major int
	Minor int
	Patch int
}

// Set takes in a SemVer version string and splits it into parts
func (version *Version) Set(v string) (err error) {
	parts := strings.Split(v, ".")

	version.Major, err = strconv.Atoi(parts[Major])
	version.Minor, err = strconv.Atoi(parts[Minor])
	version.Patch, err = strconv.Atoi(parts[Patch])

	return err
}

// Update will take in a version section and increment it
func (version *Version) Update(section string) error {
	r := reflect.ValueOf(version).Elem().FieldByName(section)
	if r.IsValid() {
		current := reflect.ValueOf(version).Elem().FieldByName(section).Int()
		r.SetInt(current + 1)

		return nil
	}

	return fmt.Errorf("invalid section [%s] passed to update", section)
}

// String will implement the stringer interface, and return the SemVer versioning string
func (version Version) String() string {
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
}

func main() {
	version := flag.String("version", "patch", "the semver version type to increment (major|minor|patch)")
	directory := flag.String("dir", "", "the absolute path to the directory to run the tagging update within, defaults to current")

	flag.Parse()

	commandSequence := fmt.Sprintf("cd %s;git tag", *directory)
	cmd := exec.Command("/bin/sh", "-c", commandSequence)

	output, _ := cmd.CombinedOutput()

	r, _ := regexp.Compile("[\\d\\.]{5}")
	matches := r.FindAllString(string(output), -1)

	currentVersion := "0.0.0"
	if len(matches) > 0 {
		currentVersion = matches[len(matches)-1]
	}

	v := Version{}
	err := v.Set(currentVersion)
	if err != nil {
		panic(err)
	}

	err = v.Update(*version)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Project [%s] will be updated to Version [%s], is this correct? (y/n)\n", *directory, v)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if input != "y\n" {
		os.Exit(0)
	}

	commandSequence = fmt.Sprintf("cd %s;git tag -a \"%s\" -m \"Release %s\";", *directory, v, v)
	cmd = exec.Command("/bin/sh", "-c", commandSequence)
	cmd.Run()

	fmt.Printf("Tags updated locally, would you like to push to remote? (y/n)")

	input, err = reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if input != "y\n" {
		os.Exit(0)
	}

	commandSequence = fmt.Sprintf("cd %s;git push origin --tags", *directory)
	cmd = exec.Command("/bin/sh", "-c", commandSequence)
	cmd.Run()
}
