arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Int_Add { name: "math/big.Int.Add"; library: "go-stdlib"; instr match { CALL math_big___Int__Add } gen { $z.Add($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
