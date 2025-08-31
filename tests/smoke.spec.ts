import { test, expect } from '@playwright/test';

test('homepage has a title', async ({ page }) => {
  // Navigate to the base URL defined in the config
  await page.goto('/');

  // Expect a title to be present in the document.
  // This is a basic check to ensure the page loaded.
  await expect(page).toHaveTitle(/.+/);
});
