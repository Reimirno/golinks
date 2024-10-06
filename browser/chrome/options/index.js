// Save the server URL when the user clicks "Save"
document.getElementById('saveButton').addEventListener('click', function () {
    const serverUrl = document.getElementById('serverUrl').value;
    chrome.storage.sync.set({ serverUrl: serverUrl }, function () {
        alert(`Server URL saved to be: ${serverUrl}`);
    });
});

// Load the stored server URL when the options page is opened
chrome.storage.sync.get('serverUrl', function (data) {
    if (data.serverUrl) {
        document.getElementById('serverUrl').value = data.serverUrl;
    }
});
