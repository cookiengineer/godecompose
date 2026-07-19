arch x86_64;
platform linux, darwin, windows, freebsd;

pattern reflect_DeepEqual { name: "reflect.DeepEqual"; library: "go-stdlib"; instr match { CALL reflect_DeepEqual } gen { reflect.DeepEqual($x, $y); } bind { x as "x"; y as "y"; } }
