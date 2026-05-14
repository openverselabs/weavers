<p align="center">
  <img src="https://i.ibb.co.com/kV5L53BM/weavers.jpg" alt="weavers" border="0">
</p>
<h1 align="center">Weavers</h1>
<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-blue.svg" alt="Language Go">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Release-v1.0.0-orange.svg" alt="Release">
</p>

A fast and highly concurrent reconnaissance tool built in Go for gathering DNS records, WHOIS data, SSL/TLS certificate details, web technologies, and common exposed files across multiple domains simultaneously.

---

## Features

- Highly concurrent reconnaissance engine
- DNS enumeration
  - A Records
  - MX Records
  - CNAME Records
  - TXT Records
- WHOIS information gathering
- SSL/TLS certificate inspection
- Web technology fingerprinting
- robots.txt detection
- sitemap.xml detection
- Output file support
- Silent mode
- Configurable timeout & concurrency

---

## Installation

### Quick Install

```bash id="93m2xd"
curl -sSL https://raw.githubusercontent.com/openverselabs/weavers/main/install.sh | bash
```

### Build From Source

```bash id="9s3m2d"
git clone https://github.com/openverselabs/weavers.git
cd weavers
go build -o weavers
```

---

## Usage

### Single Target

```bash id="m2k1sa"
./weavers -d example.com
```

### Multiple Targets

```bash id="z0x8cn"
./weavers -l domains.txt
```

### Save Output

```bash id="v1x2ka"
./weavers -l domains.txt -o output.txt
```

### Silent Mode

```bash id="s8x2mz"
./weavers -d example.com -silent
```

### Custom Concurrency

```bash id="x2m1za"
./weavers -l domains.txt -c 50
```

### Custom Timeout

```bash id="m2s8dx"
./weavers -d example.com -t 15
```

---

## Flags

| Flag | Description |
|---|---|
| `-d` | Single target domain |
| `-l` | File containing list of domains |
| `-o` | File to write output |
| `-c` | Maximum concurrency |
| `-t` | Timeout in seconds |
| `-silent` | Silent mode |

---

## Example Output

```txt id="k1z9xm"
Summary for: example.com

[+] DNS
 - A: 93.184.216.34
 - MX: mail.example.com
 - TXT: v=spf1 include:_spf.google.com ~all

[+] WHOIS
 - Registrar: Example Registrar
 - Creation Date: 2020-01-01
 - Expiry Date: 2030-01-01

[+] SSL/TLS
 - Issuer: Let's Encrypt
 - Subject: example.com
 - Expiry: 2026-01-01 00:00:00 +0000 UTC

[+] Web Technologies & Files
 - Server: nginx
 - Powered By: PHP/8.2
 - Found: /robots.txt
 - Found: /sitemap.xml
```

---

## Example

```bash id="z8x2mn"
cat domains.txt

example.com
google.com
openai.com
```

```bash id="o2x8ma"
./weavers -l domains.txt -c 30 -o result.txt
```

---

## Disclaimer

This tool is intended for authorized reconnaissance and security testing purposes only. Unauthorized usage against systems without permission may violate applicable laws and regulations.

---


## License and Contributions

* **License**: Distributed under the MIT License.
* **Contributing**: Pull requests are welcome. For major changes, please open an issue first to discuss the proposed updates.
