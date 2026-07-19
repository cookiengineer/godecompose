arch x86_64;
platform linux, darwin, windows, freebsd;

pattern xml_Marshal {
    name: "xml.Marshal"; library: "go-stdlib";
    instr match { CALL encoding_xml_Marshal }
    gen { xml.Marshal($v); } bind { v as "v"; }
}
