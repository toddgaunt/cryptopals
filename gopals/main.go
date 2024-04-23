package main

import (
	"errors"
	"fmt"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	clear = "\033[0m"
)

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

func decodeHex(input string) ([]byte, error) {
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

func encodeHex(input []byte) string {
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

func Challenge1() error {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	want := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	raw, err := decodeHex(input)
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

	x, _ := decodeHex(a)
	y, _ := decodeHex(b)
	z, _ := xorBuf(x, y)

	if got, want := encodeHex(z), c; got != want {
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
}
