arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_Name { name: "reflect.Type.Name"; library: "go-stdlib"; instr match { CALL reflect___Type__Name } gen { $t.Name(); } bind { t as "t"; } }
