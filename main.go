package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// P10 Permutation table
var p10 = []int{3, 5, 2, 7, 4, 10, 1, 9, 8, 6}

// P8 Permutation table
var p8 = []int{6, 3, 7, 4, 8, 5, 10, 9}

// Initial Permutation table
var ip = []int{2, 6, 3, 1, 4, 8, 5, 7}

// Expansion Permutation table
var ep = []int{4, 1, 2, 3, 2, 3, 4, 1}

// Final Permutation table
var ipInverse = []int{4, 1, 3, 5, 7, 2, 8, 6}

// S-Box tables
var sBox0 = [][]int{
	{1, 0, 3, 2},
	{3, 2, 1, 0},
	{0, 2, 1, 3},
	{3, 1, 3, 2},
}

var sBox1 = [][]int{
	{0, 1, 2, 3},
	{2, 0, 1, 3},
	{3, 0, 1, 0},
	{2, 1, 0, 3},
}

// Function to read permutation tables from file
func readPermutations(filename string) map[string][]int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	permutations := make(map[string][]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		name := strings.TrimSpace(parts[0])
		values := strings.Fields(parts[1])
		var permutation []int
		for _, v := range values {
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("Error converting to integer:", err)
				os.Exit(1)
			}
			permutation = append(permutation, num)
		}
		permutations[name] = permutation
	}
	return permutations
}

// Function to read key and plaintext from file
func readInput(filename string) (string, string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var key, plaintext string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		switch strings.TrimSpace(parts[0]) {
		case "Key":
			key = strings.TrimSpace(parts[1])
		case "Plaintext":
			plaintext = strings.TrimSpace(parts[1])
		}
	}
	return key, plaintext
}

// Function to perform permutation
func permute(input []int, table []int) []int {
	output := make([]int, len(table))
	for i, v := range table {
		output[i] = input[v-1]
	}
	return output
}

// Function to perform circular left shift
func leftShift(input []int, shift int) []int {
	n := len(input)
	output := make([]int, n)
	for i := 0; i < n; i++ {
		output[i] = input[(i+shift)%n]
	}
	return output
}

// Function to perform XOR operation
func xor(a, b []int) []int {
	n := len(a)
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

// Function to perform S-Box lookup
func sBoxLookup(input []int, sBox [][]int) []int {
	row := input[0]*2 + input[3]
	col := input[1]*2 + input[2]
	return []int{(sBox[row][col] >> 1) & 1, sBox[row][col] & 1}
}

// Function to generate keys
func generateKeys(key []int) (k1, k2 []int) {
	// Apply P10 permutation
	p10Key := permute(key, p10)
	// Split into Left and Right halves
	left := p10Key[:5]
	right := p10Key[5:]
	// Left shift (LS-1) on both halves
	left = leftShift(left, 1)
	right = leftShift(right, 1)
	// Combining and applying P8 permutation to generate K1
	k1 = permute(append(left, right...), p8)
	// Left shift (LS-2) on both halves
	left = leftShift(left, 2)
	right = leftShift(right, 2)
	// Combining and applying P8 permutation to generate K2
	k2 = permute(append(left, right...), p8)
	return k1, k2
}

// Function to encrypt plaintext
func encrypt(plaintext, key []int) []int {
	// Initial Permutation (IP)
	permutedPlaintext := permute(plaintext, ip)
	// Splitting into Left and Right halves
	left := permutedPlaintext[:4]
	right := permutedPlaintext[4:]
	// Expansion Permutation (EP) on Right half
	expanded := permute(right, ep)
	// XOR with Key (K1)
	xored := xor(expanded, key)
	// S-Box substitution
	sBoxOutput := append(sBoxLookup(xored[:4], sBox0), sBoxLookup(xored[4:], sBox1)...)
	// Permutation (P4)
	permuted := permute(sBoxOutput, []int{2, 4, 3, 1})
	// XOR with Left half
	newRight := xor(permuted, left)
	// Final permutation (IP-1)
	ciphertext := permute(append(newRight, right...), ipInverse)
	return ciphertext
}

func main() {
	// Read permutation tables from file
	permutations := readPermutations("permutations.txt")

	// Read key and plaintext from file
	key, plaintext := readInput("input.txt")

	// Parsing key and plaintext
	keyBits := make([]int, len(key))
	for i, char := range key {
		keyBits[i] = int(char - '0')
	}

	plaintextBits := make([]int, len(plaintext))
	for i, char := range plaintext {
		plaintextBits[i] = int(char - '0')
	}

	// Generate keys
	k1, k2 := generateKeys(keyBits)
	fmt.Println("K1:", k1)
	fmt.Println("K2:", k2)

	// Encrypt plaintext
	ciphertext := encrypt(plaintextBits, k1)
	ciphertext = encrypt(ciphertext, k2)

	// Print ciphertext
	fmt.Println("Ciphertext:", ciphertext)
}
