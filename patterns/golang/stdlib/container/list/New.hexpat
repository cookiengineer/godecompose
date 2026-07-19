arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_New { name: "container/list.New"; library: "go-stdlib"; instr match { CALL container_list_New } gen { list.New(); } }
