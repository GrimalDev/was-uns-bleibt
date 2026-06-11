# Dual-Screen Raspberry Pi Web Interface

## Overview

This project runs a simple web application on a Raspberry Pi and displays two different browser views across two physical screens.

- Single web app stack
- Two browser windows (one per display)
- No route-level monitor separation required initially
- Secondary screen can load a dedicated view directly via URL

## Initial Architecture Decision

We use **one web application** instead of two separate programs.

Why:

- Shared state and logic in one codebase
- Faster iteration and simpler deployment
- Easy dual-screen control using browser startup flags in kiosk/fullscreen mode

## Screen Strategy (Phase 1)

- Screen A (primary): load the main app view
- Screen B (secondary): load a specific view directly
- Both windows run fullscreen/kiosk on boot

At this stage, the app does not need strict route-to-monitor enforcement.

## Fullscreen / Kiosk Mode

Raspberry Pi Chromium supports fullscreen kiosk mode.

Typical launch approach:

- Start two Chromium instances/windows
- Pass kiosk/fullscreen flags
- Set window position/size so each window maps to a specific display

This enables unattended startup for a kiosk-like dual-screen setup.

## Tech Stack

- **Backend API:** Go
- **Frontend SPA:** Svelte
- **Rendering Engine:** PixiJS (WebGL)
