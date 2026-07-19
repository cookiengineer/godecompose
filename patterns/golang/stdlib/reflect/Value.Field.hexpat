arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Field { name: "reflect.Value.Field"; library: "go-stdlib"; instr match { CALL reflect___Value__Field } gen { $v.Field($i); } bind { v as "v"; i as "i"; } }
