arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Contains {
    name: "strings.Contains";
    library: "go-stdlib"; instr match { CALL strings_Contains }
    gen { strings.Contains($s, $substr); }
    bind { s as "s"; substr as "substr"; }
}
pattern strings_HasPrefix {
    name: "strings.HasPrefix";
    library: "go-stdlib"; instr match { CALL strings_HasPrefix }
    gen { strings.HasPrefix($s, $prefix); }
    bind { s as "s"; prefix as "prefix"; }
}
pattern strings_HasSuffix {
    name: "strings.HasSuffix";
    library: "go-stdlib"; instr match { CALL strings_HasSuffix }
    gen { strings.HasSuffix($s, $suffix); }
    bind { s as "s"; suffix as "suffix"; }
}
pattern strings_Join {
    name: "strings.Join";
    library: "go-stdlib"; instr match { CALL strings_Join }
    gen { strings.Join($elems, $sep); }
    bind { elems as "elems"; sep as "sep"; }
}
pattern strings_Split {
    name: "strings.Split";
    library: "go-stdlib"; instr match { CALL strings_Split }
    gen { strings.Split($s, $sep); }
    bind { s as "s"; sep as "sep"; }
}
pattern strings_ReplaceAll {
    name: "strings.ReplaceAll";
    library: "go-stdlib"; instr match { CALL strings_ReplaceAll }
    gen { strings.ReplaceAll($s, $old, $new); }
    bind { s as "s"; old as "old"; new as "new"; }
}
pattern strings_ToLower {
    name: "strings.ToLower";
    library: "go-stdlib"; instr match { CALL strings_ToLower }
    gen { strings.ToLower($s); }
    bind { s as "s"; }
}
pattern strings_ToUpper {
    name: "strings.ToUpper";
    library: "go-stdlib"; instr match { CALL strings_ToUpper }
    gen { strings.ToUpper($s); }
    bind { s as "s"; }
}
pattern strings_TrimSpace {
    name: "strings.TrimSpace";
    library: "go-stdlib"; instr match { CALL strings_TrimSpace }
    gen { strings.TrimSpace($s); }
    bind { s as "s"; }
}
pattern strings_TrimPrefix {
    name: "strings.TrimPrefix";
    library: "go-stdlib"; instr match { CALL strings_TrimPrefix }
    gen { strings.TrimPrefix($s, $prefix); }
    bind { s as "s"; prefix as "prefix"; }
}
pattern strings_TrimSuffix {
    name: "strings.TrimSuffix";
    library: "go-stdlib"; instr match { CALL strings_TrimSuffix }
    gen { strings.TrimSuffix($s, $suffix); }
    bind { s as "s"; suffix as "suffix"; }
}
pattern strings_Index {
    name: "strings.Index";
    library: "go-stdlib"; instr match { CALL strings_Index }
    gen { strings.Index($s, $substr); }
    bind { s as "s"; substr as "substr"; }
}
pattern strings_Count {
    name: "strings.Count";
    library: "go-stdlib"; instr match { CALL strings_Count }
    gen { strings.Count($s, $substr); }
    bind { s as "s"; substr as "substr"; }
}
pattern strings_Repeat {
    name: "strings.Repeat";
    library: "go-stdlib"; instr match { CALL strings_Repeat }
    gen { strings.Repeat($s, $count); }
    bind { s as "s"; count as "count"; }
}
pattern strings_Fields {
    name: "strings.Fields";
    library: "go-stdlib"; instr match { CALL strings_Fields }
    gen { strings.Fields($s); }
    bind { s as "s"; }
}
pattern strings_NewReader {
    name: "strings.NewReader";
    library: "go-stdlib"; instr match { CALL strings_NewReader }
    gen { strings.NewReader($s); }
    bind { s as "s"; }
}
pattern strings_NewReplacer {
    name: "strings.NewReplacer";
    library: "go-stdlib"; instr match { CALL strings_NewReplacer }
    gen { strings.NewReplacer($oldnew...); }
    bind { oldnew as "oldnew"; }
}
pattern strings_Builder_WriteString {
    name: "strings.Builder.WriteString";
    library: "go-stdlib"; instr match { CALL strings___Builder__WriteString }
    gen { $b.WriteString($s); }
    bind { b as "b"; s as "s"; }
}
pattern strings_Builder_WriteByte {
    name: "strings.Builder.WriteByte";
    library: "go-stdlib"; instr match { CALL strings___Builder__WriteByte }
    gen { $b.WriteByte($c); }
    bind { b as "b"; c as "c"; }
}
pattern strings_Builder_String {
    name: "strings.Builder.String";
    library: "go-stdlib"; instr match { CALL strings___Builder__String }
    gen { $b.String(); }
    bind { b as "b"; }
}
pattern strings_Builder_Reset {
    name: "strings.Builder.Reset";
    library: "go-stdlib"; instr match { CALL strings___Builder__Reset }
    gen { $b.Reset(); }
    bind { b as "b"; }
}
