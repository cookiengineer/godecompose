arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Index { name: "reflect.Value.Index"; library: "go-stdlib"; instr match { CALL reflect___Value__Index } gen { $v.Index($i); } bind { v as "v"; i as "i"; } }
