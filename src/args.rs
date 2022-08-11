use clap::{Args, Parser, Subcommand};

#[derive(Parser, Debug)]
#[clap(author, version, about)]
pub struct AuthArgs {
    #[clap(subcommand)]
    pub entity_type: EntityType,
}

#[derive(Debug, Subcommand)]
pub enum EntityType {
    /// Create a new key
    Create(CreateKey),

    /// Test key authentication
    Test(TestKey),
}

#[derive(Debug, Args)]
pub struct CreateKey {
    /// The email of the user
    pub email: String,
}

#[derive(Debug, Args)]
pub struct TestKey {
    /// User given secret
    pub secret: String,
}
