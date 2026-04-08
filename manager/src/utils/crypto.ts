import CryptoJS from 'crypto-js';

const crypto = {
	// MD5加密
	MD5(data: string) {
		return CryptoJS.MD5(data).toString();
	},

	//BASE64加解密
	BASE64: {
		encrypt(data: string) {
			return CryptoJS.enc.Base64.stringify(CryptoJS.enc.Utf8.parse(data));
		},
		decrypt(data: string) {
			return CryptoJS.enc.Base64.parse(data).toString(CryptoJS.enc.Utf8);
		}
	},

	//AES加解密
	AES: {
		encrypt(data: string, secretKey: string, config : any={}){
			if(secretKey.length % 8 != 0){
				console.warn("[SCUI error]: 秘钥长度需为8的倍数，否则解密将会失败。")
			}
			const result = CryptoJS.AES.encrypt(data, CryptoJS.enc.Utf8.parse(secretKey), {
				iv: CryptoJS.enc.Utf8.parse(config.iv || ""),
				mode: CryptoJS.mode.CBC,
				padding: CryptoJS.pad.Pkcs7
			})
			return result.toString()
		},
		decrypt(cipher: string, secretKey: string, config: any={}){
			const result = CryptoJS.AES.decrypt(cipher, CryptoJS.enc.Utf8.parse(secretKey), {
				iv: CryptoJS.enc.Utf8.parse(config.iv || ""),
				mode: CryptoJS.mode.CBC,
				padding: CryptoJS.pad.Pkcs7
			})
			return CryptoJS.enc.Utf8.stringify(result);
		},
	},
}

export default crypto;