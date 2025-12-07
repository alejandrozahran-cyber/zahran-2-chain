use sha2::{Sha256, Digest};

pub struct AES256 {
    key: [u8; 32],
}

impl AES256 {
    pub fn new(password: &str) -> Self {
        let mut hasher = Sha256::new();
        hasher.update(password.as_bytes());
        let result = hasher.finalize();
        
        let mut key = [0u8; 32];
        key.copy_from_slice(&result[..]);
        
        AES256 { key }
    }

    pub fn encrypt(&self, data: &[u8]) -> Vec<u8> {
        // Simple XOR encryption (for demo purposes)
        // In production, use real AES-256
        data.iter()
            .enumerate()
            .map(|(i, &b)| b ^ self.key[i % 32])
            .collect()
    }

    pub fn decrypt(&self, encrypted: &[u8]) -> Vec<u8> {
        // XOR is symmetric
        self.encrypt(encrypted)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_aes256() {
        let aes = AES256::new("super_secret_password");
        let plaintext = b"Hello, NUSA Chain!";
        
        let encrypted = aes.encrypt(plaintext);
        let decrypted = aes.decrypt(&encrypted);
        
        assert_eq!(plaintext, &decrypted[..]);
    }
}
