arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gob_Register { name: "gob.Register"; library: "go-stdlib"; instr match { CALL encoding_gob_Register } gen { gob.Register($value); } bind { value as "value"; } }
