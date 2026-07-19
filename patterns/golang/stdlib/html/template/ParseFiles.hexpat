arch x86_64;
platform linux, darwin, windows, freebsd;

pattern html_template_ParseFiles { name: "html/template.ParseFiles"; library: "go-stdlib"; instr match { CALL html_template_ParseFiles } gen { template.ParseFiles($filenames...); } }
