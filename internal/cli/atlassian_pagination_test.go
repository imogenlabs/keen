package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

// fakeAtlassian serves `total` items, `pageCap` per page, in Agile's
// {startAt, maxResults, isLast, values} shape (or with total instead of isLast).
type fakeAtlassian struct {
	total, pageCap int
	useTotal       bool
	calls          int
}

func (f *fakeAtlassian) GetWithHeaders(_ context.Context, _ string, params map[string]string, _ map[string]string) (json.RawMessage, error) {
	f.calls++
	start, _ := strconv.Atoi(params["startAt"])
	end := min(start+f.pageCap, f.total)
	values := []map[string]int{}
	for i := start; i < end; i++ {
		values = append(values, map[string]int{"id": i})
	}
	page := map[string]any{"startAt": start, "maxResults": f.pageCap, "values": values}
	if f.useTotal {
		page["total"] = f.total
	} else {
		page["isLast"] = end >= f.total
	}
	return json.Marshal(page)
}

func TestPaginatedGetFollowsAtlassianOffsets(t *testing.T) {
	for _, useTotal := range []bool{false, true} {
		t.Run(fmt.Sprintf("total=%v", useTotal), func(t *testing.T) {
			f := &fakeAtlassian{total: 123, pageCap: 50, useTotal: useTotal}
			// maxResults 100 is asked for; the server caps pages at 50.
			params := map[string]string{"startAt": "0", "maxResults": "100", "state": "closed"}
			data, err := paginatedGet(context.Background(), f, "/rest/agile/1.0/board/1/sprint", params, nil, true, "", "offset", "maxResults", "", "")
			if err != nil {
				t.Fatal(err)
			}
			var items []map[string]int
			if err := json.Unmarshal(data, &items); err != nil {
				t.Fatal(err)
			}
			if len(items) != 123 || f.calls != 3 {
				t.Fatalf("got %d items in %d calls, want 123 in 3", len(items), f.calls)
			}
			for i, it := range items {
				if it["id"] != i {
					t.Fatalf("item %d has id %d: pages overlapped or skipped", i, it["id"])
				}
			}
		})
	}
}

func TestPaginatedGetWithoutAllStaysOnePage(t *testing.T) {
	f := &fakeAtlassian{total: 123, pageCap: 50}
	if _, err := paginatedGet(context.Background(), f, "/x", map[string]string{"startAt": "0"}, nil, false, "", "offset", "maxResults", "", ""); err != nil {
		t.Fatal(err)
	}
	if f.calls != 1 {
		t.Fatalf("calls = %d, want 1", f.calls)
	}
}
