arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_TypeOf { name: "reflect.TypeOf"; library: "go-stdlib"; instr match { CALL reflect_TypeOf } gen { reflect.TypeOf($i); } bind { i as "i"; } }
