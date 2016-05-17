package main

import "fmt"
import "os"
import "path/filepath"
import "strings"

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "[data.txt] [data2.txt] ...")
		fmt.Println("Output: .z file(s) in the same directory.")
		os.Exit(1)
	}
	inFilenames := os.Args[1:]

	for _, inFilename := range inFilenames {
		inFilenameNoExt := strings.TrimSuffix(inFilename, filepath.Ext(inFilename))
		outFilename := inFilenameNoExt + ".z"

		params, data := parseCHIFile(inFilename)
		writeZviewFile(outFilename, params, data)
	}
}

func writeZviewFile(filename string, params map[string]float64, data []map[string]float64) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "\"ZView Calculated Data File: Version 1.1\"\r\n")
	fmt.Fprintf(f, "\"  Freq (Hz)    Ampl     Bias   Time(Sec)   Z'(a)    Z''(b)    GD   Err   Range\"\r\n")

	// We can't get exactly three digit precision exponent unless:
	// http://stackoverflow.com/questions/8773133/c-how-to-get-one-digit-exponent-with-printf
	for i, d := range data {
		// Sample:
		// 6.500000E+0004,  0.000000E+0000,  0.000000E+0000,  1.000000E+0000,  1.773600E+0003, -6.670100E+0000,  0.000000E+0000, 0, 0
		fmt.Fprintf(f, " %12.6E,  0.000000E+0000,  0.000000E+0000, %12.6E,  %12.6E,  %12.6E,  0.000000E+0000, 0, 0\r\n",
			d["Freq"], float64(i), d["Zp"], d["Zpp"]) // we don't have time info, so just use i
	}
}
