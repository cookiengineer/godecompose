arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Reader_ReadByte {
    name: "bufio.Reader.ReadByte"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadByte }
    gen { $r.ReadByte(); } bind { r as "r"; }
}
