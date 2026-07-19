arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_SearchStrings { name: "sort.SearchStrings"; library: "go-stdlib"; instr match { CALL sort_SearchStrings } gen { sort.SearchStrings($a, $x); } bind { a as "a"; x as "x"; } }
