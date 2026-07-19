arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_SearchFloat64s { name: "sort.SearchFloat64s"; library: "go-stdlib"; instr match { CALL sort_SearchFloat64s } gen { sort.SearchFloat64s($a, $x); } bind { a as "a"; x as "x"; } }
