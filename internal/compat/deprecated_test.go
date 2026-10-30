package compat

import "testing"

// A removed endpoint must still be one the specs list as deprecated. When
// Atlassian finally drops it from the specs, its commands disappear on the next
// regen: move it to RetiredCommands then.
func TestRemovedEndpointsAreDeprecated(t *testing.T) {
	for key := range RemovedEndpoints {
		if _, ok := DeprecatedEndpoints[key]; !ok {
			t.Errorf("%s is in RemovedEndpoints but not in the specs' deprecated list; move its commands to RetiredCommands", key)
		}
	}
}

func TestEndpointKey(t *testing.T) {
	if got := EndpointKey("GET", "/rest/api/3/issue/{issueIdOrKey}/comment/{id}"); got != "GET /rest/api/3/issue/{}/comment/{}" {
		t.Fatalf("EndpointKey = %q", got)
	}
}
