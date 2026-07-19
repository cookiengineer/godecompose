arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Float { name: "reflect.Value.Float"; library: "go-stdlib"; instr match { CALL reflect___Value__Float } gen { $v.Float(); } bind { v as "v"; } }
