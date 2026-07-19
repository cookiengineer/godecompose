arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Scanner_Bytes {
    name: "bufio.Scanner.Bytes"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Bytes }
    gen { $s.Bytes(); } bind { s as "s"; }
}
