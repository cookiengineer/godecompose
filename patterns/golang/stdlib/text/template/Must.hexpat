arch x86_64;
platform linux, darwin, windows, freebsd;

pattern text_template_Must { name: "text/template.Must"; library: "go-stdlib"; instr match { CALL text_template_Must } gen { template.Must($t, $err); } bind { t as "t"; err as "err"; } }
