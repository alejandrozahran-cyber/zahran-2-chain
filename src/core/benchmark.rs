use std::sync::Arc;

pub struct Benchmark {
    theoretical_tps: u64,
}

impl Benchmark {
    pub fn new() -> Self {
        Benchmark {
            theoretical_tps: 50_000,
        }
    }

    pub fn theoretical_tps(&self) -> u64 {
        self.theoretical_tps
    }

    pub fn block_time_ms(&self) -> u64 {
        500
    }
}
