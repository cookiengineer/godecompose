arch x86_64;
platform linux, darwin, windows, freebsd;

pattern xml_NewEncoder {
    name: "xml.NewEncoder"; library: "go-stdlib";
    instr match { CALL encoding_xml_NewEncoder }
    gen { xml.NewEncoder($w); } bind { w as "w"; }
}
