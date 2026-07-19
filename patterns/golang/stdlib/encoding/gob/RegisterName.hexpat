arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gob_RegisterName { name: "gob.RegisterName"; library: "go-stdlib"; instr match { CALL encoding_gob_RegisterName } gen { gob.RegisterName($name, $value); } bind { name as "name"; value as "value"; } }
