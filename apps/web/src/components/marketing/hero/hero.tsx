import { SlugBuilder } from "./slug-builder";
import { TrustBar } from "./trust-bar";

export function Hero() {
  return (
    <section className="relative pt-48 pb-24 px-4 flex flex-col items-center text-center">
      <div className="max-w-5xl mx-auto w-full">
        
        {/* Headline */}
        <h1 className="text-display-1 mb-12 text-foreground max-w-4xl mx-auto">
          Vende online como los grandes.
        </h1>

        {/* Slug Register (ahora limitado a max-w-xl) */}
        <div className="max-w-xl mx-auto mb-16">
          <SlugBuilder variant="dark" />
        </div>

        {/* Trust Bar (Carrusel) */}
        <TrustBar />
      </div>
    </section>
  );
}
