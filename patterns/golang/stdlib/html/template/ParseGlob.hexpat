arch x86_64;
platform linux, darwin, windows, freebsd;

pattern html_template_ParseGlob { name: "html/template.ParseGlob"; library: "go-stdlib"; instr match { CALL html_template_ParseGlob } gen { template.ParseGlob($pattern); } bind { pattern as "pattern"; } }
