arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bufio_Scanner_Err {
    name: "bufio.Scanner.Err"; library: "go-stdlib"; instr match { CALL bufio___Scanner__Err }
    gen { $s.Err(); } bind { s as "s"; }
}
