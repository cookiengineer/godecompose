arch x86_64;
platform linux, darwin, windows, freebsd;

pattern text_template_ParseGlob { name: "text/template.ParseGlob"; library: "go-stdlib"; instr match { CALL text_template_ParseGlob } gen { template.ParseGlob($pattern); } bind { pattern as "pattern"; } }
