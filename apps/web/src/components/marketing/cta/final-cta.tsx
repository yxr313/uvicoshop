import { SlugBuilder } from "@/components/marketing/hero/slug-builder";

export function FinalCTA() {
  return (
    <section className="py-32 px-4 bg-background">
      <div className="max-w-6xl mx-auto">
        <div className="relative bg-[#0a0a0a] rounded-[2rem] md:rounded-[3rem] py-24 px-6 md:px-16 flex flex-col items-center text-center overflow-hidden">
          
          {/* Resplandor sutil superior */}
          <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-white/20 to-transparent" />
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-1/2 h-32 bg-white/5 blur-3xl rounded-full" />

          <h2 className="relative z-10 text-display-2 text-white mb-12 max-w-3xl mx-auto tracking-tight">
            Tu tienda está a un clic.
          </h2>
          
          {/* Slug ancho para el CTA, con botón accent */}
          <div className="relative z-10 w-full max-w-2xl">
            <SlugBuilder variant="default" />
          </div>

        </div>
      </div>
    </section>
  );
}
