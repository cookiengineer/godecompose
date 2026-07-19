arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rsa_DecryptOAEP { name: "crypto/rsa.DecryptOAEP"; library: "go-stdlib"; instr match { CALL crypto_rsa_DecryptOAEP } gen { rsa.DecryptOAEP($hash, $random, $priv, $ciphertext, $label); } bind { hash as "hash"; random as "random"; priv as "priv"; ciphertext as "ciphertext"; label as "label"; } }
