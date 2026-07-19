arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_String { name: "reflect.Value.String"; library: "go-stdlib"; instr match { CALL reflect___Value__String } gen { $v.String(); } bind { v as "v"; } }
