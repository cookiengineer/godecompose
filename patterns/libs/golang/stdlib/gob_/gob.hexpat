arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gob_NewEncoder { name: "encoding/gob.NewEncoder"; library: "go-stdlib"; instr match { CALL encoding_gob_NewEncoder } gen { gob.NewEncoder($w); } bind { w as "w"; } }
pattern gob_NewDecoder { name: "encoding/gob.NewDecoder"; library: "go-stdlib"; instr match { CALL encoding_gob_NewDecoder } gen { gob.NewDecoder($r); } bind { r as "r"; } }
pattern gob_Register { name: "gob.Register"; library: "go-stdlib"; instr match { CALL encoding_gob_Register } gen { gob.Register($value); } bind { value as "value"; } }
pattern gob_RegisterName { name: "gob.RegisterName"; library: "go-stdlib"; instr match { CALL encoding_gob_RegisterName } gen { gob.RegisterName($name, $value); } bind { name as "name"; value as "value"; } }
