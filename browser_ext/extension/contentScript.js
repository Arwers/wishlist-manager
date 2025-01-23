console.log("Content script injected!");

chrome.runtime.sendMessage(
  { type: "CONTENT_SCRIPT_LOADED" },
  (response) => {
    console.log("Background responded:", response);
  }
);
