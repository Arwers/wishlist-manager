chrome.runtime.onInstalled.addListener(() => {
    console.log("Extension Installed");
  });
  
  chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.type === "PING_SERVER") {
      console.log("Background script received PING_SERVER request");
      sendResponse({ message: "Quickpick" });
    }
  });
  