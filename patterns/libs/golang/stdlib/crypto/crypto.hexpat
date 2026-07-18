arch x86_64;
platform linux, darwin, windows, freebsd;

pattern crypto_sha256_Sum256 {
    name: "crypto/sha256.Sum256";
    library: "go-stdlib";
    description: "Computes SHA256 hash";

    instr match {
        MOVQ data, AX
        CALL sha256_Sum256
    }

    gen {
        sha256.Sum256($data);
    }

    bind {
        data as "data";
    }
}

pattern crypto_sha256_New {
    name: "crypto/sha256.New";
    library: "go-stdlib";
    description: "Creates a new SHA256 hash";

    instr match {
        CALL sha256_New
    }

    gen {
        sha256.New();
    }
}

pattern crypto_md5_Sum {
    name: "crypto/md5.Sum";
    library: "go-stdlib";
    description: "Computes MD5 hash";

    instr match {
        MOVQ data, AX
        CALL md5_Sum
    }

    gen {
        md5.Sum($data);
    }

    bind {
        data as "data";
    }
}

pattern crypto_rand_Read {
    name: "crypto/rand.Read";
    library: "go-stdlib";
    description: "Reads cryptographically secure random bytes";

    instr match {
        MOVQ buf, AX
        CALL crypto_rand_Read
    }

    gen {
        crypto_rand.Read($buf);
    }

    bind {
        buf as "buf";
    }
}

pattern crypto_tls_Dial {
    name: "crypto/tls.Dial";
    library: "go-stdlib";
    description: "Dials a TLS connection";

    instr match {
        MOVQ network, AX
        MOVQ addr, BX
        MOVQ config, CX
        CALL tls_Dial
    }

    gen {
        tls.Dial($network, $addr, $config);
    }

    bind {
        network as "network";
        addr as "addr";
        config as "config";
    }
}

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

pattern crypto_x509_ParseCertificate {
    name: "crypto/x509.ParseCertificate";
    library: "go-stdlib";
    description: "Parses an X.509 certificate";

    instr match {
        MOVQ asn1, AX
        CALL x509_ParseCertificate
    }

    gen {
        x509.ParseCertificate($asn1);
    }

    bind {
        asn1 as "data";
    }
}
