arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_Float64s { name: "sort.Float64s"; library: "go-stdlib"; instr match { CALL sort_Float64s } gen { sort.Float64s($a); } bind { a as "a"; } }
