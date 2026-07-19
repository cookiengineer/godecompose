arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_PushFront { name: "list.List.PushFront"; library: "go-stdlib"; instr match { CALL container_list___List__PushFront } gen { $l.PushFront($v); } bind { l as "l"; v as "v"; } }
