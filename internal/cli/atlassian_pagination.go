// Local customization: --all follows Atlassian startAt/isLast paging.
// See .printing-press-patches.json (atlassian-offset-pagination).

package cli

import (
	"encoding/json"
	"strconv"
)

// atlassianNextStartAt reports the next startAt for Jira/Agile/Service Desk
// style offset paging, where the spec declares no cursor or has-more field
// but the response carries isLast (or total). It advances by the number of
// items actually returned, because Atlassian caps maxResults server-side.
func atlassianNextStartAt(obj map[string]json.RawMessage, params map[string]string, got int) (string, bool) {
	if got == 0 {
		return "", false
	}
	start := 0
	if s := params["startAt"]; s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return "", false
		}
		start = n
	}
	next := start + got

	if raw, ok := obj["isLast"]; ok {
		var last bool
		if json.Unmarshal(raw, &last) != nil || last {
			return "", false
		}
		return strconv.Itoa(next), true
	}
	if raw, ok := obj["total"]; ok {
		var total int
		if json.Unmarshal(raw, &total) != nil || next >= total {
			return "", false
		}
		return strconv.Itoa(next), true
	}
	return "", false
}

func rawHas(obj map[string]json.RawMessage, keys ...string) bool {
	for _, k := range keys {
		if _, ok := obj[k]; ok {
			return true
		}
	}
	return false
}
