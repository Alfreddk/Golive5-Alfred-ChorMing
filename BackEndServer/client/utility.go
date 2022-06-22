package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// List of messages
var (
	InvalidKeyEntry   = errors.New("Invalid Key Entry")
	ExcessiveKeyEntry = errors.New("Excessive Key Entry")
	KeyNotFound       = errors.New("Key Entry Not Found")
	InvalidCapacity   = errors.New("Invalid Capacity Entry")
)

const numSetStr string = "0123456789" // define the character set for numeric

// Get a string from bufio.Reader
// Remove leading and trailing blanks
// convert to string, check len of string
// Possible returns for (str, err)
// ("", nil)     empty read (just carriage return)
// (string, nil) read with message string
// ("", err != nil) error reading
func readString() (string, error) {

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n') // read until delimited (delimiter included)

	if err == nil {
		line = strings.TrimSpace(line) // remove leading and trailing spaces

		if len(line) == 0 { // check for blank entry
			//			return "", errors.New("Empty Input String")
			return "", nil
		}
		return line, nil
	} else {
		return "", errors.New("Read String Error")
	}
}

// returns the index of the match in the character Set of 1 character
// return err != nil, value -1 if no match
// return err != nil, value > 0, N is the number of characters received
// return err = nil, value = -1 means empty string
// otherwise return err = nil, value = index of the match character

func get1KeyMatch(charSet string) (int, error) {

	keyStr, err := readString()
	if err != nil {
		return -1, InvalidKeyEntry
	}

	// Empty string always treated as valid
	keylen := len(keyStr)
	if keylen == 0 {
		return -1, nil
	}

	// Non empty string but more than 1 character
	if keylen > 1 {
		return keylen, ExcessiveKeyEntry // return value = len of the string
	}

	// search the 1 char against the character set
	i := strings.Index(charSet, keyStr)
	if i >= 0 {
		return i, nil
	}
	// char not found in the character set
	return i, KeyNotFound
}

// returns the number of match of N input characters tested contiguously are within Character Set
// return
// value -1, err != nil if at least one of the n char is not in the character set
// value -2, err != nil, if more than n characters entered
// value 0,  err = nil, if empty character is entered
// value n,  err = nil, if at most n characters are found in character set
/// otherwise return the number of matched characters found
func checkAtMostNKeyMatch(keyStr string, charSet string, n int) (int, error) {

	// keyStr, err := readString()
	// if err != nil {
	// 	return -1, InvalidKeyEntry
	// }

	keylen := len(keyStr)
	// entered string len = 0, meaning empty entry
	// return 0 for index (indicating a empty match)
	if keylen == 0 {
		return 0, nil
	}

	// entered string is longer than n
	if keylen > n {
		return -2, ExcessiveKeyEntry
	}

	// found n characters
	// Test every input character against Character Set
	i := 0
	for _, v := range keyStr {
		indexMatch := strings.Index(charSet, (string)(v))
		if indexMatch == -1 {
			return -1, InvalidKeyEntry // at least one of the n character is not in the Character set
		}
		i++
	}
	return i, nil // returns the number of match characters (should be the same as n)
}

// color for display
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

//setColor set the color of the display
func setColor(color string) {
	fmt.Printf("%s", color)
}
