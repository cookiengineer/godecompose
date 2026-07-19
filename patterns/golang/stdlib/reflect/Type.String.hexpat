arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_String { name: "reflect.Type.String"; library: "go-stdlib"; instr match { CALL reflect___Type__String } gen { $t.String(); } bind { t as "t"; } }
