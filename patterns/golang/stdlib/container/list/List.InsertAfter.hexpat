arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_InsertAfter { name: "list.List.InsertAfter"; library: "go-stdlib"; instr match { CALL container_list___List__InsertAfter } gen { $l.InsertAfter($v, $mark); } bind { l as "l"; v as "v"; mark as "mark"; } }
