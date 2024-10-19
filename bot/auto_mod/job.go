package auto_mod

import (
	dgo "github.com/bwmarrin/discordgo"
	"log/slog"
	"regexp"
)

func (a *AutoMod) ExecuteBackgroundJob() error {
	a.logger.Debug("executing background job")
	return a.syncRules()
}

func (a *AutoMod) syncRules() error {
	rules, err := a.session.AutoModerationRules(a.guildId.String())
	if err != nil {
		return err
	}

	a.emptyCachedRules()
	a.extractRulesMetadata(rules)

	a.logger.Debug("rules synced",
		slog.Int("keywordsFilter", len(a.keywordFilter)),
		slog.Int("regexPatterns", len(a.regexPatterns)),
	)

	return nil
}

func (a *AutoMod) extractRulesMetadata(rules []*dgo.AutoModerationRule) {
	for _, rule := range rules {
		if rule.Enabled == nil || !*rule.Enabled {
			continue
		}

		if rule.EventType != dgo.AutoModerationEventMessageSend {
			continue
		}

		a.processTriggerMetadata(rule.TriggerMetadata)
	}
}

// Function to convert a word with wildcards into a regex pattern
func (a *AutoMod) convertToRegexPattern(word string) string {
	//// Escape special characters, and replace "*" with ".*" to match any characters
	//pattern := "^" + regexp.QuoteMeta(word) + "$"
	//// Replace escaped asterisk (\*) with regex equivalent (.*)
	//pattern = regexp.MustCompile(`\\\*`).ReplaceAllString(pattern, ".*")

	// Escape special characters, and replace "*" with ".*" to match any characters
	pattern := regexp.QuoteMeta(word) // Escape regex special chars in word
	// Replace escaped asterisk (\*) with regex equivalent (.*) to match any characters
	pattern = regexp.MustCompile(`\\\*`).ReplaceAllString(pattern, ".*")
	return pattern
}

func (a *AutoMod) processTriggerMetadata(metadata *dgo.AutoModerationTriggerMetadata) {
	if metadata == nil {
		return
	}

	a.processTriggerMetadataKeywordFilter(metadata.KeywordFilter)

	a.processTriggerMetadataRegexPatterns(metadata.RegexPatterns)
}

func (a *AutoMod) processTriggerMetadataKeywordFilter(filters []string) {
	if len(filters) == 0 {
		return
	}

	for _, filter := range filters {
		pattern := a.convertToRegexPattern(filter)
		//regex, err := regexp.Compile(pattern)
		//if err != nil {
		//	a.logger.Error("keyword filter parsing failed",
		//		slog.Any("err", filter),
		//	)
		//	continue
		//}

		a.keywordFilter[filter] = pattern
	}
}

func (a *AutoMod) processTriggerMetadataRegexPatterns(patterns []string) {
	if len(patterns) == 0 {
		return
	}

	for _, pattern := range patterns {
		a.regexPatterns[pattern] = pattern
	}
}
