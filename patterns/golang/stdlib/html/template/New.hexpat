arch x86_64;
platform linux, darwin, windows, freebsd;

pattern html_template_New { name: "html/template.New"; library: "go-stdlib"; instr match { CALL html_template_New } gen { template.New($name); } bind { name as "name"; } }
