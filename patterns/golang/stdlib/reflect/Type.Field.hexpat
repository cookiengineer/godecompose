arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_Field { name: "reflect.Type.Field"; library: "go-stdlib"; instr match { CALL reflect___Type__Field } gen { $t.Field($i); } bind { t as "t"; i as "i"; } }
