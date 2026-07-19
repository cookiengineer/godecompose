arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_Ints { name: "sort.Ints"; library: "go-stdlib"; instr match { CALL sort_Ints } gen { sort.Ints($a); } bind { a as "a"; } }
