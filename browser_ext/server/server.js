const express = require("express");
const cors = require("cors");

const app = express();
const PORT = 3000;

app.use(cors());

app.get("/ping", (req, res) => {
  res.json({ message: "pong" });
});

app.listen(PORT, () => {
  console.log(`Server listening on http://localhost:${PORT}`);
});
