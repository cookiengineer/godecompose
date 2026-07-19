arch x86_64;
platform linux, darwin, windows, freebsd;

pattern text_template_Execute { name: "text/template.Template.Execute"; library: "go-stdlib"; instr match { CALL text_template___Template__Execute } gen { $t.Execute($wr, $data); } bind { t as "t"; wr as "wr"; data as "data"; } }
