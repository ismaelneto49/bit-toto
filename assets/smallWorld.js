const fs = require("fs");

function generateSmallWorld(numNodes, extraEdges = 2, basePort = 5000) {
  const connections = {};

  const nodeName = (i) => `127.0.0.1:${basePort + i}`;

  // Initialize
  for (let i = 0; i < numNodes; i++) {
    connections[nodeName(i)] = [];
  }

  // Step 1: Ring backbone
  for (let i = 0; i < numNodes; i++) {
    const next = (i + 1) % numNodes;
    const prev = (i - 1 + numNodes) % numNodes;

    addConnection(connections, nodeName(i), nodeName(next));
    addConnection(connections, nodeName(i), nodeName(prev));
  }

  // Step 2: Add random long-distance connections
  for (let i = 0; i < numNodes; i++) {
    const current = nodeName(i);

    // Possible targets: all except self and existing neighbors
    const possibleTargets = [];
    for (let j = 0; j < numNodes; j++) {
      const candidate = nodeName(j);
      if (candidate !== current && !connections[current].includes(candidate)) {
        possibleTargets.push(candidate);
      }
    }

    // Shuffle possible targets
    shuffle(possibleTargets);

    // Pick up to `extraEdges` targets
    const chosen = possibleTargets.slice(0, extraEdges);

    for (const target of chosen) {
      addConnection(connections, current, target);
    }
  }

  return { connections };
}

// Helper: add bidirectional edge
function addConnection(connections, a, b) {
  if (!connections[a].includes(b)) connections[a].push(b);
  if (!connections[b].includes(a)) connections[b].push(a);
}

// Fisher–Yates shuffle
function shuffle(array) {
  for (let i = array.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [array[i], array[j]] = [array[j], array[i]];
  }
}

const numNodes = process.argv[2] ? parseInt(process.argv[2]) : 10;
const graph = generateSmallWorld(numNodes, 2, 5000);

fs.writeFileSync("graph.json", JSON.stringify(graph, null, 2));
console.log(`Generated small-world graph with ${numNodes} nodes -> graph.json`);
