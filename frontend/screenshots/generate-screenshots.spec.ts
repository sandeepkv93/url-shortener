import { test, expect } from '@playwright/test'

test.describe('Comprehensive UI Screenshots', () => {
  test.beforeEach(async ({ page }) => {
    // Set viewport size for consistent screenshots
    await page.setViewportSize({ width: 1280, height: 720 })
  })

  test('01 - Main Application Pages', async ({ page }) => {
    // Home Page
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/01-home-page.png',
      fullPage: true
    })

    // Dashboard Page (no auth required for screenshot)
    await page.goto('/dashboard')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/02-dashboard-page.png',
      fullPage: true
    })

    // Analytics Page
    await page.goto('/analytics')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/03-analytics-page.png',
      fullPage: true
    })

    // Profile Page
    await page.goto('/profile')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/04-profile-page.png',
      fullPage: true
    })

    // 404 Not Found Page
    await page.goto('/nonexistent-page')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/05-not-found-page.png',
      fullPage: true
    })
  })

  test('02 - Component Demo Page', async ({ page }) => {
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)

    // Demo navigation overview
    await page.screenshot({
      path: 'screenshots/06-demo-navigation.png',
      fullPage: true
    })

    // Login Form
    await page.click('text=Login Form')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/07-login-form.png',
      fullPage: true
    })

    // Register Form
    await page.click('text=Register Form')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/08-register-form.png',
      fullPage: true
    })

    // Password Reset
    await page.click('text=Password Reset')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/09-password-reset.png',
      fullPage: true
    })

    // Common Components
    await page.click('text=Common Components')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/10-common-components.png',
      fullPage: true
    })
  })

  test('03 - Form Interactions and States', async ({ page }) => {
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    
    // Login Form Validation
    await page.click('text=Login Form')
    await page.waitForTimeout(500)
    
    const submitButton = page.locator('button[type="submit"]')
    await submitButton.click()
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/11-login-form-validation.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })

    // Login form filled state
    await page.fill('input[type="email"]', 'user@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/12-login-form-filled.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })

    // Register Form with weak password
    await page.click('text=Register Form')
    await page.waitForTimeout(500)
    await page.fill('input[name="name"]', 'John Doe')
    await page.fill('input[name="email"]', 'john@example.com')
    await page.fill('input[name="password"]', 'weak')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/13-register-form-weak-password.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 700 }
    })

    // Register form with strong password
    await page.fill('input[name="password"]', 'StrongPassword123!')
    await page.fill('input[name="confirmPassword"]', 'StrongPassword123!')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/14-register-form-strong-password.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 700 }
    })

    // Password reset form
    await page.click('text=Password Reset')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/15-password-reset-initial.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })

    // Password reset filled
    await page.fill('input[type="email"]', 'user@example.com')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/16-password-reset-filled.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })
  })

  test('04 - Mobile Responsive Views', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 })
    
    // Home page mobile
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/17-home-mobile.png',
      fullPage: true
    })

    // Dashboard mobile
    await page.goto('/dashboard')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/18-dashboard-mobile.png',
      fullPage: true
    })

    // Analytics mobile
    await page.goto('/analytics')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/19-analytics-mobile.png',
      fullPage: true
    })

    // Demo page mobile
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/20-demo-mobile.png',
      fullPage: true
    })

    // Login form mobile
    await page.click('text=Login Form')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/21-login-form-mobile.png',
      fullPage: true
    })

    // Register form mobile
    await page.click('text=Register Form')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/22-register-form-mobile.png',
      fullPage: true
    })
  })

  test('05 - Navigation and Layout Components', async ({ page }) => {
    // Desktop header
    await page.setViewportSize({ width: 1280, height: 720 })
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    await page.screenshot({
      path: 'screenshots/23-header-desktop.png',
      fullPage: false,
      clip: { x: 0, y: 0, width: 1280, height: 100 }
    })

    // Footer desktop
    await page.screenshot({
      path: 'screenshots/24-footer-desktop.png',
      fullPage: false,
      clip: { x: 0, y: 580, width: 1280, height: 140 }
    })

    // Mobile header and navigation
    await page.setViewportSize({ width: 375, height: 667 })
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    
    await page.screenshot({
      path: 'screenshots/25-header-mobile.png',
      fullPage: false,
      clip: { x: 0, y: 0, width: 375, height: 100 }
    })

    // Try to open mobile menu if it exists
    try {
      const menuButton = page.locator('button[aria-label*="menu" i], button:has-text("☰"), button:has-text("Menu"), [data-testid="mobile-menu"]').first()
      if (await menuButton.isVisible()) {
        await menuButton.click()
        await page.waitForTimeout(500)
        await page.screenshot({
          path: 'screenshots/26-mobile-menu-open.png',
          fullPage: false,
          clip: { x: 0, y: 0, width: 375, height: 400 }
        })
      }
    } catch (error) {
      // Mobile menu might not exist or be implemented differently
      console.log('Mobile menu not found or clickable')
    }
  })

  test('06 - Loading and Error States', async ({ page }) => {
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.click('text=Common Components')
    await page.waitForTimeout(1000)

    // Focus on loading components section
    await page.screenshot({
      path: 'screenshots/27-loading-components.png',
      fullPage: false,
      clip: { x: 0, y: 200, width: 1280, height: 600 }
    })

    // Error state (404 page)
    await page.goto('/this-page-does-not-exist')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/28-error-404-state.png',
      fullPage: true
    })
  })

  test('07 - Tablet Responsive Views', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 })
    
    // Home page tablet
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/29-home-tablet.png',
      fullPage: true
    })

    // Dashboard tablet
    await page.goto('/dashboard')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: 'screenshots/30-dashboard-tablet.png',
      fullPage: true
    })

    // Demo page tablet
    await page.goto('/demo')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(500)
    await page.screenshot({
      path: 'screenshots/31-demo-tablet.png',
      fullPage: true
    })
  })
})