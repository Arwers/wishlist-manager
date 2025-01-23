document.getElementById("pingBtn").addEventListener("click", async () => {
    const resultDiv = document.getElementById("result");
    try {
      const response = await fetch("http://localhost:3000/ping");
      const data = await response.json();
      resultDiv.textContent = `Server says: ${data.message}`;
    } catch (error) {
      console.error("Error fetching from server:", error);
      resultDiv.textContent = "Failed to ping server.";
    }
  });
  