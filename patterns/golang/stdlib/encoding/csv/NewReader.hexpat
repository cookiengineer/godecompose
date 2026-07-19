arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_NewReader { name: "encoding/csv.NewReader"; library: "go-stdlib"; instr match { CALL encoding_csv_NewReader } gen { csv.NewReader($r); } bind { r as "r"; } }
