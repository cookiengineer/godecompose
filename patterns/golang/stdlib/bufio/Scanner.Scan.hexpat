arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Scanner_Scan {
    name: "bufio.Scanner.Scan"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Scan }
    gen { $s.Scan(); } bind { s as "s"; }
}
