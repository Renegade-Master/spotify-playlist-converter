package youtube

import (
	"sort"

	"google.golang.org/api/youtube/v3"
)

type WeightedPlaylistItem struct {
	Result *youtube.PlaylistItem
	Weight int
}

type ByPlaylistItem func(wpi1, wpi2 *WeightedPlaylistItem) bool

func (by ByPlaylistItem) SortPlaylistItem(results []WeightedPlaylistItem) {
	wpiSorter := &weightedPlaylistItemSorter{
		weightedPlaylistItems: results,
		by:                    by,
	}

	sort.Sort(wpiSorter)
}

type weightedPlaylistItemSorter struct {
	weightedPlaylistItems []WeightedPlaylistItem
	by                    func(wpi1, wpi2 *WeightedPlaylistItem) bool
}

func (t weightedPlaylistItemSorter) Len() int {
	return len(t.weightedPlaylistItems)
}

func (t weightedPlaylistItemSorter) Less(i, j int) bool {
	return t.by(&t.weightedPlaylistItems[i], &t.weightedPlaylistItems[j])
}

func (t weightedPlaylistItemSorter) Swap(i, j int) {
	t.weightedPlaylistItems[i], t.weightedPlaylistItems[j] = t.weightedPlaylistItems[j], t.weightedPlaylistItems[i]
}

// --- //

type WeightedSearchResult struct {
	Result *youtube.SearchResult
	Weight int
}

type BySearchResult func(wpi1, wpi2 *WeightedSearchResult) bool

func (by BySearchResult) SortSearchResult(results []WeightedSearchResult) {
	wpiSorter := &weightedSearchResultSorter{
		weightedSearchResults: results,
		by:                    by,
	}

	sort.Sort(wpiSorter)
}

type weightedSearchResultSorter struct {
	weightedSearchResults []WeightedSearchResult
	by                    func(wpi1, wpi2 *WeightedSearchResult) bool
}

func (t weightedSearchResultSorter) Len() int {
	return len(t.weightedSearchResults)
}

func (t weightedSearchResultSorter) Less(i, j int) bool {
	return t.by(&t.weightedSearchResults[i], &t.weightedSearchResults[j])
}

func (t weightedSearchResultSorter) Swap(i, j int) {
	t.weightedSearchResults[i], t.weightedSearchResults[j] = t.weightedSearchResults[j], t.weightedSearchResults[i]
}
