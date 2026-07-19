arch x86_64;
platform linux, darwin, windows, freebsd;

pattern text_template_ExecuteTemplate { name: "text/template.Template.ExecuteTemplate"; library: "go-stdlib"; instr match { CALL text_template___Template__ExecuteTemplate } gen { $t.ExecuteTemplate($wr, $name, $data); } bind { t as "t"; wr as "wr"; name as "name"; data as "data"; } }
