arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_Back { name: "list.List.Back"; library: "go-stdlib"; instr match { CALL container_list___List__Back } gen { $l.Back(); } bind { l as "l"; } }
