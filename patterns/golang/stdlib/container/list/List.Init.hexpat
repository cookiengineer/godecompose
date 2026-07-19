arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_Init { name: "list.List.Init"; library: "go-stdlib"; instr match { CALL container_list___List__Init } gen { $l.Init(); } bind { l as "l"; } }
