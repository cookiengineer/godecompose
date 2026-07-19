arch x86_64;
platform linux, darwin, windows, freebsd;

pattern xml_MarshalIndent {
    name: "xml.MarshalIndent"; library: "go-stdlib";
    instr match { CALL encoding_xml_MarshalIndent }
    gen { xml.MarshalIndent($v, $prefix, $indent); } bind { v as "v"; prefix as "prefix"; indent as "indent"; }
}
