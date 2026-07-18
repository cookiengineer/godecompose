arch x86_64;
platform linux, darwin, windows, freebsd;

pattern csv_NewReader { name: "encoding/csv.NewReader"; library: "go-stdlib"; instr match { CALL encoding_csv_NewReader } gen { csv.NewReader($r); } bind { r as "r"; } }
pattern csv_NewWriter { name: "encoding/csv.NewWriter"; library: "go-stdlib"; instr match { CALL encoding_csv_NewWriter } gen { csv.NewWriter($w); } bind { w as "w"; } }
pattern csv_Reader_Read { name: "csv.Reader.Read"; library: "go-stdlib"; instr match { CALL encoding_csv___Reader__Read } gen { $r.Read(); } bind { r as "r"; } }
pattern csv_Reader_ReadAll { name: "csv.Reader.ReadAll"; library: "go-stdlib"; instr match { CALL encoding_csv___Reader__ReadAll } gen { $r.ReadAll(); } bind { r as "r"; } }
pattern csv_Writer_Write { name: "csv.Writer.Write"; library: "go-stdlib"; instr match { CALL encoding_csv___Writer__Write } gen { $w.Write($record); } bind { w as "w"; record as "record"; } }
pattern csv_Writer_WriteAll { name: "csv.Writer.WriteAll"; library: "go-stdlib"; instr match { CALL encoding_csv___Writer__WriteAll } gen { $w.WriteAll($records); } bind { w as "w"; records as "records"; } }
pattern csv_Writer_Flush { name: "csv.Writer.Flush"; library: "go-stdlib"; instr match { CALL encoding_csv___Writer__Flush } gen { $w.Flush(); } bind { w as "w"; } }
