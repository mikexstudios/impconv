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
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "[chi-data.txt]")
		os.Exit(1)
	}
	in := os.Args[1]

	// Since files are small, read whole thing into string
	whole, err := ioutil.ReadFile(in)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(whole), "\n")
	fmt.Println(len(lines))
	// fmt.Println(lines[len(lines)-2])

	// We want to match the following:
	var InitE, HighFreq, LowFreq, Amplitude float64
	re := make(map[string]*regexp.Regexp)
	// Init E (V) = 0.2
	re["InitE"], _ = regexp.Compile(`^Init E \(V\) =\s*(.+)\s*`)
	// High Frequency (Hz) = 1e+5
	re["HighFreq"], _ = regexp.Compile(`^High Frequency \(Hz\) =\s*(.+)\s*`)
	// Low Frequency (Hz) = 1
	re["LowFreq"], _ = regexp.Compile(`^Low Frequency \(Hz\) =\s*(.+)\s*`)
	// Imp SF -> ignore
	// Amplitude (V) = 0.005
	re["Amplitude"], _ = regexp.Compile(`^Amplitude \(V\) =\s*(.+)\s*`)
	// Quiet Time (sec) = 0 -> ignore
	// Freq/Hz, Z'/ohm, Z"/ohm, Z/ohm, Phase/deg
	re["Header"], _ = regexp.Compile(`^Freq/Hz,.+`)

	// fmt.Println(re["InitE"].FindStringSubmatch("Init E (V) = 0.2"))

	var inData bool = false
	for _, line := range lines {
		line = strings.TrimSpace(line) //Remove \n, \r, etc.

		if re["Header"].MatchString(line) {
			inData = true
		}

		if !inData {
			// fmt.Printf("h")
			// Check for various key lines
			if sm := re["InitE"].FindStringSubmatch(line); sm != nil {
				InitE, err = strconv.ParseFloat(sm[1], 64)
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

	}

	fmt.Println(InitE, HighFreq, LowFreq, Amplitude)
}
