arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_Kind { name: "reflect.Type.Kind"; library: "go-stdlib"; instr match { CALL reflect___Type__Kind } gen { $t.Kind(); } bind { t as "t"; } }
