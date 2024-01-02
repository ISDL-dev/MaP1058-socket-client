package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	input := flag.String("i", "", "input file")
	output := flag.String("o", "output/transposed.csv", "output file")
	trimIndex := flag.Bool("trim-index", false, "trim index")
	trimPoint := flag.Bool("trim-point", false, "trim point")
	flag.Parse()

	if input == nil || *input == "" {
		flag.Usage()
		handleErr(fmt.Errorf("input file is required"))
	}

	file, err := os.Open(*input)
	if err != nil {
		handleErr(fmt.Errorf("failed to open input file: %w", err))
	}
	defer file.Close()
	r := csv.NewReader(file)
	all, err := r.ReadAll()
	if err != nil {
		handleErr(fmt.Errorf("failed to read csv records all: %w", err))
	}

	tf, err := os.Create(*output)
	if err != nil {
		handleErr(fmt.Errorf("failed to create output file: %w", err))
	}
	defer tf.Close()

	if *trimIndex {
		all = all[1:]
	}
	transposed := transpose(all)
	if *trimPoint {
		r1 := transposed[0]
		if r1[0] == "1" {
			transposed = transposed[1:]
		}
	}

	w := csv.NewWriter(tf)
	if err := w.WriteAll(transposed); err != nil {
		os.Remove(*output)
		handleErr(fmt.Errorf("failed to write transposed csv records all: %w", err))
	}
}

func transpose(matrix [][]string) [][]string {
	rowLen := len(matrix)
	colLen := len(matrix[0])

	result := make([][]string, colLen)
	for i := range result {
		result[i] = make([]string, rowLen)
	}

	for i, row := range matrix {
		for j, val := range row {
			result[j][i] = val
		}
	}

	return result
}

func handleErr(err error) {
	if err != nil {
		return
	}
	log.Println(err)
	os.Exit(1)
}
