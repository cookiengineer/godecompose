arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_ValueOf { name: "reflect.ValueOf"; library: "go-stdlib"; instr match { CALL reflect_ValueOf } gen { reflect.ValueOf($i); } bind { i as "i"; } }
