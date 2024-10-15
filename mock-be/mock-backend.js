// mock-backend.js
const express = require("express");
const app = express();
const PORT = 3000;

// Middleware to parse JSON bodies
app.use(express.json());

// In-memory database for simplicity
const users = {
  5655972163: {
    telegramId: "5655972163",
    name: "John Doe",
    accessToken: "mock-token-12345",
  },
};

// Endpoint to validate Telegram user and return an access token
app.post("/auth", (req, res) => {
  const { telegramId } = req.body;

  if (!telegramId) {
    return res.status(400).json({ message: "Telegram ID is required." });
  }

  const user = users[telegramId];

  if (user) {
    // User exists, return access token
    return res.status(200).json({ accessToken: user.accessToken });
  } else {
    // User not found
    return res.status(403).json({ message: "Access denied. User not found." });
  }
});

// Endpoint to validate the access token
app.get("/protected", (req, res) => {
  const authHeader = req.headers.authorization;

  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    return res
      .status(401)
      .json({ message: "Authorization token is required." });
  }

  const token = authHeader.split(" ")[1];

  // Check if the token exists in our database
  const user = Object.values(users).find((user) => user.accessToken === token);

  if (user) {
    return res.status(200).json({ message: "Access granted.", user });
  } else {
    return res.status(403).json({ message: "Invalid access token." });
  }
});

// Start the mock server
app.listen(PORT, () => {
  console.log(`Mock backend running on http://localhost:${PORT}`);
});
