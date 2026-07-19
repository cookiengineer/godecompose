arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_Front { name: "list.List.Front"; library: "go-stdlib"; instr match { CALL container_list___List__Front } gen { $l.Front(); } bind { l as "l"; } }
