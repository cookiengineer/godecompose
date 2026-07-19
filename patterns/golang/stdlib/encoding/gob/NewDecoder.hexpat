arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gob_NewDecoder { name: "encoding/gob.NewDecoder"; library: "go-stdlib"; instr match { CALL encoding_gob_NewDecoder } gen { gob.NewDecoder($r); } bind { r as "r"; } }
