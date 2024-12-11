package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func addToSlice(s []int, value, amount int) []int {
	if amount < 0 {
		amount = amount * -1
		value = 0
	}

	for range amount {
		s = append(s, value)
	}
	return s
}

func calcCheckSum(disc []int) int {
	total := 0
	for i, id := range disc {
		total += i * id
	}
	return total
}

func part1(data []string) {
	ordered := []int{}

	var extraSpace int
	j := len(data) - 1
	i := 1

	endblock, err := strconv.Atoi(data[j])
	if err != nil {
		log.Fatal(err)
	}
	idI, idJ := 0, int(j/2)

	for i <= j {
		if endblock <= 0 && i > 1 {
			j -= 2
			if j == 0 {
				break
			}
			idJ = int(j / 2)
			endblock, err = strconv.Atoi(data[j])
			if err != nil {
				log.Fatal(err)
			}
		}
		if extraSpace > 0 {
			if endblock >= extraSpace {
				ordered = addToSlice(ordered, idJ, extraSpace)
				endblock -= extraSpace
				extraSpace = 0
			} else {
				extraSpace = extraSpace - endblock
				ordered = addToSlice(ordered, idJ, endblock)
				endblock = 0
			}

		} else {
			takenSpace, err := strconv.Atoi(data[i-1])
			if err != nil {
				log.Fatal(err)
			}
			emptySpace, err := strconv.Atoi(data[i])
			if err != nil {
				log.Fatal(err)
			}
			ordered = addToSlice(ordered, idI, takenSpace)
			if j < i {
				endblock = 0
				break
			}

			if endblock > 0 {
				if endblock <= emptySpace {
					ordered = addToSlice(ordered, idJ, endblock)
					emptySpace -= endblock
					endblock = 0
				} else {
					endblock -= emptySpace
					ordered = addToSlice(ordered, idJ, emptySpace)
					emptySpace = 0
				}

			}
			if emptySpace > 0 {
				if endblock >= emptySpace {
					ordered = addToSlice(ordered, idJ, emptySpace)
					endblock -= emptySpace
				} else {
					extraSpace = emptySpace - endblock
					ordered = addToSlice(ordered, idJ, endblock)
					endblock = 0
				}
			}
			i += 2
			idI++
		}

	}
	if endblock > 0 {
		ordered = addToSlice(ordered, idJ, endblock)
	}

	log.Printf("Final len %v", len(ordered))
	checksum := calcCheckSum(ordered)
	log.Printf("checksum %v", checksum)
}

func part2(data []string) {
	ordered := []int{}

	i := 1

	idI := 0

	for i < len(data) {
		takenSpace, err := strconv.Atoi(data[i-1])
		if err != nil {
			log.Fatal(err)
		}
		emptySpace, err := strconv.Atoi(data[i])
		if err != nil {
			log.Fatal(err)
		}
		ordered = addToSlice(ordered, idI, takenSpace)
		j := len(data) - 1
		for j > i {
			idJ := int(j / 2)
			endblock, err := strconv.Atoi(data[j])
			if err != nil {
				log.Fatal(err)
			}
			if endblock > 0 && endblock <= emptySpace {
				ordered = addToSlice(ordered, idJ, endblock)
				emptySpace -= endblock
				data[j] = fmt.Sprintf("%v", endblock*-1)

			}
			j -= 2
		}
		if emptySpace > 0 {
			ordered = addToSlice(ordered, 0, emptySpace)
		}
		i += 2
		idI++
	}

	log.Printf("Final len %v", len(ordered))
	checksum := calcCheckSum(ordered)
	log.Printf("checksum %v", checksum)
}

func findtotal(data []string) {
	i := 0
	total := 0
	for i < len(data) {
		block, err := strconv.Atoi(data[i])
		if err != nil {
			log.Fatal(err)
		}
		total += block
		i += 2
	}
	log.Printf("Expect find len to be: %v", total)
}

func main() {

	//f, err := os.Open("./testInput.txt")
	f, err := os.Open("./day9Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	disc := strings.Split(string(data), "")
	log.Printf("Len = %v", len(disc))
	findtotal(disc)

	part1(disc)
	part2(disc)

}
