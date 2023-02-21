package hashing

const (
	block_size  = 64
	result_size = 20

	ipad = 0x36
	opad = 0x5c
)

func HMAC_SHA1(message string, key string) string {
	keyByte := []byte(key) //

	if len(keyByte) > block_size {
		keyByte = []byte(SHA1(key))
	}

	if len(keyByte) < block_size {
		keyByte = append(keyByte, make([]byte, block_size-len(keyByte))...)
	}

	var ikeypad []byte
	var okeypad []byte
	for i := 0; i < len(keyByte); i++ {
		ikeypad = append(ikeypad, keyByte[i]^ipad)
		okeypad = append(okeypad, keyByte[i]^opad)
	}

	return SHA1(string(okeypad) + SHA1(string(ikeypad)+message))
}
