/*
 *    Copyright (c) 2025 [renegade@renegade-master.com]
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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

// --- //

type WeightedSimpleSearchResult struct {
	Result string
	Weight int
	Id     string
}

type BySimpleSearchResult func(wpi1, wpi2 *WeightedSimpleSearchResult) bool

func (by BySimpleSearchResult) SortSimpleSearchResult(results []WeightedSimpleSearchResult) {
	wpiSorter := &weightedSimpleSearchResultSorter{
		weightedSimpleSearchResults: results,
		by:                          by,
	}

	sort.Sort(wpiSorter)
}

type weightedSimpleSearchResultSorter struct {
	weightedSimpleSearchResults []WeightedSimpleSearchResult
	by                          func(wpi1, wpi2 *WeightedSimpleSearchResult) bool
}

func (t weightedSimpleSearchResultSorter) Len() int {
	return len(t.weightedSimpleSearchResults)
}

func (t weightedSimpleSearchResultSorter) Less(i, j int) bool {
	return t.by(&t.weightedSimpleSearchResults[i], &t.weightedSimpleSearchResults[j])
}

func (t weightedSimpleSearchResultSorter) Swap(i, j int) {
	t.weightedSimpleSearchResults[i], t.weightedSimpleSearchResults[j] = t.weightedSimpleSearchResults[j], t.weightedSimpleSearchResults[i]
}
