#!/usr/bin/env node
const { spawn } = require("child_process");
const path = require("path");
const os = require("os");

const platform = os.platform();
const arch = os.arch();

let binaryName;

if (platform === "linux" && arch === "x64") {
  binaryName = "langforge-linux-amd64";
} else if (platform === "linux" && arch === "arm64") {
  binaryName = "langforge-linux-arm64";
} else if (platform === "darwin" && arch === "x64") {
  binaryName = "langforge-macos-amd64";
} else if (platform === "darwin" && arch === "arm64") {
  binaryName = "langforge-macos-arm64";
} else if (platform === "win32" && arch === "x64") {
  binaryName = "langforge-windows-amd64.exe";
} else if (platform === "win32" && arch === "arm64") {
  binaryName = "langforge-windows-arm64.exe";
} else {
  console.error("Unsupported platform or architecture:", platform, arch);
  process.exit(1);
}

const binaryPath = path.join(__dirname, "bin", binaryName);
const args = process.argv.slice(2);

const child = spawn(binaryPath, args, { stdio: "inherit" });

child.on("error", (error) => {
  console.error("Failed to start subprocess:", error);
});

child.on("exit", (code, signal) => {
  if (code !== null) {
    process.exit(code);
  } else if (signal) {
    process.exit(signal === "SIGINT" ? 0 : 1);
  }
});
