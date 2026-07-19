arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Reader_ReadString {
    name: "bufio.Reader.ReadString"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadString }
    gen { $r.ReadString($delim); } bind { r as "r"; delim as "delim"; }
}
