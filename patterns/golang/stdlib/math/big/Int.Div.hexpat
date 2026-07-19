arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Int_Div { name: "math/big.Int.Div"; library: "go-stdlib"; instr match { CALL math_big___Int__Div } gen { $z.Div($x, $y); } bind { z as "z"; x as "x"; y as "y"; } }
