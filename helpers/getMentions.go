package helpers

import "strings"

func GetMentions(text string) []string {
	var mentions []string
	for _, word := range strings.Fields(text) {
		if strings.HasPrefix(word, "@") {
			mentions = append(mentions, word)
		}
	}
	return mentions
}
