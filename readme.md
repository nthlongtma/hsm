# init token
softhsm2-util --init-token --slot 0 --label "test-hsm"

# delete token
softhsm2-util --delete-token --token "test-hsm"