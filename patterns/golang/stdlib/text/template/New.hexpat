arch x86_64;
platform linux, darwin, windows, freebsd;

pattern text_template_New { name: "text/template.New"; library: "go-stdlib"; instr match { CALL text_template_New } gen { template.New($name); } bind { name as "name"; } }
