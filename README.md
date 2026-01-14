# Port-Scanner

A simple and fast **TCP port scanner written in Go**, designed to scan a target host over a user-defined port range and report open ports with a clean terminal interface.

This project was built to explore **network programming**, **concurrency**, and **system-level tooling** using Golang.

---

## Features

- Fast TCP port scanning  
- Custom port range selection  
- Interactive CLI prompts  
- Scan progress indicator  
- ASCII banner for clear startup output  
- Written in pure Go (no external dependencies)

---

## Project Structure
```
.
├── portscanner.go # Main application logic
├── go.mod # Go module definition
├── README.md # Project documentation
└── .idea/ # IDE configuration (optional)
`
```

---

## 🛠️ Requirements

- **Go 1.20+** installed  
    https://go.dev/dl/

---

## 🚀 Getting Started

### 1️⃣ Clone the repository
```bash
git clone https://github.com/Mohamed-Elwasila/Port-Scanner.git
cd Port-Scanner
```
### 2️⃣ Run the scanner
```bash
go run portscanner.go
```

