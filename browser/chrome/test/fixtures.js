import path from 'path';
import puppeteer from 'puppeteer';

const setup = async () => {
    const browser = await puppeteer.launch({
        // headless: true,
        // slowMo: 200,
        args: [
            `--disable-extensions-except=${path.join(process.cwd(), 'src')}`,
            `--load-extension=${path.join(process.cwd(), 'src')}`,
            '--disable-features=DialMediaRouteProvider',
        ],
    });
    const page = await browser.newPage();
    return { browser, page };
};

const teardown = async ({ browser, page }) => {
    await browser.close();
};

const getExtensionId = async (browser) => {
    const extensionTarget = await browser.waitForTarget(
        // Assumes that there is only one service worker created by the extension and its URL ends with background.js.
        target =>
            target.type() === 'service_worker' && target.url().endsWith('background.js')
    );
    const partialExtensionUrl = extensionTarget.url() || '';
    const [, , extensionId] = partialExtensionUrl.split('/');
    return extensionId;
};

export { getExtensionId, setup, teardown };

