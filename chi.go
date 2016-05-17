package main

// import "fmt"
// import "os"
// import "path/filepath"
import "io/ioutil"
import "strings"
import "regexp"
import "strconv"

func parseCHIFile(filename string) (params map[string]float64, data []map[string]float64) {
	// Since files are small, read whole thing into string
	whole, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(whole), "\n")

	// We want to match the following parameters:
	// Init E (V) = 0.2
	// High Frequency (Hz) = 1e+5
	// Low Frequency (Hz) = 1
	// Imp SF -> ignore
	// Amplitude (V) = 0.005
	// Quiet Time (sec) = 0 -> ignore
	// Freq/Hz, Z'/ohm, Z"/ohm, Z/ohm, Phase/deg
	re := make(map[string]*regexp.Regexp)
	re["InitE"], _ = regexp.Compile(`^Init E \(V\) =\s*(.+)\s*`)
	re["HighFreq"], _ = regexp.Compile(`^High Frequency \(Hz\) =\s*(.+)\s*`)
	re["LowFreq"], _ = regexp.Compile(`^Low Frequency \(Hz\) =\s*(.+)\s*`)
	re["Amplitude"], _ = regexp.Compile(`^Amplitude \(V\) =\s*(.+)\s*`)
	re["Header"], _ = regexp.Compile(`^Freq/Hz,.+`)

	// We store the data as a slice of maps (with keys as columns):
	params = make(map[string]float64)
	data = make([]map[string]float64, 0)

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
				params["InitE"], _ = strconv.ParseFloat(sm[1], 64)
			}
			if sm := re["HighFreq"].FindStringSubmatch(line); sm != nil {
				// ParseFloat can handle "scientific" formats, e.g., 1e-3
				params["HighFreq"], _ = strconv.ParseFloat(sm[1], 64)
			}
			if sm := re["LowFreq"].FindStringSubmatch(line); sm != nil {
				params["LowFreq"], _ = strconv.ParseFloat(sm[1], 64)
			}
			if sm := re["Amplitude"].FindStringSubmatch(line); sm != nil {
				params["Amplitude"], _ = strconv.ParseFloat(sm[1], 64)
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

	return
}
