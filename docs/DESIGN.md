---
name: Aqueous Narrative
colors:
  surface: '#f8f9ff'
  surface-dim: '#d7dae1'
  surface-bright: '#f8f9ff'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f1f3fb'
  surface-container: '#ebeef5'
  surface-container-high: '#e5e8ef'
  surface-container-highest: '#e0e2ea'
  on-surface: '#181c21'
  on-surface-variant: '#404751'
  inverse-surface: '#2d3136'
  inverse-on-surface: '#eef1f8'
  outline: '#707883'
  outline-variant: '#c0c7d3'
  surface-tint: '#0061a3'
  primary: '#005f9f'
  on-primary: '#ffffff'
  primary-container: '#0078c7'
  on-primary-container: '#fdfcff'
  inverse-primary: '#9ecaff'
  secondary: '#a43b38'
  on-secondary: '#ffffff'
  secondary-container: '#fc7d76'
  on-secondary-container: '#711517'
  tertiary: '#7a5500'
  on-tertiary: '#ffffff'
  tertiary-container: '#9a6c00'
  on-tertiary-container: '#fffbff'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#d1e4ff'
  primary-fixed-dim: '#9ecaff'
  on-primary-fixed: '#001d36'
  on-primary-fixed-variant: '#00497c'
  secondary-fixed: '#ffdad7'
  secondary-fixed-dim: '#ffb3ad'
  on-secondary-fixed: '#410004'
  on-secondary-fixed-variant: '#842323'
  tertiary-fixed: '#ffdeab'
  tertiary-fixed-dim: '#fbbc44'
  on-tertiary-fixed: '#271900'
  on-tertiary-fixed-variant: '#5f4100'
  background: '#f8f9ff'
  on-background: '#181c21'
  surface-variant: '#e0e2ea'
typography:
  display-lg:
    fontFamily: Libre Caslon Text
    fontSize: 48px
    fontWeight: '400'
    lineHeight: 56px
    letterSpacing: -0.02em
  display-lg-mobile:
    fontFamily: Libre Caslon Text
    fontSize: 36px
    fontWeight: '400'
    lineHeight: 44px
  headline-md:
    fontFamily: Libre Caslon Text
    fontSize: 32px
    fontWeight: '400'
    lineHeight: 40px
  body-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Plus Jakarta Sans
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  label-sm:
    fontFamily: Plus Jakarta Sans
    fontSize: 12px
    fontWeight: '600'
    lineHeight: 16px
    letterSpacing: 0.05em
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  canvas-margin: 2.5rem
  gutter-fluid: 2rem
  organic-gap: 1.5rem
  stack-sm: 0.75rem
---

## Brand & Style
This design system captures the fluid, expressive nature of watercolor painting. It is designed to evoke a sense of creativity, tranquility, and organic beauty, targeting art-centric platforms, lifestyle journals, or premium wellness applications. 

The aesthetic blends **Minimalism** with **Glassmorphism**, utilizing "pigment washes" and soft blurs to simulate paint bleeding into paper. The UI prioritizes negative space (the "canvas") to allow the vibrant, translucent colors to breathe. Interaction should feel soft and intentional, avoiding the mechanical rigidity of traditional enterprise software in favor of an airy, artisanal atmosphere.

## Colors
The palette is built on "The Wash"—a series of translucent, vibrant tones that mimic pigments suspended in water. 

- **Primary (Cerulean):** A deep yet bright blue used for core interactions and brand presence.
- **Secondary (Coral):** A warm, energetic pink for highlights and secondary calls to action.
- **Tertiary (Amber):** A sunny, optimistic yellow for warnings, notifications, or decorative accents.
- **Accent (Seafoam):** A lush, calming green for success states or specialized categories.
- **Surface:** The background is an off-white `background_paper`, providing a textured, natural foundation. 

Use varying opacities (20% to 80%) for container backgrounds to create the "overlap" effect common in watercolor layering.

## Typography
The typography balances the traditional authority of an editorial serif with the friendly accessibility of a modern sans-serif. 

**Libre Caslon Text** is utilized for headlines to ground the fluid visuals with a sense of literary refinement. **Plus Jakarta Sans** provides a soft, rounded contrast for body copy and UI labels, ensuring high legibility even when layered over colorful washes. Text should primarily appear in `neutral_ink` to maintain high contrast against the paper-tinted backgrounds.

## Layout & Spacing
The layout follows an **airy and organic** philosophy. Instead of a rigid, dense grid, it uses a generous fluid layout with wide margins to mimic the composition of a painting.

- **Desktop:** A 12-column fluid grid with wide 32px gutters. Content is often offset or asymmetrical to avoid a "boxed-in" feel.
- **Mobile:** A single-column layout with 20px side margins.
- **Rhythm:** Use "Organic Gaps" for vertical spacing between different content types, allowing elements to feel like they are floating on the canvas rather than stacked in a list.

## Elevation & Depth
Depth is created through **Tonal Layers** and **Backdrop Blurs** rather than traditional drop shadows.

- **Pigment Layering:** Use semi-transparent backgrounds (e.g., 15% opacity Primary Blue) for containers. When containers overlap, their colors should multiply, simulating physical paint layering.
- **Soft Diffusion:** For elevated elements like modals or floating menus, use a high-radius backdrop blur (20px) with a very faint, tinted shadow (e.g., 5% opacity of the element’s primary color) to suggest a soft lift off the paper.
- **Outline:** Elements use ultra-thin (1px) borders in a slightly darker shade of the surface color to define boundaries without adding visual weight.

## Shapes
Shapes are defined by "softness." Sharp corners are avoided to maintain the fluid, liquid metaphor. Standard components use a 0.5rem radius, while larger interactive "puddle" areas or decorative containers may use higher, irregular radii to simulate natural water droplets.

## Components
- **Buttons:** Primary buttons feature a soft gradient "wash" from the primary color to a slightly lighter tint. They should have a subtle hover effect that expands the "bleed" of the color.
- **Chips:** Highly translucent backgrounds (10% opacity) with a solid 1px border of the same hue. Use for tags and categories.
- **Input Fields:** Minimalist design—bottom border only, mimicking a ruled line on paper, or a very light, soft-rounded background tint that darkens slightly on focus.
- **Cards:** Use "The Wash" backgrounds. Cards should not have heavy shadows; instead, use a slight color shift or a very soft, diffused inner glow to signify presence.
- **Checkboxes & Radios:** These should appear "hand-drawn" in spirit—rounded, using the Primary color for the fill with a white checkmark.
- **Lists:** Separated by generous whitespace rather than divider lines. If dividers are necessary, use a "deckle edge" or a very faint, feathered line.