package main

import "fmt"
import "os"
import "path/filepath"
import "io/ioutil"
import "strings"
import "regexp"
import "strconv"

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "[data.txt]")
		fmt.Println("Output: data.z file in the same directory.")
		os.Exit(1)
	}
	inFilename := os.Args[1]
	inFilenameNoExt := strings.TrimSuffix(inFilename, filepath.Ext(inFilename))
	outFilename := inFilenameNoExt + ".z"

	// Since files are small, read whole thing into string
	whole, err := ioutil.ReadFile(inFilename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(whole), "\n")
	// fmt.Println(len(lines))
	// fmt.Println(lines[len(lines)-2])

	// We want to match the following parameters:
	// Init E (V) = 0.2
	// High Frequency (Hz) = 1e+5
	// Low Frequency (Hz) = 1
	// Imp SF -> ignore
	// Amplitude (V) = 0.005
	// Quiet Time (sec) = 0 -> ignore
	// Freq/Hz, Z'/ohm, Z"/ohm, Z/ohm, Phase/deg
	var InitE, HighFreq, LowFreq, Amplitude float64
	re := make(map[string]*regexp.Regexp)
	re["InitE"], _ = regexp.Compile(`^Init E \(V\) =\s*(.+)\s*`)
	re["HighFreq"], _ = regexp.Compile(`^High Frequency \(Hz\) =\s*(.+)\s*`)
	re["LowFreq"], _ = regexp.Compile(`^Low Frequency \(Hz\) =\s*(.+)\s*`)
	re["Amplitude"], _ = regexp.Compile(`^Amplitude \(V\) =\s*(.+)\s*`)
	re["Header"], _ = regexp.Compile(`^Freq/Hz,.+`)

	// We store the data as a slice of maps (with keys as columns):
	data := make([]map[string]float64, 0)

	var inData bool = false
	for _, line := range lines {
		line = strings.TrimSpace(line) //Remove \n, \r, etc.
		if line == "" {
			continue
		}

		if re["Header"].MatchString(line) {
			inData = true
			continue
		}

		if !inData {
			// fmt.Printf("h")
			// Check for various key lines
			if sm := re["InitE"].FindStringSubmatch(line); sm != nil {
				InitE, _ = strconv.ParseFloat(sm[1], 64)
			}
			if sm := re["HighFreq"].FindStringSubmatch(line); sm != nil {
				// ParseFloat can handle "scientific" formats, e.g., 1e-3
				HighFreq, _ = strconv.ParseFloat(sm[1], 64)
			}
			if sm := re["LowFreq"].FindStringSubmatch(line); sm != nil {
				LowFreq, _ = strconv.ParseFloat(sm[1], 64)
			}
			if sm := re["Amplitude"].FindStringSubmatch(line); sm != nil {
				Amplitude, _ = strconv.ParseFloat(sm[1], 64)
			}

			continue
		}
		// fmt.Printf("d")

		// Now parse data for each line
		d := make(map[string]float64, 5)
		s := strings.Split(line, ",")
		d["Freq"], _ = strconv.ParseFloat(strings.TrimSpace(s[0]), 64)
		d["Zp"], _ = strconv.ParseFloat(strings.TrimSpace(s[1]), 64)
		d["Zpp"], _ = strconv.ParseFloat(strings.TrimSpace(s[2]), 64)
		d["Z"], _ = strconv.ParseFloat(strings.TrimSpace(s[3]), 64)
		d["Phase"], _ = strconv.ParseFloat(strings.TrimSpace(s[4]), 64)
		//fmt.Println(d)
		data = append(data, d)
	}

	// fmt.Println(InitE, HighFreq, LowFreq, Amplitude)
	// fmt.Println(data)

	// Now output Zview's older format
	f, err := os.Create(outFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close() //at end of main()

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

	// Need to use these vars
	fmt.Println(InitE, HighFreq, LowFreq, Amplitude)
}
