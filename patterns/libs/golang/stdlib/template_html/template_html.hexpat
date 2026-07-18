arch x86_64;
platform linux, darwin, windows, freebsd;

pattern html_template_New { name: "html/template.New"; library: "go-stdlib"; instr match { CALL html_template_New } gen { template.New($name); } bind { name as "name"; } }
pattern html_template_Must { name: "html/template.Must"; library: "go-stdlib"; instr match { CALL html_template_Must } gen { template.Must($t, $err); } bind { t as "t"; err as "err"; } }
pattern html_template_ParseFiles { name: "html/template.ParseFiles"; library: "go-stdlib"; instr match { CALL html_template_ParseFiles } gen { template.ParseFiles($filenames...); } }
pattern html_template_ParseGlob { name: "html/template.ParseGlob"; library: "go-stdlib"; instr match { CALL html_template_ParseGlob } gen { template.ParseGlob($pattern); } bind { pattern as "pattern"; } }
pattern html_template_Execute { name: "html/template.Template.Execute"; library: "go-stdlib"; instr match { CALL html_template___Template__Execute } gen { $t.Execute($wr, $data); } bind { t as "t"; wr as "wr"; data as "data"; } }
pattern html_template_ExecuteTemplate { name: "html/template.Template.ExecuteTemplate"; library: "go-stdlib"; instr match { CALL html_template___Template__ExecuteTemplate } gen { $t.ExecuteTemplate($wr, $name, $data); } bind { t as "t"; wr as "wr"; name as "name"; data as "data"; } }
