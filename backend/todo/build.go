package todo

import (
	"encoding/json"
	"fmt"
)

var Build = struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}{
	Version: "dev",
	Commit:  "NA",
	Date:    "NA",
}

func BuildDetails() string {
	return fmt.Sprintf("Build Details:\n\tVersion:\t%s\n\tCommit:\t\t%s\n\tDate:\t\t%s", Build.Version, Build.Commit, Build.Date)
}

func BuildJSON() []byte {
	b, _ := json.Marshal(Build)
	return b
}
