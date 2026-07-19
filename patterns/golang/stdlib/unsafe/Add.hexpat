arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_Add { name: "unsafe.Add"; library: "go-stdlib"; instr match { CALL unsafe_Add } gen { unsafe.Add($p, $len); } bind { p as "p"; len as "len"; } }
