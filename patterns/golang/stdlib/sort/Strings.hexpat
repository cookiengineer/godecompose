arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_Strings { name: "sort.Strings"; library: "go-stdlib"; instr match { CALL sort_Strings } gen { sort.Strings($a); } bind { a as "a"; } }
