arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rsa_GenerateKey { name: "crypto/rsa.GenerateKey"; library: "go-stdlib"; instr match { CALL crypto_rsa_GenerateKey } gen { rsa.GenerateKey($random, $bits); } bind { random as "random"; bits as "bits"; } }
