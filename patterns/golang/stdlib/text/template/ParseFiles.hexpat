arch x86_64;
platform linux, darwin, windows, freebsd;

pattern text_template_ParseFiles { name: "text/template.ParseFiles"; library: "go-stdlib"; instr match { CALL text_template_ParseFiles } gen { template.ParseFiles($filenames...); } }
