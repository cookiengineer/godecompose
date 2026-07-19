arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_IsSorted { name: "sort.IsSorted"; library: "go-stdlib"; instr match { CALL sort_IsSorted } gen { sort.IsSorted($data); } bind { data as "data"; } }
