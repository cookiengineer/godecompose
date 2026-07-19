arch x86_64;
platform linux, darwin, windows, freebsd;

pattern html_template_ExecuteTemplate { name: "html/template.Template.ExecuteTemplate"; library: "go-stdlib"; instr match { CALL html_template___Template__ExecuteTemplate } gen { $t.ExecuteTemplate($wr, $name, $data); } bind { t as "t"; wr as "wr"; name as "name"; data as "data"; } }
