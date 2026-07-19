arch x86_64;
platform linux, darwin, windows, freebsd;

pattern html_template_Execute { name: "html/template.Template.Execute"; library: "go-stdlib"; instr match { CALL html_template___Template__Execute } gen { $t.Execute($wr, $data); } bind { t as "t"; wr as "wr"; data as "data"; } }
