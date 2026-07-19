arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_Reader_Read { name: "csv.Reader.Read"; library: "go-stdlib"; instr match { CALL encoding_csv___Reader__Read } gen { $r.Read(); } bind { r as "r"; } }
