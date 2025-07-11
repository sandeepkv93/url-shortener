import { test, expect } from '@playwright/test'

test.describe('Component Screenshots', () => {
  test.beforeEach(async ({ page }) => {
    // Set viewport size for consistent screenshots
    await page.setViewportSize({ width: 1280, height: 720 })
  })

  test('Home page screenshot', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    
    // Wait for fonts to load
    await page.waitForTimeout(1000)
    
    await page.screenshot({
      path: 'screenshots/01-home-page.png',
      fullPage: true
    })
  })

  test('Demo page screenshots', async ({ page }) => {
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)

    // Take screenshot of demo page navigation
    await page.screenshot({
      path: 'screenshots/02-demo-navigation.png',
      fullPage: false
    })

    // Login Form
    await page.click('text=Login Form')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/03-login-form.png',
      fullPage: true
    })

    // Register Form
    await page.click('text=Register Form')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/04-register-form.png',
      fullPage: true
    })

    // Password Reset
    await page.click('text=Password Reset')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/05-password-reset.png',
      fullPage: true
    })

    // Common Components
    await page.click('text=Common Components')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/06-common-components.png',
      fullPage: true
    })
  })

  test('Login form interactions', async ({ page }) => {
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.click('text=Login Form')
    await page.waitForTimeout(500)

    // Empty form validation
    const submitButton = page.locator('button[type="submit"]')
    await submitButton.click()
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/07-login-form-validation.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })

    // Fill form partially to show interaction
    await page.fill('input[type="email"]', 'user@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/08-login-form-filled.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })
  })

  test('Register form interactions', async ({ page }) => {
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.click('text=Register Form')
    await page.waitForTimeout(500)

    // Fill form to show password strength indicator
    await page.fill('input[name="name"]', 'John Doe')
    await page.fill('input[name="email"]', 'john@example.com')
    await page.fill('input[name="password"]', 'weak')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/09-register-form-weak-password.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 700 }
    })

    // Strong password
    await page.fill('input[name="password"]', 'StrongPassword123!')
    await page.fill('input[name="confirmPassword"]', 'StrongPassword123!')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/10-register-form-strong-password.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 700 }
    })
  })

  test('Password reset flow', async ({ page }) => {
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.click('text=Password Reset')
    await page.waitForTimeout(500)

    // Initial forgot password form
    await page.screenshot({
      path: 'screenshots/11-password-reset-initial.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })

    // Fill email and show validation
    await page.fill('input[type="email"]', 'user@example.com')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/12-password-reset-filled.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })
  })

  test('Mobile responsive views', async ({ page }) => {
    // Mobile viewport
    await page.setViewportSize({ width: 375, height: 667 })
    
    // Home page mobile
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/13-home-mobile.png',
      fullPage: true
    })

    // Demo page mobile
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/14-demo-mobile.png',
      fullPage: true
    })

    // Login form mobile
    await page.click('text=Login Form')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/15-login-form-mobile.png',
      fullPage: true
    })
  })

  test('Header navigation states', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 720 })
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')

    // Desktop header
    await page.screenshot({
      path: 'screenshots/16-header-desktop.png',
      fullPage: false,
      clip: { x: 0, y: 0, width: 1280, height: 100 }
    })

    // Mobile header with menu
    await page.setViewportSize({ width: 375, height: 667 })
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    
    // Show mobile menu button
    await page.screenshot({
      path: 'screenshots/17-header-mobile.png',
      fullPage: false,
      clip: { x: 0, y: 0, width: 375, height: 100 }
    })

    // Open mobile menu
    const menuButton = page.locator('button', { hasText: '' }).first()
    await menuButton.click()
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/18-header-mobile-menu.png',
      fullPage: false,
      clip: { x: 0, y: 0, width: 375, height: 400 }
    })
  })
})