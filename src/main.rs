use args::{CreateKey, TestKey};
use clap::Parser;
use google_authenticator::{
    create_secret, qr_code_url, verify_code, ErrorCorrectionLevel, GA_AUTH,
};
use std::io::{stdin, stdout, Write};

use crate::args::{AuthArgs, EntityType};

mod args;

fn main() {
    let args = AuthArgs::parse();

    match args.entity_type {
        EntityType::Create(args) => handle_create_command(args),
        EntityType::Test(args) => handle_test_command(args),
    };
}

fn handle_create_command(args: CreateKey) {
    let secret = create_secret!();

    let qr_url = qr_code_url!(
        secret.as_str(),
        args.email.as_str(),
        "MyAuth",
        200,
        200,
        ErrorCorrectionLevel::High
    );

    println!("Your secret: {}", secret);
    println!("QR Code: {}", qr_url);

    // let qrcode = qr_code!(
    //     secret.as_str(),
    //     args.email.as_str(),
    //     "MyAuth",
    //     200,
    //     200,
    //     ErrorCorrectionLevel::High
    // )
    // .unwrap();
    // println!("svg: {}", qrcode);
}

fn handle_test_command(args: TestKey) {
    println!("Enter 000000 to exit");
    loop {
        print!("Enter code> ");
        stdout().flush().unwrap();
        let mut user_text = String::new();

        stdin()
            .read_line(&mut user_text)
            .expect("Failed to read line");

        let code = user_text.trim(); // trim out new line

        if code.eq("000000") {
            println!("B Y E !!!");
            break;
        }

        if verify_code!(args.secret.as_str(), code) {
            println!("[{}] => 💪  🧤  🍏", code)
        } else {
            println!("[{}] => 🤦‍♀️  🚨  🧨", code)
        }
        println!("");
    }
}
