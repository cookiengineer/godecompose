arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_Type_Elem { name: "reflect.Type.Elem"; library: "go-stdlib"; instr match { CALL reflect___Type__Elem } gen { $t.Elem(); } bind { t as "t"; } }
