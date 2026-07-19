arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Len { name: "reflect.Value.Len"; library: "go-stdlib"; instr match { CALL reflect___Value__Len } gen { $v.Len(); } bind { v as "v"; } }
