arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Int_SetString { name: "math/big.Int.SetString"; library: "go-stdlib"; instr match { CALL math_big___Int__SetString } gen { $z.SetString($s, $base); } bind { z as "z"; s as "s"; base as "base"; } }
