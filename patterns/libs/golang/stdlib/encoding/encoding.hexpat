arch x86_64;
platform linux, darwin, windows, freebsd;

pattern base64_StdEncoding_EncodeToString {
    name: "base64.StdEncoding.EncodeToString"; library: "go-stdlib";
    instr match { CALL encoding_base64___Encoding__EncodeToString }
    gen { base64.StdEncoding.EncodeToString($src); } bind { src as "src"; }
}
pattern base64_StdEncoding_DecodeString {
    name: "base64.StdEncoding.DecodeString"; library: "go-stdlib";
    instr match { CALL encoding_base64___Encoding__DecodeString }
    gen { base64.StdEncoding.DecodeString($s); } bind { s as "s"; }
}
pattern base64_URLEncoding_EncodeToString {
    name: "base64.URLEncoding.EncodeToString"; library: "go-stdlib";
    instr match { CALL encoding_base64___Encoding__EncodeToString }
    gen { base64.URLEncoding.EncodeToString($src); } bind { src as "src"; }
}

pattern hex_EncodeToString {
    name: "hex.EncodeToString"; library: "go-stdlib";
    instr match { CALL encoding_hex_EncodeToString }
    gen { hex.EncodeToString($src); } bind { src as "src"; }
}
pattern hex_DecodeString {
    name: "hex.DecodeString"; library: "go-stdlib";
    instr match { CALL encoding_hex_DecodeString }
    gen { hex.DecodeString($s); } bind { s as "s"; }
}
pattern hex_NewEncoder {
    name: "hex.NewEncoder"; library: "go-stdlib";
    instr match { CALL encoding_hex_NewEncoder }
    gen { hex.NewEncoder($w); } bind { w as "w"; }
}
pattern hex_NewDecoder {
    name: "hex.NewDecoder"; library: "go-stdlib";
    instr match { CALL encoding_hex_NewDecoder }
    gen { hex.NewDecoder($r); } bind { r as "r"; }
}

pattern xml_Marshal {
    name: "xml.Marshal"; library: "go-stdlib";
    instr match { CALL encoding_xml_Marshal }
    gen { xml.Marshal($v); } bind { v as "v"; }
}
pattern xml_Unmarshal {
    name: "xml.Unmarshal"; library: "go-stdlib";
    instr match { CALL encoding_xml_Unmarshal }
    gen { xml.Unmarshal($data, $v); } bind { data as "data"; v as "v"; }
}
pattern xml_NewEncoder {
    name: "xml.NewEncoder"; library: "go-stdlib";
    instr match { CALL encoding_xml_NewEncoder }
    gen { xml.NewEncoder($w); } bind { w as "w"; }
}
pattern xml_NewDecoder {
    name: "xml.NewDecoder"; library: "go-stdlib";
    instr match { CALL encoding_xml_NewDecoder }
    gen { xml.NewDecoder($r); } bind { r as "r"; }
}

pattern binary_Read {
    name: "binary.Read"; library: "go-stdlib";
    instr match { CALL encoding_binary_Read }
    gen { binary.Read($r, $order, $data); } bind { r as "r"; order as "order"; data as "data"; }
}
pattern binary_Write {
    name: "binary.Write"; library: "go-stdlib";
    instr match { CALL encoding_binary_Write }
    gen { binary.Write($w, $order, $data); } bind { w as "w"; order as "order"; data as "data"; }
}

pattern regexp_Compile {
    name: "regexp.Compile"; library: "go-stdlib";
    instr match { CALL regexp_Compile }
    gen { regexp.Compile($expr); } bind { expr as "expr"; }
}
pattern regexp_MustCompile {
    name: "regexp.MustCompile"; library: "go-stdlib";
    instr match { CALL regexp_MustCompile }
    gen { regexp.MustCompile($expr); } bind { expr as "expr"; }
}
pattern regexp_MatchString {
    name: "regexp.MatchString"; library: "go-stdlib";
    instr match { CALL regexp_MatchString }
    gen { regexp.MatchString($pattern, $s); } bind { pattern as "pattern"; s as "s"; }
}
pattern regexp_Regexp_MatchString {
    name: "regexp.Regexp.MatchString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__MatchString }
    gen { $re.MatchString($s); } bind { re as "re"; s as "s"; }
}
pattern regexp_Regexp_FindString {
    name: "regexp.Regexp.FindString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__FindString }
    gen { $re.FindString($s); } bind { re as "re"; s as "s"; }
}
pattern regexp_Regexp_FindAllString {
    name: "regexp.Regexp.FindAllString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__FindAllString }
    gen { $re.FindAllString($s, $n); } bind { re as "re"; s as "s"; n as "n"; }
}
pattern regexp_Regexp_ReplaceAllString {
    name: "regexp.Regexp.ReplaceAllString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__ReplaceAllString }
    gen { $re.ReplaceAllString($src, $repl); } bind { re as "re"; src as "src"; repl as "repl"; }
}
pattern regexp_Regexp_Split {
    name: "regexp.Regexp.Split"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__Split }
    gen { $re.Split($s, $n); } bind { re as "re"; s as "s"; n as "n"; }
}

pattern path_filepath_Join {
    name: "path/filepath.Join"; library: "go-stdlib";
    instr match { CALL path_filepath_Join }
    gen { filepath.Join($elem...); }
}
pattern path_filepath_Base {
    name: "path/filepath.Base"; library: "go-stdlib";
    instr match { CALL path_filepath_Base }
    gen { filepath.Base($path); } bind { path as "path"; }
}
pattern path_filepath_Dir {
    name: "path/filepath.Dir"; library: "go-stdlib";
    instr match { CALL path_filepath_Dir }
    gen { filepath.Dir($path); } bind { path as "path"; }
}
pattern path_filepath_Ext {
    name: "path/filepath.Ext"; library: "go-stdlib";
    instr match { CALL path_filepath_Ext }
    gen { filepath.Ext($path); } bind { path as "path"; }
}
pattern path_filepath_Abs {
    name: "path/filepath.Abs"; library: "go-stdlib";
    instr match { CALL path_filepath_Abs }
    gen { filepath.Abs($path); } bind { path as "path"; }
}
pattern path_filepath_Walk {
    name: "path/filepath.Walk"; library: "go-stdlib";
    instr match { CALL path_filepath_Walk }
    gen { filepath.Walk($root, $fn); } bind { root as "root"; fn as "fn"; }
}
pattern path_filepath_Glob {
    name: "path/filepath.Glob"; library: "go-stdlib";
    instr match { CALL path_filepath_Glob }
    gen { filepath.Glob($pattern); } bind { pattern as "pattern"; }
}
