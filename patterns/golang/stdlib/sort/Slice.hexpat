arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_Slice { name: "sort.Slice"; library: "go-stdlib"; instr match { CALL sort_Slice } gen { sort.Slice($s, $less); } bind { s as "s"; less as "less"; } }
