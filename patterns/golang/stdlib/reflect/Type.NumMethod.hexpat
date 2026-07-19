arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_NumMethod { name: "reflect.Type.NumMethod"; library: "go-stdlib"; instr match { CALL reflect___Type__NumMethod } gen { $t.NumMethod(); } bind { t as "t"; } }
