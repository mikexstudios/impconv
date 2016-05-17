package main

import "fmt"
import "os"
import "path/filepath"
import "strings"

// import "io/ioutil"
// import "regexp"
// import "strconv"

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "[data.txt]")
		fmt.Println("Output: data.dta file in the same directory.")
		os.Exit(1)
	}
	inFilename := os.Args[1]
	inFilenameNoExt := strings.TrimSuffix(inFilename, filepath.Ext(inFilename))
	outFilename := inFilenameNoExt + ".dta"

	params, data := parseCHIFile(inFilename)
	// fmt.Println(data)

	// Now output Gamry's DTA format
	f, err := os.Create(outFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close() //at end of main()

	fmt.Fprintf(f, "EXPLAIN\r\n")
	fmt.Fprintf(f, "TAG	EISPOT\r\n")
	fmt.Fprintf(f, "TITLE	LABEL	Potentiostatic EIS	Test &Identifier\r\n")
	fmt.Fprintf(f, "\r\n")

	// We can't get exactly three digit precision exponent unless:
	// http://stackoverflow.com/questions/8773133/c-how-to-get-one-digit-exponent-with-printf
	fmt.Fprintf(f, "VDC	POTEN	%11.5E	F	DC &Voltage (V)\r\n", params["InitE"])
	fmt.Fprintf(f, "FREQINIT	QUANT	%11.5E	Initial Fre&q. (Hz)\r\n", params["HighFreq"])
	fmt.Fprintf(f, "FREQFINAL	QUANT	%11.5E	Final Fre&q. (Hz)\r\n", params["LowFreq"])
	// PTSPERDEC	QUANT	1.00000E+001	Points/&decade
	fmt.Fprintf(f, "VAC	QUANT	%11.5E	AC &Voltage (mV rms)\r\n", params["Amplitude"])
	// AREA	QUANT	1.00000E+000	&Area (cm^2)
	// CONDIT	TWOPARAM	F	1.50000E+001	0.00000E+000	Conditionin&g	Time(s)	E(V)
	// DELAY	TWOPARAM	F	1.00000E+002	0.00000E+000	Init. De&lay	Time(s)	Stab.(mV/s)
	// SPEED	SELECTOR	1	&Optimize for:
	// ZGUESS	QUANT	2.00000E+002	E&stimated Z (ohms)
	// EOC	QUANT	0.1358522	Open Circuit (V)

	fmt.Fprintf(f, "ZCURVE	TABLE\r\n")
	fmt.Fprintf(f, "	Pt	Time	Freq	Zreal	Zimag	Zsig	Zmod	Zphz	Idc	Vdc	IERange\r\n")
	fmt.Fprintf(f, "	#	s	Hz	ohm	ohm	V	ohm	°	A	V	#\r\n")
	for i, d := range data {
		fmt.Fprintf(f, "\t%d\t%d\t%f\t%f\t%f\t1\t%f\t%f\t0.000000E-000\t0.000000\t10\r\n",
			i, i, //we don't have time information, so just use #
			d["Freq"], d["Zp"], d["Zpp"], d["Z"], d["Phase"])
	}

}
