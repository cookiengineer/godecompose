arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_SetString { name: "reflect.Value.SetString"; library: "go-stdlib"; instr match { CALL reflect___Value__SetString } gen { $v.SetString($x); } bind { v as "v"; x as "x"; } }
