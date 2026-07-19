arch x86_64;
platform linux, darwin, windows, freebsd;

pattern xml_EscapeText {
    name: "xml.EscapeText"; library: "go-stdlib";
    instr match { CALL encoding_xml_EscapeText }
    gen { xml.EscapeText($w, $text); } bind { w as "w"; text as "text"; }
}
