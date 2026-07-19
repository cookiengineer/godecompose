arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_hmac_New {
    name: "crypto/hmac.New";
    library: "go-stdlib";
    description: "Creates a new HMAC hash";

    instr match {
        MOVQ h, AX
        MOVQ key, BX
        CALL hmac_New
    }

    gen {
        hmac.New($h, $key);
    }

    bind {
        h as "h";
        key as "key";
    }
}
