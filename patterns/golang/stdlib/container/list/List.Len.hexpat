arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_Len { name: "list.List.Len"; library: "go-stdlib"; instr match { CALL container_list___List__Len } gen { $l.Len(); } bind { l as "l"; } }
