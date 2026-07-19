arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_Reader_ReadAll { name: "csv.Reader.ReadAll"; library: "go-stdlib"; instr match { CALL encoding_csv___Reader__ReadAll } gen { $r.ReadAll(); } bind { r as "r"; } }
