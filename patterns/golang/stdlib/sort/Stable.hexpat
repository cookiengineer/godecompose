arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_Stable { name: "sort.Stable"; library: "go-stdlib"; instr match { CALL sort_Stable } gen { sort.Stable($data); } bind { data as "data"; } }
