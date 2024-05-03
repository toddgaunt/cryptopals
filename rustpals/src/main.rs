use std::str;

const BASE64_TABLE: &[u8] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/".as_bytes();

const RED: &str = "\x1b[31m";
const GREEN: &str = "\x1b[32m";
const CLEAR: &str = "\x1b[0m";

fn base64_encode(input: &[u8]) -> Result<Vec<u8>, &str> {
	if input.len()%3 != 0 {
		return Err("input length must be divisible by 3")
	}

	let mut output = Vec::new();
	let mut i = 0;
	while i < input.len() {
		let mut chunk: [u8; 4] = [0; 4];

		// First 6 bits
		chunk[0] = (input[i+0] & 0xFC) >> 2;
		// 2 bits + 4 bits = 6 bits
		chunk[1] = ((input[i+0] & 0x03) << 4) | ((input[i+1] & 0xF0) >> 4);
		// 4 bits + 2 bits = 6 bits
		chunk[2] = ((input[i+1] & 0x0F) << 2) | ((input[i+2] & 0xC0) >> 6);
		// Last 6 bits
		chunk[3] = input[i+2] & 0x3F;

		output.push(BASE64_TABLE[chunk[0] as usize]);
		output.push(BASE64_TABLE[chunk[1] as usize]);
		output.push(BASE64_TABLE[chunk[2] as usize]);
		output.push(BASE64_TABLE[chunk[3] as usize]);

		i += 3
	}
	Ok(output)
}

fn chunk_str(input: &str) -> Vec<(char, char)> {
	let bytes = input.as_bytes();
	let mut v = Vec::new();
	let mut i = 0;
	while i < bytes.len() {
		v.push((bytes[i] as char, bytes[i+1] as char));
		i += 2
	};
	v

}

fn hex2byte(ch: char) -> Result<u8, String> {
	let c = ch as u8;
	if 48 <= c && c <= 57 {
		return Ok(c - 48)
	} else if 65 <= c && c <=70 {
		return Ok(c - 55)
	} else if 97 <= c && c <= 102 {
		return Ok(c - 87)
	}
	return Err(String::from(format!("invalid hex character {}", ch)))
}

fn hex_encode(input: Vec<u8>) -> String {
	let mut output = String::new();
	for b in input {
		output.push_str(&format!("{:02x}", b));
	}
	output
}

fn hex_decode(input: &str) -> Result<Vec<u8>, String> {
	if input.len()%2 != 0 {
		return Err(String::from("input length must be divisible by 2"))
	}

	let mut output = Vec::new();
	let chunks = chunk_str(input);
	for chunk in chunks {
		let c1 = hex2byte(chunk.0);
		let c2 = hex2byte(chunk.1);
		match (c1, c2) {
			(Ok(a), Ok(b)) => output.push((a << 4) | b),
			(Err(a), Ok(_)) => return Err(a),
			(Ok(_), Err(b)) => return Err(b),
			(Err(_), Err(_)) => return Err(String::from(format!("invalid hex characters {} and {}", chunk.0, chunk.1))),
		}
	}

	Ok(output)
}

fn xor_buf(x: Vec<u8>, y: Vec<u8>) -> Result<Vec<u8>, String> {
	if x.len() != y.len() {
		return Err(String::from("buffers must be the same length"))
	}

	let mut output = Vec::new();
	for i in 0..x.len() {
		output.push(x[i] ^ y[i]);
	}

	return Ok(output)
}

fn print_result(n: i32, result: Result<(), String>) {
	match result {
		Err(s) => print!("Challenge {n}: {RED}{s}{CLEAR}\n"),
		Ok(_) => print!("Challenge {n}: {GREEN}OK{CLEAR}\n"),
	}
}

fn challenge1() -> Result<(), String> {
	let input = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d";
	let decoded = hex_decode(input);
	let bytes = decoded.unwrap();
	let encoded = base64_encode(bytes.as_slice());
	let encoded = encoded.unwrap();
	let encoded_str = str::from_utf8(&encoded);
	let got = encoded_str.unwrap();

	let want =  "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t";
	if want != got {
		return Err(String::from(format!("got {}, want {}", got, want)))
	}

	return Ok(())
}

fn challenge2() -> Result<(), String> {
	let a = "1c0111001f010100061a024b53535009181c";
	let b = "686974207468652062756c6c277320657965";
	let c = "746865206b696420646f6e277420706c6179";

	let x = hex_decode(a);
	let y = hex_decode(b);
	let z = xor_buf(x?, y?);

	let got = hex_encode(z?);
	let want = c;

	if got != want {
		return Err(String::from(format!("got {}, want {}", got, want)))
	}

	return Ok(())
}

fn main() {
	print!("Rustpals!\n");

	let result1 = challenge1();
	print_result(1, result1);

	let result2 = challenge2();
	print_result(2, result2);
}
