arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_IsValid { name: "reflect.Value.IsValid"; library: "go-stdlib"; instr match { CALL reflect___Value__IsValid } gen { $v.IsValid(); } bind { v as "v"; } }
