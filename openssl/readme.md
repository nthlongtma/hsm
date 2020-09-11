# Generate 2048 bit Private key
$ openssl genrsa -out private.pem 2048
# Separate the public part from the Private key file.
$ openssl rsa -in private.pem -pubout > public.pem