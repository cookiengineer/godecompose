arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Close { name: "reflect.Value.Close"; library: "go-stdlib"; instr match { CALL reflect___Value__Close } gen { $v.Close(); } bind { v as "v"; } }
