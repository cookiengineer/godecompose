arch x86_64;
platform linux, darwin, windows, freebsd;

pattern html_template_Must { name: "html/template.Must"; library: "go-stdlib"; instr match { CALL html_template_Must } gen { template.Must($t, $err); } bind { t as "t"; err as "err"; } }
