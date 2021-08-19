package sbot

import (
	"fmt"

	refs "go.mindeco.de/ssb-refs"
)

// stub for disabled mode
type disabledMetaFeeds struct{}

var errMetafeedsDisabled = fmt.Errorf("sbot: metafeeds are disabled")

func (disabledMetaFeeds) CreateSubFeed(mount refs.FeedRef, purpose string, format refs.RefAlgo) (refs.FeedRef, error) {
	return refs.FeedRef{}, errMetafeedsDisabled
}

func (disabledMetaFeeds) TombstoneSubFeed(_, _ refs.FeedRef) error {
	return errMetafeedsDisabled
}

func (disabledMetaFeeds) ListSubFeeds(mount refs.FeedRef) ([]SubfeedListEntry, error) {
	return nil, errMetafeedsDisabled
}

func (disabledMetaFeeds) Publish(as refs.FeedRef, content interface{}) (refs.MessageRef, error) {
	return refs.MessageRef{}, errMetafeedsDisabled
}
