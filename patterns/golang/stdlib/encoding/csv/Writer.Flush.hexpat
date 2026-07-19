arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_Writer_Flush { name: "csv.Writer.Flush"; library: "go-stdlib"; instr match { CALL encoding_csv___Writer__Flush } gen { $w.Flush(); } bind { w as "w"; } }
