chrome.webNavigation.onBeforeNavigate.addListener(function (details) {
    const url = new URL(details.url);

    // match 'go/keyword', but does not match 'go/'
    if (url.hostname === 'go' && url.pathname.length > 1) {
        // Extract the keyword after 'go/'
        const keyword = url.pathname.substring(1);

        // Get the server URL from Chrome storage
        chrome.storage.sync.get('serverUrl', function (data) {
            let serverUrl = data.serverUrl || 'http://localhost:8080'; // Default to localhost:8080 if not set

            // Construct the redirect URL using the configured server
            const redirectUrl = `${serverUrl}/${keyword}`;

            // Update the tab with the redirect URL
            chrome.tabs.update(details.tabId, { url: redirectUrl });
        });
    }
}, { url: [{ urlMatches: 'http://go/*' }] });
