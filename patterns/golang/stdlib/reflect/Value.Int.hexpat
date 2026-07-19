arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Int { name: "reflect.Value.Int"; library: "go-stdlib"; instr match { CALL reflect___Value__Int } gen { $v.Int(); } bind { v as "v"; } }
