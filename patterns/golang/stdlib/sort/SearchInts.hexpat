arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_SearchInts { name: "sort.SearchInts"; library: "go-stdlib"; instr match { CALL sort_SearchInts } gen { sort.SearchInts($a, $x); } bind { a as "a"; x as "x"; } }
