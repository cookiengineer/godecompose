arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_Writer_Write { name: "csv.Writer.Write"; library: "go-stdlib"; instr match { CALL encoding_csv___Writer__Write } gen { $w.Write($record); } bind { w as "w"; record as "record"; } }
