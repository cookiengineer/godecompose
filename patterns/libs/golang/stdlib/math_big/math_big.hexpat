arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Int_SetString { name: "math/big.Int.SetString"; library: "go-stdlib"; instr match { CALL math_big___Int__SetString } gen { $z.SetString($s, $base); } bind { z as "z"; s as "s"; base as "base"; } }
pattern big_Int_String { name: "math/big.Int.String"; library: "go-stdlib"; instr match { CALL math_big___Int__String } gen { $x.String(); } bind { x as "x"; } }
pattern big_Int_Add { name: "math/big.Int.Add"; library: "go-stdlib"; instr match { CALL math_big___Int__Add } gen { $z.Add($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
pattern big_Int_Mul { name: "math/big.Int.Mul"; library: "go-stdlib"; instr match { CALL math_big___Int__Mul } gen { $z.Mul($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
pattern big_Int_Div { name: "math/big.Int.Div"; library: "go-stdlib"; instr match { CALL math_big___Int__Div } gen { $z.Div($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
pattern big_Int_Mod { name: "math/big.Int.Mod"; library: "go-stdlib"; instr match { CALL math_big___Int__Mod } gen { $z.Mod($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
pattern big_Rat_SetString { name: "math/big.Rat.SetString"; library: "go-stdlib"; instr match { CALL math_big___Rat__SetString } gen { $r.SetString($s); } bind { r as "r"; s as "s"; } }
pattern big_Float_SetString { name: "math/big.Float.SetString"; library: "go-stdlib"; instr match { CALL math_big___Float__SetString } gen { $f.SetString($s); } bind { f as "f"; s as "s"; } }
