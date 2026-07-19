arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Int_Mul { name: "math/big.Int.Mul"; library: "go-stdlib"; instr match { CALL math_big___Int__Mul } gen { $z.Mul($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
