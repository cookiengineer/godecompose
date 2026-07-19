arch x86_64;
platform linux, darwin, windows, freebsd;

pattern xml_Unmarshal {
    name: "xml.Unmarshal"; library: "go-stdlib";
    instr match { CALL encoding_xml_Unmarshal }
    gen { xml.Unmarshal($data, $v); } bind { data as "data"; v as "v"; }
}
