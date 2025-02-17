package keystore_test

import (
	"testing"
	"testing/fstest"

	"github.com/google/uuid"
	"github.com/hamidoujand/sales/pkg/keystore"
)

const privatePEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAvQUB+/al4VKJFmmpx+6mFePHeZo2a43gu5rgd7qnV/c5+ZbA
KdPyuiCoUPLjfNELwSOQVi/mt/vxlnvlCwduFeZvozGP3+qEh6bVrPQF/A4RFG0f
RG9CIH0TbfjVtGgLMQMdxbXzVhrwGT0uRBlUeqwDd7e0eqB2uIL7EA8Sc6OthHkf
SMdvoUKkqWEA2xdrkuNj2aHCcxi9jC0mxa6jN0CXB0MKBLcJqAsPnPW2366aMmk8
PziIRO66HKTNMQcmAaHTII2o66fgQg0Fdj0dvM9bzY1UE2hdP9TMYzFa5721HTx4
pezu4SbER4KCO8oirGGwljeiZDreoP8OpkYgOwIDAQABAoIBAB35Q9HJUI+1D2Ea
+13lhbfV6YVqg3O1yWvmiO7jjfLglPRzx+A6KHUEhbxkb9eUrMkBUzufl/YYATzs
Q6tmj7nwU0atLtQCs+Zw+dRV0/ce4e17ymgHPpS5UNHxEi5sC05H4Lo/+qjuV6Gg
9ou8+o0DZv9ehcOmW30x5A8tXK8ygOqhturnR9HhvpyikNBr1ij6r2Y9y5xxXQiv
a79ghEV/i7OVhQ4hWZxlGVaMH0bYUxGlpsLdOrTT7awrN67gv1Dvk8edZGjJuMXV
2Rju+XTkwZlkgxAGhwEQmhb0bEibbhBf2ggQvyVm9qjXaoAM881aR1DJHYjYNiWU
18nLtBECgYEA7tEbcSwN/WkEdXFZl2DWFGR5KESrbkpWuL9GoBYI2s4FBAuX1JrS
qj0qhTbjxhuWFNJhgEkSKNucN3T6IMk3KudnuqzNgWLPlhvmrd+KSd1CF+849H0A
1F7CI9pw6GJpuBU+S8xDFj+gBDQbN/WqgQa5GhO7z1C9gWT3U3iw59ECgYEAyp6o
lpv5cAX6NZ03sboWnMMd5clT1v4IReeUTWT3N2U1rpGquZEZGzoWA3FBJ31F4pNJ
T1ddZsNUGWUlyAEuHJYyJWE8C9eq8ELhoS+CLgvP9aRqEGDywrDgIcbzPCrCUOlz
KQvlWs0ZhoAoLxRqP7XjyiB7fxP9Flhrj2nyVksCgYEAwg69I9tOiuq5Ks2upWmU
zAFQyj3yp65Uhc84DoGZNGNQhBb/i007fgYx9QnDUIm+DLFfdSTrUrQRXqb5UYbw
AzcCfRhJ7adjU3Dco9EPyDG4sUY8m76v2+IcE5I/STYe/eyVMHaM1RliZ8gHjhNc
N3hFFUGPzUiolOp8ZyGdbuECgYEAqZKfBW0EFPzrqnMpaVSUGB4zp1wXDpcL1XU6
aItXWsUZaEAA4czNdjvmsHrYTHRLSJR7hitXv+k5OQet1vUl4kbRMPdviXm1Vd6j
doKMMH0yTiKLoamBge8FpT8b0f73IUA/YNrT2GpOMoKPHte3FBrlyQPmVzQjW9Ak
NKI2boECgYBpm80kcnmHk3kGYUa99x+zSS0fEqDlJF6X2RKZIkPpov+ERmhLfKcx
3o3IUEmwP81KudbPo/lbdmAMOY0CfCRBvkeY0cC3GiJoXnjMVy4RvO0pnI4GhWZo
GyYYKBpGB9Ng99Adgc7huA7yJ1CF85BSIY2+ka7KZf02aSKGx+ALGg==
-----END RSA PRIVATE KEY-----`

func TestKeyStore(t *testing.T) {
	kid := uuid.NewString()
	filename := kid + "-private.pem"

	file := fstest.MapFile{
		Data: []byte(privatePEM),
	}

	fs := fstest.MapFS{
		filename: &file,
	}

	ks := keystore.New()
	if err := ks.LoadKeys(fs); err != nil {
		t.Fatalf("failed to load files: %s", err)
	}

	_, err := ks.PrivateKey(kid)
	if err != nil {
		t.Fatalf("failed to get private key: %s", err)
	}

	_, err = ks.PublicKey(kid)
	if err != nil {
		t.Fatalf("failed to get public key: %s", err)
	}
}
