arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sort_Reverse { name: "sort.Reverse"; library: "go-stdlib"; instr match { CALL sort_Reverse } gen { sort.Reverse($data); } bind { data as "data"; } }
