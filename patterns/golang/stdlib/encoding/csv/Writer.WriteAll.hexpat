arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_Writer_WriteAll { name: "csv.Writer.WriteAll"; library: "go-stdlib"; instr match { CALL encoding_csv___Writer__WriteAll } gen { $w.WriteAll($records); } bind { w as "w"; records as "records"; } }
