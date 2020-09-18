# init token
softhsm2-util --init-token --slot 0 --label "test-hsm"

# delete token
softhsm2-util --delete-token --token "test-hsm"

# pkcs11-tool
## show slot list
pkcs11-tool --module ./module/libsofthsm2.so -L

## show mechanism
pkcs11-tool --module ./module/libsofthsm2.so -M