arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_InsertBefore { name: "list.List.InsertBefore"; library: "go-stdlib"; instr match { CALL container_list___List__InsertBefore } gen { $l.InsertBefore($v, $mark); } bind { l as "l"; v as "v"; mark as "mark"; } }
