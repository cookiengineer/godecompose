arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Reader_ReadBytes {
    name: "bufio.Reader.ReadBytes"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadBytes }
    gen { $r.ReadBytes($delim); } bind { r as "r"; delim as "delim"; }
}
