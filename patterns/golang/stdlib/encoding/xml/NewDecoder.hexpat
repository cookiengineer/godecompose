arch x86_64;
platform linux, darwin, windows, freebsd;

pattern xml_NewDecoder {
    name: "xml.NewDecoder"; library: "go-stdlib";
    instr match { CALL encoding_xml_NewDecoder }
    gen { xml.NewDecoder($r); } bind { r as "r"; }
}
