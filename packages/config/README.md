# @uvico/config

Sistema de diseño, configuraciones de TypeScript y ESLint de UVICO.

## Consumo en apps/web

En `apps/web/src/app/globals.css`:

    @import "tailwindcss";
    @import "@uvico/config/tailwind";

En `apps/web/tsconfig.json`:

    {
      "extends": "@uvico/config/typescript/nextjs.json"
    }

En `apps/web/eslint.config.mjs`:

    import { nextjs } from "@uvico/config/eslint/nextjs";
    export default [...nextjs];

## ADN visual

- Fondo: blanco puro, sin modo oscuro.
- Texto: negro.
- Acento: reservado solo para CTAs, focus y bordes interactivos.
- Tipografía: Bricolage Grotesque (display), Instrument Sans (body).
- Motion: 200-300ms con easing suave.

## Reglas

- Cero emojis, badges, subtítulos en títulos, botones en títulos.
- Máximo un elemento accent por viewport.
- prefers-reduced-motion respetado siempre.
