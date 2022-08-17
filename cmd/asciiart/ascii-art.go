package Asciiart

import (
	"fmt"
	"os"
	"strings"
)

func Asciiart(s string, banner string) bool {
	// if isAllNewLine(s) {
	// 	return false
	// }
	input := strings.Split(s, "\r\n")
	// ignore the \n
	// res := []string{}
	// for _, v := range input {
	// 	tempStr := strings.Split(v, "\\n")
	// 	res = append(res, tempStr...)
	// }
	// input = res
	if NonAsciiCheck(input) {
		fmt.Println("ERROR: Exeptional characters!")
		return false
	}

	sint := WriteIndexes(input)

	template, err := os.ReadFile("cmd/asciiart/" + banner + ".txt")
	if err != nil {
		fmt.Println(err)
		return false
	}

	dict := strings.Split(string(template), "\n")
	indexes := CountIndexes(dict)

	// if len(template) != 6623 {
	// 	fmt.Println("ERROR: Corrupted file.")
	// 	return
	// }

	CreateBannerAndWriteToFile(input, dict, sint, indexes)
	return true
}

func CreateBannerAndWriteToFile(input []string, dict []string, sint [][]int, indexes []int) {
	f, err := os.Create("cmd/output.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for i := range input {
		if len(sint[i]) == 0 {
			// write "\n" to the file
			if _, err := f.WriteString("\n"); err != nil {
				panic(err)
			}

			// fmt.Println()
			continue
		}

		banner := WordsToAscii(dict, sint[i], indexes)

		// write banner to the file
		for _, v := range banner {
			if _, err := f.WriteString(v); err != nil {
				panic(err)
			}
		}

		// PrintBanner(banner)
	}
}

func NonAsciiCheck(input []string) bool {
	for _, v := range input {
		for _, b := range v {
			if !((32 <= b && b <= 126) || b == 13 || b == 10) {
				return true
			}
		}
	}
	return false
}

func LineisNewline(s string) bool {
	for _, v := range s {
		if v != '\n' {
			return false
		}
	}
	return true
}

func isAllNewLine(input string) bool {
	isAllNewLine := true
	for i := range input {
		if !(input[i] == '\\' && input[i+1] == 'n') {
			isAllNewLine = false
		}
		i++
	}
	if isAllNewLine {
		for i := 0; i < len(input)-1; i += 2 {
			fmt.Println()
		}
	}
	return isAllNewLine
}

func CountIndexes(a []string) []int {
	var indexes []int
	j := 0
	for i := range a {
		// fmt.Println(i, a[i], LineisNewline(a[i]))
		if LineisNewline(a[i]) {
			if i != len(a)-1 {
				indexes = append(indexes, i)
				j++
			}
		}
	}
	return indexes
}

func WriteIndexes(input []string) [][]int {
	sint := make([][]int, len(input))
	for i := 0; i < len(input); i++ {
		sint[i] = make([]int, 0)
	}
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[i]); j++ {
			sint[i] = append(sint[i], int(input[i][j]-' '))
		}
	}

	return sint
}

func PrintBanner(banner [8]string) {
	for _, v := range banner {
		fmt.Print(v)
	}
}

func WordsToAscii(dict []string, sint []int, indexes []int) [8]string {
	var output [8]string
	for i := 0; i < len(output); i++ {
		output[i] = ""
		for _, v := range sint {
			output[i] += dict[i+indexes[v]+1]
		}
		output[i] += "\n"
	}
	return output
}
