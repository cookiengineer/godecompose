arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Elem { name: "reflect.Value.Elem"; library: "go-stdlib"; instr match { CALL reflect___Value__Elem } gen { $v.Elem(); } bind { v as "v"; } }
