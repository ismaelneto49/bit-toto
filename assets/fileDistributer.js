const fs = require("fs");
const path = require("path");

// Usage: node fileDistributer.js <numNodes> <filesFolder> <nodesFolder>
const numNodes = parseInt(process.argv[2]) || 10;
const filesFolder = process.argv[3] || "files";
const nodesFolder = process.argv[4] || "nodes";
const basePort = 5000;
const randomExtras = 0;

// Ensure folders exist
fs.mkdirSync(filesFolder, { recursive: true });
fs.mkdirSync(nodesFolder, { recursive: true });

// Step 1: Create N files in filesFolder
for (let i = 0; i < numNodes; i++) {
  const fileName = `file${i}.txt`;
  const filePath = path.join(filesFolder, fileName);
  fs.writeFileSync(filePath, fileName, "utf-8");
}

// Step 2 & 3: Create node subfolders and copy files
for (let i = 0; i < numNodes; i++) {
  const nodeName = `127.0.0.1:${basePort + i}`;
  const nodeFolder = path.join(nodesFolder, nodeName);

  fs.mkdirSync(nodeFolder, { recursive: true });

  // Copy own file
  const ownFile = `file${i}.txt`;
  fs.copyFileSync(
    path.join(filesFolder, ownFile),
    path.join(nodeFolder, ownFile)
  );

  // Pick 5 random other files
  const allFiles = Array.from({ length: numNodes }, (_, j) => `file${j}.txt`);
  const candidates = allFiles.filter((f) => f !== ownFile);

  shuffle(candidates);
  const chosen = candidates.slice(0, randomExtras);

  for (const file of chosen) {
    fs.copyFileSync(
      path.join(filesFolder, file),
      path.join(nodeFolder, file)
    );
  }
}

console.log(
  `Created ${numNodes} files in "${filesFolder}" and node folders in "${nodesFolder}".`
);

// Helper: shuffle array (Fisher–Yates)
function shuffle(arr) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [arr[i], arr[j]] = [arr[j], arr[i]];
  }
}
