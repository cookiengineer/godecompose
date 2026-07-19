arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Value_Bool { name: "reflect.Value.Bool"; library: "go-stdlib"; instr match { CALL reflect___Value__Bool } gen { $v.Bool(); } bind { v as "v"; } }
