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
pattern base64_RawStdEncoding_EncodeToString {
    name: "base64.RawStdEncoding.EncodeToString"; library: "go-stdlib";
    instr match { CALL encoding_base64___Encoding__EncodeToString }
    gen { base64.RawStdEncoding.EncodeToString($src); } bind { src as "src"; }
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
pattern hex_Dump {
    name: "hex.Dump"; library: "go-stdlib";
    instr match { CALL encoding_hex_Dump }
    gen { hex.Dump($data); } bind { data as "data"; }
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
pattern xml_MarshalIndent {
    name: "xml.MarshalIndent"; library: "go-stdlib";
    instr match { CALL encoding_xml_MarshalIndent }
    gen { xml.MarshalIndent($v, $prefix, $indent); } bind { v as "v"; prefix as "prefix"; indent as "indent"; }
}
pattern xml_Escape {
    name: "xml.Escape"; library: "go-stdlib";
    instr match { CALL encoding_xml_Escape }
    gen { xml.Escape($s); } bind { s as "s"; }
}
pattern xml_EscapeText {
    name: "xml.EscapeText"; library: "go-stdlib";
    instr match { CALL encoding_xml_EscapeText }
    gen { xml.EscapeText($w, $text); } bind { w as "w"; text as "text"; }
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
pattern binary_Size {
    name: "binary.Size"; library: "go-stdlib";
    instr match { CALL encoding_binary_Size }
    gen { binary.Size($v); } bind { v as "v"; }
}
pattern binary_Varint {
    name: "binary.Varint"; library: "go-stdlib";
    instr match { CALL encoding_binary_Varint }
    gen { binary.Varint($buf); } bind { buf as "buf"; }
}
pattern binary_PutUvarint {
    name: "binary.PutUvarint"; library: "go-stdlib";
    instr match { CALL encoding_binary_PutUvarint }
    gen { binary.PutUvarint($buf, $x); } bind { buf as "buf"; x as "x"; }
}
pattern binary_PutVarint {
    name: "binary.PutVarint"; library: "go-stdlib";
    instr match { CALL encoding_binary_PutVarint }
    gen { binary.PutVarint($buf, $x); } bind { buf as "buf"; x as "x"; }
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
pattern regexp_CompilePOSIX {
    name: "regexp.CompilePOSIX"; library: "go-stdlib";
    instr match { CALL regexp_CompilePOSIX }
    gen { regexp.CompilePOSIX($expr); } bind { expr as "expr"; }
}
pattern regexp_Match {
    name: "regexp.Match"; library: "go-stdlib";
    instr match { CALL regexp_Match }
    gen { regexp.Match($pattern, $b); } bind { pattern as "pattern"; b as "b"; }
}
pattern regexp_QuoteMeta {
    name: "regexp.QuoteMeta"; library: "go-stdlib";
    instr match { CALL regexp_QuoteMeta }
    gen { regexp.QuoteMeta($s); } bind { s as "s"; }
}
pattern regexp_Regexp_ReplaceAllStringFunc {
    name: "regexp.Regexp.ReplaceAllStringFunc"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__ReplaceAllStringFunc }
    gen { $re.ReplaceAllStringFunc($src, $repl); } bind { re as "re"; src as "src"; repl as "repl"; }
}
pattern regexp_Regexp_SubexpNames {
    name: "regexp.Regexp.SubexpNames"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__SubexpNames }
    gen { $re.SubexpNames(); } bind { re as "re"; }
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
pattern path_filepath_Clean {
    name: "path/filepath.Clean"; library: "go-stdlib";
    instr match { CALL path_filepath_Clean }
    gen { filepath.Clean($path); } bind { path as "path"; }
}
pattern path_filepath_Rel {
    name: "path/filepath.Rel"; library: "go-stdlib";
    instr match { CALL path_filepath_Rel }
    gen { filepath.Rel($basepath, $targpath); } bind { basepath as "basepath"; targpath as "targpath"; }
}
pattern path_filepath_Split {
    name: "path/filepath.Split"; library: "go-stdlib";
    instr match { CALL path_filepath_Split }
    gen { filepath.Split($path); } bind { path as "path"; }
}
pattern path_filepath_WalkDir {
    name: "path/filepath.WalkDir"; library: "go-stdlib";
    instr match { CALL path_filepath_WalkDir }
    gen { filepath.WalkDir($root, $fn); } bind { root as "root"; fn as "fn"; }
}
pattern path_filepath_Match {
    name: "path/filepath.Match"; library: "go-stdlib";
    instr match { CALL path_filepath_Match }
    gen { filepath.Match($pattern, $name); } bind { pattern as "pattern"; name as "name"; }
}
pattern path_filepath_IsAbs {
    name: "path/filepath.IsAbs"; library: "go-stdlib";
    instr match { CALL path_filepath_IsAbs }
    gen { filepath.IsAbs($path); } bind { path as "path"; }
}
pattern path_filepath_EvalSymlinks {
    name: "path/filepath.EvalSymlinks"; library: "go-stdlib";
    instr match { CALL path_filepath_EvalSymlinks }
    gen { filepath.EvalSymlinks($path); } bind { path as "path"; }
}
