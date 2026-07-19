arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Reader_ReadLine {
    name: "bufio.Reader.ReadLine"; library: "go-stdlib"; instr match { CALL bufio___Reader__ReadLine }
    gen { $r.ReadLine(); } bind { r as "r"; }
}
