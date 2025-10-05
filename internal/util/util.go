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

package util

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/zmb3/spotify/v2"
)

const MaxDistance int = 5

var badPhrases = []string{
	"\\[", "\\]", "\\(", "\\)", "-", ",",
	"clip",
	"hd",
	"lyrics",
	"video",
	"with",
	"official",
}

// CheckIfCommandExists checks if executable 'e' is in PATH
func CheckIfCommandExists(e ...string) bool {
	anyCommandFound := false

	for _, command := range e {
		_, err := exec.LookPath(command)

		if err != nil {
			anyCommandFound = true
			break
		}
	}

	return anyCommandFound
}

func RemoveIndexString(original []string, index int) []string {
	log.Printf("Removing Item [%d] [%s]", index, original[index])

	modified := make([]string, 0)
	modified = append(modified, original[:index]...)

	return append(modified, original[index+1:]...)
}

func RemoveIndexTrack(original []spotify.PlaylistTrack, index int) []spotify.PlaylistTrack {
	log.Printf("Removing Item [%d] [%s]", index, original[index].Track.Name)

	modified := make([]spotify.PlaylistTrack, 0)
	modified = append(modified, original[:index]...)

	return append(modified, original[index+1:]...)
}

func LevenshteinDistance(stringA, stringB string) int {
	cleanedA := strings.ToLower(stringA)
	cleanedB := strings.ToLower(stringB)

	cleanedA = cleanTitle(cleanedA)
	cleanedB = cleanTitle(cleanedB)

	distance := levenshtein.ComputeDistance(cleanedA, cleanedB)
	log.Printf("Levenshtein Distance between strings \n[%s] and \n[%s]\nis [%d]", cleanedA, cleanedB, distance)

	return distance
}

func cleanTitle(original string) string {
	newString := original

	for _, replacement := range badPhrases {

		reString := fmt.Sprintf("(?i)%s", replacement)
		re := regexp.MustCompile(reString)

		newString = re.ReplaceAllString(newString, "")
	}

	reString := "(\\s){2,}"
	re := regexp.MustCompile(reString)

	newString = re.ReplaceAllString(newString, " ")

	return strings.Trim(newString, " .")
}
