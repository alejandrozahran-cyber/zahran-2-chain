use std::collections::HashMap;

pub struct Proposal {
    pub id: u64,
    pub title: String,
    pub description: String,
    pub votes_for: u64,
    pub votes_against: u64,
    pub status: ProposalStatus,
}

pub enum ProposalStatus {
    Active,
    Passed,
    Rejected,
    Executed,
}

pub struct Governance {
    proposals: HashMap<u64, Proposal>,
    next_proposal_id: u64,
}

impl Governance {
    pub fn new() -> Self {
        Governance {
            proposals: HashMap::new(),
            next_proposal_id: 1,
        }
    }

    pub fn create_proposal(&mut self, title: String, description: String) -> u64 {
        let id = self.next_proposal_id;
        self.next_proposal_id += 1;
        
        let proposal = Proposal {
            id,
            title,
            description,
            votes_for: 0,
            votes_against: 0,
            status: ProposalStatus::Active,
        };
        
        self.proposals.insert(id, proposal);
        id
    }

    pub fn vote(&mut self, proposal_id: u64, vote_for: bool, power: u64) -> Result<(), String> {
        if let Some(proposal) = self.proposals.get_mut(&proposal_id) {
            if vote_for {
                proposal. votes_for += power;
            } else {
                proposal.votes_against += power;
            }
            Ok(())
        } else {
            Err("Proposal not found".to_string())
        }
    }

    pub fn get_proposal(&self, proposal_id: u64) -> Option<&Proposal> {
        self. proposals.get(&proposal_id)
    }

    pub fn execute_proposal(&mut self, proposal_id: u64) -> Result<(), String> {
        if let Some(proposal) = self.proposals.get_mut(&proposal_id) {
            if proposal.votes_for > proposal.votes_against {
                proposal.status = ProposalStatus::Executed;
                Ok(())
            } else {
                Err("Proposal not passed".to_string())
            }
        } else {
            Err("Proposal not found".to_string())
        }
    }
}
