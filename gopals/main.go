package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	clear = "\033[0m"
)

var asciiTable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// encodeBase64 encodes a slice of bytes into a base64 string. Padding is not
// supported.
func encodeBase64(input []byte) ([]byte, error) {
	if n := len(input); n%3 != 0 {
		return nil, errors.New("input length must be divisible by 3")
	}

	var output []byte
	for i := 0; i < len(input); i += 3 {
		// Chunk is a container for 4 encoded 6-bit values, 24 bits or 3 8-bit
		// bytes total.
		var chunk [4]byte

		// First 6 bits
		chunk[0] = ((input[i+0] & 0xFC) >> 2)
		// 2 bits + 4 bits = 6 bits
		chunk[1] = ((input[i+0] & 0x03) << 4) | ((input[i+1] & 0xF0) >> 4)
		// 4 bits + 2 bits = 6 bits
		chunk[2] = ((input[i+1] & 0x0F) << 2) | ((input[i+2] & 0xC0) >> 6)
		// Last 6 bits
		chunk[3] = (input[i+2] & 0x3F)

		output = append(output, base64Table[chunk[0]])
		output = append(output, base64Table[chunk[1]])
		output = append(output, base64Table[chunk[2]])
		output = append(output, base64Table[chunk[3]])
	}

	return output, nil
}

func hexDecode(input string) ([]byte, error) {
	if n := len(input); n%2 != 0 {
		return nil, errors.New("input length must be divisible by 2")
	}

	var runes = []rune(input)
	var output []byte
	for i := 0; i < len(input); i += 2 {
		n1, ok := hex2byte(runes[i])
		if !ok {
			return nil, fmt.Errorf("invalid hex character at index %d", i)
		}
		n2, ok := hex2byte(runes[i+1])
		if !ok {
			return nil, fmt.Errorf("invalid hex character at index %d", i+1)
		}
		n := (n1 << 4) | n2
		output = append(output, n)
	}

	return output, nil
}

func hex2byte(c rune) (byte, bool) {
	if c >= '0' && c <= '9' {
		return byte(c - '0'), true
	} else if c >= 'a' && c <= 'f' {
		return byte(c - 'a' + 0xa), true
	} else if c >= 'A' && c <= 'F' {
		return byte(c - 'A' + 0xa), true
	}

	return 0, false
}

func byte2hex(n byte) []rune {
	n1 := (n & 0xF0) >> 4
	n2 := n & 0x0F

	result := [2]rune{0, 0}

	if n1 <= 0x9 {
		result[0] = rune(n1 + byte('0'))
	} else if 0xa <= n1 && n1 <= 0xf {
		result[0] = rune(n1 + byte('a') - 0xa)
	}

	if n2 <= 0x9 {
		result[1] = rune(n2 + byte('0'))
	} else if 0xa <= n2 && n2 <= 0xf {
		result[1] = rune(n2 + byte('a') - 0xa)
	}

	return result[:]
}

func hexEncode(input []byte) string {
	var output []rune
	for i := 0; i < len(input); i++ {
		output = append(output, byte2hex(input[i])...)
	}

	return string(output)
}

func xorBuf(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, errors.New("buffers must be of equal length")
	}

	n := len(a)

	var output = make([]byte, n)
	for i := 0; i < n; i++ {
		output[i] = a[i] ^ b[i]
	}

	return output, nil
}

func xorByte(s []byte, b byte) []byte {
	var output []byte
	for i := 0; i < len(s); i++ {
		output = append(output, s[i]^b)
	}
	return output
}

func Challenge1() error {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	want := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	raw, err := hexDecode(input)
	if err != nil {
		return err
	}

	got, err := encodeBase64(raw)
	if err != nil {
		return err
	}

	if string(got) != want {
		return fmt.Errorf("got %s, want %s", got, want)
	}

	return nil
}

func Challenge2() error {
	a := "1c0111001f010100061a024b53535009181c"
	b := "686974207468652062756c6c277320657965"
	c := "746865206b696420646f6e277420706c6179"

	x, _ := hexDecode(a)
	y, _ := hexDecode(b)
	z, _ := xorBuf(x, y)

	if got, want := hexEncode(z), c; got != want {
		return fmt.Errorf("got %s, want %s", got, want)
	}

	return nil
}

func scoreText1(text []byte) int {
	n := 0
	list := "ETAOIN SHRDLU"
	for _, c := range text {
		if strings.ContainsRune(list, rune(c)) || strings.ContainsRune(strings.ToLower(list), rune(c)) {
			n++
		}
	}

	return n
}

func Challenge3() error {
	score := scoreText1

	// This message has been encrypted with a single character xor
	enc := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	raw, err := hexDecode(enc)
	if err != nil {
		return err
	}

	key := 0
	highScore := 0
	for i := 0; i < 256; i++ {
		dec := xorByte(raw, byte(i))
		s := score(dec)
		if s > highScore {
			highScore = s
			key = i
		}
	}

	got := string(xorByte(raw, byte(key)))
	want := "Cooking MC's like a pound of bacon"
	if got != want {
		return fmt.Errorf("got %s, want %s", got, want)
	}

	return nil
}

func readLines(filename string) ([]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		lines = append(lines, line)
	}

	return lines, nil
}

func Challenge4() error {
	scoreFn := scoreText1

	lines, err := readLines("4.txt")
	if err != nil {
		return err
	}

	type score struct {
		key   byte
		score int
	}

	rawLines := [][]byte{}
	for _, enc := range lines {
		raw, err := hexDecode(enc)
		if err != nil {
			return err
		}
		rawLines = append(rawLines, raw)
	}

	highestRank := ""
	rankedLines := map[string]score{}
	for i := range rawLines {
		var p score
		for k := 0; k < 256; k++ {
			key := byte(k)
			dec := xorByte(rawLines[i], key)
			s := scoreFn(dec)
			if s > p.score {
				p = score{key: key, score: s}
			}
			if s > rankedLines[highestRank].score {
				highestRank = lines[i]
			}
		}
		rankedLines[lines[i]] = p
	}

	winner := rankedLines[highestRank]
	raw, err := hexDecode(highestRank)
	if err != nil {
		return err
	}

	decoded := xorByte(raw, winner.key)

	/*
	// Print results
	for line, score := range(rankedLines) {
		fmt.Printf("%s: %v\n", line, score)
	}
	fmt.Printf("Winner: score=%d %s=%s\n", winner.score, highestRank, decoded)
	*/

	if got, want := string(decoded), "Now that the party is jumping\n"; got != want {
		return fmt.Errorf("got %s, want %s", got, want)
	}

	return nil
}

func printResult(n int, err error) {
	fmt.Printf("Challenge %d: ", n)
	if err != nil {
		fmt.Printf("%s%v%s\n", red, err, clear)
	} else {
		fmt.Printf("%sOK%s\n", green, clear)
	}
}

func main() {
	var err error

	fmt.Println("Gopals!")

	err = Challenge1()
	printResult(1, err)

	err = Challenge2()
	printResult(2, err)

	err = Challenge3()
	printResult(3, err)

	err = Challenge4()
	printResult(4, err)
}
