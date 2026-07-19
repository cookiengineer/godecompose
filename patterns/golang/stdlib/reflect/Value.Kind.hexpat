arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Kind { name: "reflect.Value.Kind"; library: "go-stdlib"; instr match { CALL reflect___Value__Kind } gen { $v.Kind(); } bind { v as "v"; } }
