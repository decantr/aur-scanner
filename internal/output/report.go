package output

type Report struct {
	Source                  string   `json:"source"`
	CheckedPackageCount     int      `json:"checked_package_count"`
	CompromisedPackageCount int      `json:"compromised_package_count"`
	Matches                 []string `json:"matches"`
	Safe                    bool     `json:"safe"`
}

type TextOptions struct {
	Quiet bool
	Color bool
}
