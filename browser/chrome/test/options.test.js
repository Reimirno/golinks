import { getExtensionId, setup, teardown } from './fixtures';

let browser, page;

beforeAll(async () => {
    const result = await setup();
    browser = result.browser;
    page = result.page;
});

afterAll(async () => {
    await teardown({ browser, page });
});

test("Displays options page", (async () => {
    const extensionId = await getExtensionId(browser);
    const extensionUrl = `chrome-extension://${extensionId}/options/index.html`;

    await page.goto(extensionUrl, { waitUntil: ['domcontentloaded', "networkidle2"] });

    const optionsHeading = await page.$eval('h2', (e => e.textContent));
    expect(optionsHeading).toEqual('Configure Golink Server');

    const serverUrlLabel = await page.$eval('label[for="serverUrl"]', (e => e.textContent));
    expect(serverUrlLabel).toEqual('Server URL:');

    const serverUrlInput = await page.$eval('#serverUrl', (e => e.placeholder));
    expect(serverUrlInput).toEqual('e.g. http://localhost:8080');

    const saveButton = await page.$eval('#saveButton', (e => e.textContent));
    expect(saveButton).toEqual('Save');
}));
