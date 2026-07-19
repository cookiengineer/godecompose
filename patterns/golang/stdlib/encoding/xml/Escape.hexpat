arch x86_64;
platform linux, darwin, windows, freebsd;

pattern xml_Escape {
    name: "xml.Escape"; library: "go-stdlib";
    instr match { CALL encoding_xml_Escape }
    gen { xml.Escape($s); } bind { s as "s"; }
}
