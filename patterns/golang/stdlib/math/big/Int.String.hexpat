arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Int_String { name: "math/big.Int.String"; library: "go-stdlib"; instr match { CALL math_big___Int__String } gen { $x.String(); } bind { x as "x"; } }
