package main

import (
	"fmt"
)

type lzTuple struct {
	pos int
	len int
	c   byte
}

func (t *lzTuple) String() string {
	return fmt.Sprintf("p: %v, l: %v, c: %s", t.pos, t.len, string(t.c))
}

func main() {
	var input string
	fmt.Println("Type a string to be encoded:")
	fmt.Scanln(&input)
	fmt.Printf("Your string: %s\n", input)

	// Encode
	ts := encode(input)
	fmt.Println("Encoded:")
	for i, t := range ts {
		fmt.Printf("%v: %s\n", i+1, t)
	}

	// Decode
	output := decode(ts)
	fmt.Printf("Decoded: %s\n", output)
	if output != input {
		fmt.Println("Compression failed")
	}

	// The above works, but the compression rate is awful lol. The next step is
	// to store the tuples space-efficiently.
}

func encode(in string) []*lzTuple {
	ts := make([]*lzTuple, 0)
	index := 0
	for index < len(in) {
		t := lzTuple{}
		pos, length := match(in[:index], in[index:])
		var c byte
		if index+length < len(in) {
			c = in[index+length]
		}
		t.pos, t.len, t.c = pos, length, c
		index += length + 1
		ts = append(ts, &t)
	}
	return ts
}

func match(search, lookAhead string) (pos, length int) {
	if search == "" {
		return
	}

	// This is definitely not efficient, but works for now. i is for counting backwards in
	// the already checked "search" string, while j is for counting forwards in the look ahead
	// string.
	for i, j := len(search)-1, 0; i >= 0 && j < len(search); {
		if i+j < len(search) && j < len(lookAhead) && search[i+j] == lookAhead[j] {
			j++
		} else {
			if j > length {
				length = j
				pos = len(search) - i
			}
			j = 0
			i--
		}
	}
	return
}

func decode(tuples []*lzTuple) (decoded string) {
	for _, t := range tuples {
		if t.pos > 0 {
			correctIndex := len(decoded) - t.pos
			decoded += decoded[correctIndex : correctIndex+t.len]
		}
		if t.c != byte(0) {
			decoded += string(t.c)
		}
	}
	return
}
