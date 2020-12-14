package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindStringInArray(t *testing.T) {
	fruits := []string{"apple", "lime", "papaya", "plum", "tangerine"}

	// If array contains given string - non-negative index must be returned
	const fruit = "plum"
	result := FindStringInArray(fruit, fruits)
	assert.Greater(t, result, -1, fmt.Sprintf("String '%s' should be found in the given array %v", fruit, fruits))

	// If array does not contain given string, -1 must be returned
	const vegetable = "pepper"
	result = FindStringInArray(vegetable, fruits)
	assert.Equal(t, -1, result, fmt.Sprintf("String '%s' should not be found in the given array %v", fruit, fruits))
}

func TestFindStringInArrayByPrefix(t *testing.T) {
	countries := []string{"Estonia", "Japan", "Montenegro", "Portugal", "Serbia"}

	// If array contains a string with given prefix:
	// - non-negative index must be returned
	// - the returned index must correspond to the argument position in the original array
	prefix := "Mon"
	countryIndex := FindStringInArrayByPrefix(prefix, countries)
	assert.Greater(t, countryIndex, -1, fmt.Sprintf("String with prefix '%s' should be found in the given array %v", prefix, countries))
	assert.Equal(t, 2, countryIndex, fmt.Sprintf("String with prefix '%s' should have index %d in the given array %v", prefix, countryIndex, countries))

	// If array does not contain a string with given prefix - -1 must be returned
	prefix = "Fin"
	countryIndex = FindStringInArrayByPrefix(prefix, countries)
	assert.Equal(t, -1, countryIndex, fmt.Sprintf("String with prefix '%s' should not be found in the given array %v", prefix, countries))
}

func TestFindStringDuplicates(t *testing.T) {
	// If array contains duplicates - an array of duplicates must be returned
	mountains := []string{"Everest", "K2", "Skil Brum", "Elbrus", "Everest", "Kazbek", "Kazbek"}
	duplicates := []string{"Everest", "Kazbek"}
	foundDuplicates := FindStringDuplicates(mountains)
	assert.Equal(t, duplicates, foundDuplicates, fmt.Sprintf("Array %v must have duplicates", mountains))

	// If array does not contain duplicates - an empty array must be returned
	rivers := []string{"Nile", "Amazon", "Mississippi", "Yenisei", "Huang He"}
	duplicates = FindStringDuplicates(rivers)
	assert.Equal(t, 0, len(duplicates), fmt.Sprintf("Array %v must not have duplicates", rivers))
}

func TestRemoveStringDuplicates(t *testing.T) {
	// If array contains duplicates - an array of all elements except duplicates must be returned
	mountains := []string{"Everest", "K2", "Skil Brum", "Elbrus", "Everest", "Kazbek", "Kazbek"}
	uniqueMountains := []string{"Everest", "K2", "Skil Brum", "Elbrus", "Kazbek"}
	arrayWithNoDuplicates := RemoveStringDuplicates(mountains)
	assert.Equal(t, uniqueMountains, arrayWithNoDuplicates, fmt.Sprintf("Result array %v must have duplicates", mountains))

	// If array does not contain duplicates - same array must be returned
	rivers := []string{"Nile", "Amazon", "Mississippi", "Yenisei", "Huang He"}
	arrayWithNoDuplicates = RemoveStringDuplicates(rivers)
	assert.Equal(t, rivers, arrayWithNoDuplicates, fmt.Sprintf("Array %v must not have duplicates", rivers))
}

func TestContainsOneOfStrings(t *testing.T) {
	const initialString = "Go hang a salami, I'm a lasagna hog"
	const existingSubstring = "salam"
	const nonexistentSubstring = "aleykum"

	// If string contains one of provided substrings - 'true' must be returned
	result := ContainsOneOfStrings(initialString, existingSubstring, nonexistentSubstring)
	assert.True(t, result, fmt.Sprintf("String \"%s\" must contain a substring \"%s\"", initialString, existingSubstring))

	// If string does not contain any of provided substrings - 'false' must be returned
	result = ContainsOneOfStrings(initialString, nonexistentSubstring, nonexistentSubstring)
	assert.False(t, result, fmt.Sprintf("String \"%s\" must not contain a substring \"%s\"", initialString, nonexistentSubstring))
}

func TestIsAnyOfStringsInArray(t *testing.T) {
	array := []string{"Go", "hang", "a", "salami", "I'm", "a", "lasagna", "hog"}
	const existingElement = "salami"
	const nonexistentElement = "aleykum"

	// If string array includes one of provided elements - 'true' must be returned
	result := IsAnyOfStringsInArray(array, existingElement, nonexistentElement)
	assert.True(t, result, fmt.Sprintf("Array %v must contain an element \"%s\"", array, existingElement))

	// If string array does not include any of provided elements - 'false' must be returned
	result = IsAnyOfStringsInArray(array, nonexistentElement, nonexistentElement)
	assert.False(t, result, fmt.Sprintf("String \"%s\" must not contain an element \"%s\"", array, nonexistentElement))
}
