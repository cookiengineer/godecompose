arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gob_NewEncoder { name: "encoding/gob.NewEncoder"; library: "go-stdlib"; instr match { CALL encoding_gob_NewEncoder } gen { gob.NewEncoder($w); } bind { w as "w"; } }
