arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_Search { name: "sort.Search"; library: "go-stdlib"; instr match { CALL sort_Search } gen { sort.Search($n, $f); } bind { n as "n"; f as "f"; } }
