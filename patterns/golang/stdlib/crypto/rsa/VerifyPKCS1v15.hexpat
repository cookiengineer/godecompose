arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rsa_VerifyPKCS1v15 { name: "crypto/rsa.VerifyPKCS1v15"; library: "go-stdlib"; instr match { CALL crypto_rsa_VerifyPKCS1v15 } gen { rsa.VerifyPKCS1v15($pub, $hash, $hashed, $sig); } bind { pub as "pub"; hash as "hash"; hashed as "hashed"; sig as "sig"; } }
