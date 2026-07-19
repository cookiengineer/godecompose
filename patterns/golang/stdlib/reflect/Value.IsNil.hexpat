arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_IsNil { name: "reflect.Value.IsNil"; library: "go-stdlib"; instr match { CALL reflect___Value__IsNil } gen { $v.IsNil(); } bind { v as "v"; } }
