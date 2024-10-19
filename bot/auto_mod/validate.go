package auto_mod

import (
	"errors"
	"regexp"
)

func (a *AutoMod) ValidateUserInputAgainstServerRules(input string) error {
	if input == "" {
		return nil
	}

	err := a.matchKeywordFilters(input)
	if err != nil {
		return err
	}

	return a.matchRegexPatterns(input)
}

func (a *AutoMod) matchKeywordFilters(input string) error {
	for keyword, regexPattern := range a.keywordFilter {
		matched, _ := regexp.MatchString(regexPattern, input)
		if matched {
			return errors.New(keyword)
		}
	}

	return nil
}

func (a *AutoMod) matchRegexPatterns(input string) error {
	for regexName, regexPattern := range a.regexPatterns {
		matched, _ := regexp.MatchString(regexPattern, input)
		if matched {
			return errors.New(regexName)
		}
	}
	return nil
}
