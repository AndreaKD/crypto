package main

import (
	"fmt"
)

const hex = "0123456789abcdef"
const base64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var hexRunes = []rune(hex)
var base64Runes = []rune(base64)

const inputAsHex = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
const expectedOutput = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

func main() {
	fmt.Println("\nHello!!\n")

	hexToBase64()

	tests()

	fmt.Println("\nThe end!!")

}

func hexToBase64() {
	//convert ASCII to hex bytes (numbers 0-15)
	//convert hex bytes to 4 digit binary word
	//repeat and concatenate to produce raw binary data

	var rawData = make([]byte, 0)
	var input = []rune(inputAsHex)
	for i := 1; i <= len(input)/2; i++ {
		rhsPos := (i * 2) - 1
		lhsPos := rhsPos - 1
		hexByte := hexToByte(hexRuneTo4BitWord(input[lhsPos]), hexRuneTo4BitWord(input[rhsPos]))

		slice := make([]byte, 1)
		slice[0] = hexByte
		rawData = append(rawData, slice...)
	}

	// fmt.Println("raw data: ", rawData)
	fmt.Println("raw data as string: ", string(rawData))

	//chop data into 6 bit words
	//convert 6 bit words into base 64 bytes (numbers 0-65)
	//convert base64 bytes to ASCII
	var base64words = make([]byte, 0)
	// var base64Runes = make([]rune, 0)
	// base64 string will be
	for i := 0; i < len(rawData); i = i + 3 {

		var threeBytesOfRawData [3]byte

		copy(threeBytesOfRawData[:], rawData[i:(i+3)]) //convert slice to fixed size array
		var fourBytesOf6BitsEach [4]byte = bytesTo6BitWords(threeBytesOfRawData)
		var fourBytes = fourBytesOf6BitsEach[:] //convert fixed size array to slice
		base64words = append(base64words, fourBytes...)
	}
	// fmt.Println("base64 bytes: ", base64words)
	var base64 = make([]rune, 0)
	for i := 0; i < len(base64words); i++ {
		base64 = append(base64, base64RuneFromBase64word(base64words[i]))
	}

	fmt.Println("expected as string:     [", expectedOutput, "]")
	fmt.Println("base64 bytes as string: [", string(base64), "]")
}

//returns 4 bit word from hex rune
func hexRuneTo4BitWord(input rune) byte {
	var numericalVal = int(0)
	for i, num := range hexRunes {
		if input == num {
			numericalVal = int(i)
		}
	}
	return byte(numericalVal)
}

//returns 8 bits from 2 4 bit words
func hexToByte(lhs byte, rhs byte) byte {

	var lhsShifted byte = lhs << 4
	var result byte = lhsShifted + rhs
	// fmt.Println(strconv.FormatInt(int64(result), 2))
	return result
}

func base64RuneFromBase64word(base64word byte) rune {
	return base64Runes[base64word]
}

//chops up a chunk of raw data into 6 bit words (3 bytes at a time)
func bytesTo6BitWords(rawData [3]byte) [4]byte {
	// fmt.Println("            : ", strconv.FormatInt(int64(rawData[0]), 2), strconv.FormatInt(int64(rawData[1]), 2), strconv.FormatInt(int64(rawData[2]), 2))
	byte1 := rawData[0] >> 2                                       //first 6 bits of rawData[0]
	byte2 := ((rawData[0] &^ 0b11111100) << 4) + (rawData[1] >> 4) //last 2 bits of rawData[0] and first 4 of [2]
	byte3 := ((rawData[1] &^ 0b11110000) << 2) + (rawData[2] >> 6) //last 4 bits of rawData[1] and first 2 of [2]
	byte4 := rawData[2] &^ 0b11000000                              //remove first 2 bits from [2]

	// fmt.Println("converted to: ", strconv.FormatInt(int64(byte1), 2), strconv.FormatInt(int64(byte2), 2), strconv.FormatInt(int64(byte3), 2), strconv.FormatInt(int64(byte4), 2))
	// fmt.Println("            : ", [4]byte{byte1, byte2, byte3, byte4})
	return [4]byte{byte1, byte2, byte3, byte4}
}

func tests() {
	fmt.Println("\nTESTS\ntestbytesTo6BitWords")
	testbytesTo6BitWords()
	fmt.Println("testHexToByte")
	testHexToByte()
}

func testbytesTo6BitWords() {

	input1 := [3]byte{0b01001101, 0b11100010, 0b10010101}
	expectedOutput1 := [4]byte{0b010011, 0b011110, 0b001010, 0b010101}
	result1 := bytesTo6BitWords(input1)
	assertequals(expectedOutput1[:], result1[:])

	input2 := [3]byte{0b01001001, 0b00100111, 0b01101101}
	expectedOutput2 := [4]byte{0b010010, 0b010010, 0b011101, 0b101101}
	result2 := bytesTo6BitWords(input2)
	assertequals(expectedOutput2[:], result2[:])

	input3 := [3]byte{0b01110010, 0b01100001, 0b01101001}
	expectedOutput3 := [4]byte{0b011100, 0b100110, 0b000101, 0b101001}
	result3 := bytesTo6BitWords(input3)
	assertequals(expectedOutput3[:], result3[:])
}

func testHexToByte() {
	result1 := [1]byte{hexToByte(0b00000100, 0b00001001)}
	expectedOutput1 := [1]byte{0b01001001}
	assertequals(expectedOutput1[:], result1[:])
}

func assertequals(expected []byte, actual []byte) {
	var pass bool = true
	if len(expected) != len(actual) {
		pass = false
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			pass = false
		}
	}
	if pass {
		fmt.Println(" PASS")
	} else {
		fmt.Println(" FAIL expected: ", expected, " but was: ", actual)
	}
}
