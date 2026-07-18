arch x86_64;
platform linux, darwin, windows, freebsd;

pattern text_template_New { name: "text/template.New"; library: "go-stdlib"; instr match { CALL text_template_New } gen { template.New($name); } bind { name as "name"; } }
pattern text_template_Must { name: "text/template.Must"; library: "go-stdlib"; instr match { CALL text_template_Must } gen { template.Must($t, $err); } bind { t as "t"; err as "err"; } }
pattern text_template_ParseFiles { name: "text/template.ParseFiles"; library: "go-stdlib"; instr match { CALL text_template_ParseFiles } gen { template.ParseFiles($filenames...); } }
pattern text_template_ParseGlob { name: "text/template.ParseGlob"; library: "go-stdlib"; instr match { CALL text_template_ParseGlob } gen { template.ParseGlob($pattern); } bind { pattern as "pattern"; } }
pattern text_template_Execute { name: "text/template.Template.Execute"; library: "go-stdlib"; instr match { CALL text_template___Template__Execute } gen { $t.Execute($wr, $data); } bind { t as "t"; wr as "wr"; data as "data"; } }
pattern text_template_ExecuteTemplate { name: "text/template.Template.ExecuteTemplate"; library: "go-stdlib"; instr match { CALL text_template___Template__ExecuteTemplate } gen { $t.ExecuteTemplate($wr, $name, $data); } bind { t as "t"; wr as "wr"; name as "name"; data as "data"; } }
