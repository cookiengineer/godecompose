arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rsa_GenerateKey { name: "crypto/rsa.GenerateKey"; library: "go-stdlib"; instr match { CALL crypto_rsa_GenerateKey } gen { rsa.GenerateKey($random, $bits); } bind { random as "random"; bits as "bits"; } }
pattern rsa_EncryptOAEP { name: "crypto/rsa.EncryptOAEP"; library: "go-stdlib"; instr match { CALL crypto_rsa_EncryptOAEP } gen { rsa.EncryptOAEP($hash, $random, $pub, $msg, $label); } bind { hash as "hash"; random as "random"; pub as "pub"; msg as "msg"; label as "label"; } }
pattern rsa_DecryptOAEP { name: "crypto/rsa.DecryptOAEP"; library: "go-stdlib"; instr match { CALL crypto_rsa_DecryptOAEP } gen { rsa.DecryptOAEP($hash, $random, $priv, $ciphertext, $label); } bind { hash as "hash"; random as "random"; priv as "priv"; ciphertext as "ciphertext"; label as "label"; } }
pattern rsa_SignPKCS1v15 { name: "crypto/rsa.SignPKCS1v15"; library: "go-stdlib"; instr match { CALL crypto_rsa_SignPKCS1v15 } gen { rsa.SignPKCS1v15($random, $priv, $hash, $hashed); } bind { random as "random"; priv as "priv"; hash as "hash"; hashed as "hashed"; } }
pattern rsa_VerifyPKCS1v15 { name: "crypto/rsa.VerifyPKCS1v15"; library: "go-stdlib"; instr match { CALL crypto_rsa_VerifyPKCS1v15 } gen { rsa.VerifyPKCS1v15($pub, $hash, $hashed, $sig); } bind { pub as "pub"; hash as "hash"; hashed as "hashed"; sig as "sig"; } }
