arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_Size { name: "reflect.Type.Size"; library: "go-stdlib"; instr match { CALL reflect___Type__Size } gen { $t.Size(); } bind { t as "t"; } }
