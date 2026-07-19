arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_aes_NewCipher {
    name: "crypto/aes.NewCipher";
    library: "go-stdlib";
    description: "Creates a new AES cipher";

    instr match {
        MOVQ key, AX
        CALL aes_NewCipher
    }

    gen {
        aes.NewCipher($key);
    }

    bind {
        key as "key";
    }
}
