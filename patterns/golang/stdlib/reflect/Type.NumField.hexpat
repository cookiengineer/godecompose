arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_NumField { name: "reflect.Type.NumField"; library: "go-stdlib"; instr match { CALL reflect___Type__NumField } gen { $t.NumField(); } bind { t as "t"; } }
