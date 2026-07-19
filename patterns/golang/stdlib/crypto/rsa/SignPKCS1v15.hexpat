arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rsa_SignPKCS1v15 { name: "crypto/rsa.SignPKCS1v15"; library: "go-stdlib"; instr match { CALL crypto_rsa_SignPKCS1v15 } gen { rsa.SignPKCS1v15($random, $priv, $hash, $hashed); } bind { random as "random"; priv as "priv"; hash as "hash"; hashed as "hashed"; } }
