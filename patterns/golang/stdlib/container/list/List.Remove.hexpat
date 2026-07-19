arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_Remove { name: "list.List.Remove"; library: "go-stdlib"; instr match { CALL container_list___List__Remove } gen { $l.Remove($e); } bind { l as "l"; e as "e"; } }
