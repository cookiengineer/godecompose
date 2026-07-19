arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Int_Mod { name: "math/big.Int.Mod"; library: "go-stdlib"; instr match { CALL math_big___Int__Mod } gen { $z.Mod($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
