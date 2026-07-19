arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_MoveToFront { name: "list.List.MoveToFront"; library: "go-stdlib"; instr match { CALL container_list___List__MoveToFront } gen { $l.MoveToFront($e); } bind { l as "l"; e as "e"; } }
