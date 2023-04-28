package npm

import (
	"fmt"
	"strconv"
	"strings"
)

type InvalidNPMPackageNameError struct {
	Message string
}

func (e InvalidNPMPackageNameError) Error() string {
	return fmt.Sprintf("Invalid npm package name: %s", e.Message)
}

func IsNPMPackageNameValid(name string) error {
	if len(name) < 1 || len(name) > 214 {
		return InvalidNPMPackageNameError{"Package name length must be between 1 and 214 characters"}
	}

	// check valid scope format, if present
	name, err := checkValidScopeFormat(name)
	if err != nil {
		return err
	}

	// check first character is lowercase letter
	if name[0] < 'a' || name[0] > 'z' {
		return InvalidNPMPackageNameError{"Package name must start with a lowercase letter"}
	}

	err2 := validateCharacters(name)
	if err2 != nil {
		return err2
	}

	return nil
}

func validateCharacters(name string) error {
	for i := 0; i < len(name); i++ {
		c := name[i]

		// check character is valid
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return InvalidNPMPackageNameError{"Package name must contain only lowercase letters, numbers, hyphens (-), and underscores (_)"}
		}

		// check no consecutive hyphens or underscores
		if i > 0 && (c == '_' || c == '-') && name[i-1] == c {
			return InvalidNPMPackageNameError{"Package name must not contain consecutive hyphens (-) or underscores (_)"}
		}

		// check no leading/trailing hyphens or underscores
		if (i == 0 || i == len(name)-1) && (c == '_' || c == '-') {
			return InvalidNPMPackageNameError{"Package name must not start or end with a hyphen (-) or underscore (_)"}
		}
	}
	return nil
}

func checkValidScopeFormat(name string) (string, error) {
	if name[0] == '@' {
		scopeEnd := -1
		for i := 1; i < len(name); i++ {
			if name[i] == '/' {
				scopeEnd = i
				break
			}
			if !((name[i] >= 'a' && name[i] <= 'z') || (name[i] >= '0' && name[i] <= '9') || name[i] == '-') {
				return "", InvalidNPMPackageNameError{"Package scope must contain only lowercase letters, numbers, and hyphens (-)"}
			}
		}
		if scopeEnd == -1 {
			return "", InvalidNPMPackageNameError{"Scoped package name must contain a slash (/) after the scope"}
		}
		if scopeEnd == 1 {
			return "", InvalidNPMPackageNameError{"Scoped package name must contain a non-empty scope"}
		}
		if scopeEnd == len(name)-1 {
			return "", InvalidNPMPackageNameError{"Scoped package name must contain a package name after the scope"}
		}
		name = name[scopeEnd+1:]
	}
	return name, nil
}

type InvalidNPMVersionError struct {
	Message string
}

func (e InvalidNPMVersionError) Error() string {
	return fmt.Sprintf("Invalid npm version: %s", e.Message)
}

func IsNPMVersionValid(version string) error {
	if version == "" {
		return InvalidNPMVersionError{"Version cannot be empty"}
	}
	if version == "latest" {
		return nil
	}

	// check valid version format
	parts := strings.Split(version, ".")
	if len(parts) < 1 || len(parts) > 3 {
		return InvalidNPMVersionError{"Version must be in the format X.Y.Z or X.Y or X, where X, Y, and Z are integers"}
	}

	for _, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return InvalidNPMVersionError{"Version must be in the format X.Y.Z or X.Y or X, where X, Y, and Z are integers"}
		}
	}

	return nil
}
