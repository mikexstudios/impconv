impconv
=======

impconv is an electrochemical impedance file format converter. Currently, it
converts impedance data from CH Instrument's .txt format to Gamry's
.dta format or Zview's .z format.

## Usage

**[Download binaries on the releases page!][1]**

```
go get github.com/mikexstudios/impconv
go build chi2gamry.go chi.go
go build chi2zview.go chi.go
./chi2gamry data.txt [data2.txt ...]
./chi2zview data.txt [data2.txt ...]
```
`data.dta` and `data.z` will be created in the same directory. For convenience,
you may also drag and drop `data.txt` (or more files) on to the executable.

Cross-compile with (see [syslist][2]):
```
GOOS=windows GOARCH=386 go build chi2gamry.go chi.go
GOOS=windows GOARCH=386 go build chi2zview.go chi.go
./chi2gamry.exe data.txt [data2.txt ...]
./chi2zview.exe data.txt [data2.txt ...]
```
(or drag and drop `data.txt` or more files)

[1]: https://github.com/mikexstudios/impconv/releases
[2]: https://github.com/golang/go/blob/master/src/go/build/syslist.go 

## Sample formats

### CH Instrument

Using the "Convert to txt" function of CH Instrument's software gives the
following for an "AC Impedance" measurement:

```
Init E (V) = 0.2
High Frequency (Hz) = 1e+5
Low Frequency (Hz) = 1
Imp SF
Amplitude (V) = 0.005
Quiet Time (sec) = 0

Freq/Hz, Z'/ohm, Z"/ohm, Z/ohm, Phase/deg

8.252e+4, 1.079e+2, -2.715e+0, 1.080e+2, -1.4
6.812e+4, 1.077e+2, -3.172e+0, 1.078e+2, -1.7
5.615e+4, 1.076e+2, -3.552e+0, 1.077e+2, -1.9
...
```

### Gamry

Gamry's "Potentiostatic EIS" automatically provides an ASCII file:

```
EXPLAIN
TAG	EISPOT
TITLE	LABEL	Potentiostatic EIS	Test &Identifier
DATE	LABEL	5/11/2016	Date
TIME	LABEL	17:47:03	Time
	
PSTAT	PSTAT	In the hood	Potentiostat
VDC	POTEN	1.30000E+000	F	DC &Voltage (V)
FREQINIT	QUANT	1.00000E+005	Initial Fre&q. (Hz)
FREQFINAL	QUANT	1.00000E+000	Final Fre&q. (Hz)
PTSPERDEC	QUANT	1.00000E+001	Points/&decade
VAC	QUANT	1.00000E+001	AC &Voltage (mV rms)
AREA	QUANT	1.00000E+000	&Area (cm^2)
CONDIT	TWOPARAM	F	1.50000E+001	0.00000E+000	Conditionin&g	Time(s)	E(V)
DELAY	TWOPARAM	F	1.00000E+002	0.00000E+000	Init. De&lay	Time(s)	Stab.(mV/s)
SPEED	SELECTOR	1	&Optimize for:
ZGUESS	QUANT	2.00000E+002	E&stimated Z (ohms)
EOC	QUANT	0.1358522	Open Circuit (V)
ZCURVE	TABLE
	Pt	Time	Freq	Zreal	Zimag	Zsig	Zmod	Zphz	Idc	Vdc	IERange
	#	s	Hz	ohm	ohm	V	ohm	°	A	V	#
	0	1	100019.5	224.6075	-3.767681	1	224.6391	-0.961018	2.402966E-006	1.299216	10
	1	3	79511.72	224.712	-4.283262	1	224.7528	-1.091989	1.847788E-006	1.299205	10
	2	4	63105.47	225.1894	-4.847088	1	225.2416	-1.233072	1.641699E-006	1.299213	10
	3	6	50214.84	225.5566	-5.513721	1	225.624	-1.400314	1.570993E-006	1.299209	10
	4	7	39902.34	226.2954	-6.136346	1	226.3786	-1.553282	1.454366E-006	1.299218	10
	5	8	31699.22	226.885	-7.053867	1	226.9946	-1.780755	1.42282E-006	1.299213	10
...
```

This is already simplified and many lines (from `PSTAT` to `EOC`) can be omitted
and Echem Analyst software will still display the Bode and Nyquit plots.


### Zview

The older Zview data format is very simple:

```
"ZView Calculated Data File: Version 1.1"
"  Freq (Hz)    Ampl     Bias   Time(Sec)   Z'(a)    Z''(b)    GD   Err   Range"
 6.500000E+0004,  0.000000E+0000,  0.000000E+0000,  1.000000E+0000,  1.773600E+0003, -6.670100E+0000,  0.000000E+0000, 0, 0
 5.793100E+0004,  0.000000E+0000,  0.000000E+0000,  2.000000E+0000,  1.893600E+0003, -1.203200E+0001,  0.000000E+0000, 0, 0
 5.163100E+0004,  0.000000E+0000,  0.000000E+0000,  3.000000E+0000,  1.807600E+0003, -3.376700E+0001,  0.000000E+0000, 0, 0
 4.601600E+0004,  0.000000E+0000,  0.000000E+0000,  4.000000E+0000,  1.806000E+0003, -3.669100E+0001,  0.000000E+0000, 0, 0
...
```
