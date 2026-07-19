arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_MoveToBack { name: "list.List.MoveToBack"; library: "go-stdlib"; instr match { CALL container_list___List__MoveToBack } gen { $l.MoveToBack($e); } bind { l as "l"; e as "e"; } }
