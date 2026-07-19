arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_PushBack { name: "list.List.PushBack"; library: "go-stdlib"; instr match { CALL container_list___List__PushBack } gen { $l.PushBack($v); } bind { l as "l"; v as "v"; } }
