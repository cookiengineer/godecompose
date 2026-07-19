arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_NewWriter { name: "encoding/csv.NewWriter"; library: "go-stdlib"; instr match { CALL encoding_csv_NewWriter } gen { csv.NewWriter($w); } bind { w as "w"; } }
