arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Interface { name: "reflect.Value.Interface"; library: "go-stdlib"; instr match { CALL reflect___Value__Interface } gen { $v.Interface(); } bind { v as "v"; } }
