pub struct UpgradeManager {
    current_version: String,
    pending_upgrade: Option<String>,
}

impl UpgradeManager {
    pub fn new(version: String) -> Self {
        UpgradeManager {
            current_version: version,
            pending_upgrade: None,
        }
    }

    pub fn propose_upgrade(&mut self, new_version: String) {
        println!("📢 Upgrade proposed: {} → {}", self.current_version, new_version);
        self.pending_upgrade = Some(new_version);
    }

    pub fn execute_upgrade(&mut self) -> Result<String, String> {
        if let Some(new_version) = self.pending_upgrade.take() {
            println!("✅ Executing forkless upgrade to {}", new_version);
            self.current_version = new_version. clone();
            Ok(new_version)
        } else {
            Err("No pending upgrade". to_string())
        }
    }

    pub fn current_version(&self) -> &str {
        &self.current_version
    }
}
