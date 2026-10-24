import { ArrowRight, Heart, Utensils, Home, HeartPulse, ShieldCheck, Quote } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function YamalNetworkPage() {
  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-7xl mx-auto">
        
        {/* HERO */}
        <div className="max-w-4xl mb-24">
          <div className="inline-flex items-center gap-2 text-accent mb-6">
            <Heart className="w-5 h-5 fill-current" />
            <span className="text-sm font-bold tracking-widest uppercase">Fundación Humanitaria</span>
          </div>
          <h1 className="text-display-1 mb-8 text-foreground tracking-tight">
            La belleza de una sociedad se mide en cómo cuida a los suyos.
          </h1>
          <p className="text-body-lg text-muted-foreground mb-10 max-w-2xl">
            Yamal Network es el brazo humanitario de Uvico. Nacimos para llevar asistencia directa, alimentación y refugio a las familias venezolanas que enfrentan situaciones de vulnerabilidad extrema.
          </p>
          <div className="flex flex-col sm:flex-row gap-4">
            <a href="#donar">
              <Button variant="default" size="lg" className="w-full sm:w-auto">Hacer una donación</Button>
            </a>
            <a href="mailto:yamal@uvico.shop?subject=Quiero ser voluntario">
              <Button variant="outline" size="lg" className="w-full sm:w-auto">Ser voluntario</Button>
            </a>
          </div>
        </div>

        {/* MÉTRICAS DE IMPACTO (Fase Piloto - Transparencia Real) */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 border-y border-border py-12 mb-32">
          <div>
            <p className="text-display-2 font-display text-foreground mb-2">120+</p>
            <p className="text-body font-medium text-muted-foreground">Familias impactadas (Fase Piloto)</p>
          </div>
          <div>
            <p className="text-display-2 font-display text-foreground mb-2">3</p>
            <p className="text-body font-medium text-muted-foreground">Viviendas reconstruidas</p>
          </div>
          <div>
            <p className="text-display-2 font-display text-foreground mb-2">450+</p>
            <p className="text-body font-medium text-muted-foreground">Pacientes atendidos en jornadas</p>
          </div>
        </div>

        {/* PROGRAMAS */}
        <div className="mb-32">
          <h2 className="text-display-2 text-foreground mb-12">Nuestros programas.</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
            <div>
              <Utensils className="w-8 h-8 text-foreground mb-6" />
              <h3 className="text-heading-lg font-display text-foreground mb-3">Alimentación y Nutrición</h3>
              <p className="text-body text-muted-foreground">
                Entregamos canastillas básicas mensuales y operamos comedores comunitarios en zonas de alta vulnerabilidad, garantizando al menos una comida balanceada al día para niños y adultos mayores.
              </p>
            </div>
            <div>
              <Home className="w-8 h-8 text-foreground mb-6" />
              <h3 className="text-heading-lg font-display text-foreground mb-3">Refugio y Hábitat</h3>
              <p className="text-body text-muted-foreground">
                Trabajamos en la reconstrucción de techos, paredes y sistemas de agua potable en hogares afectados por las lluvias o el abandono, devolviendo la dignidad y seguridad a las familias.
              </p>
            </div>
            <div>
              <HeartPulse className="w-8 h-8 text-foreground mb-6" />
              <h3 className="text-heading-lg font-display text-foreground mb-3">Salud y Bienestar</h3>
              <p className="text-body text-muted-foreground">
                Organizamos jornadas médicas con especialistas voluntarios, donamos medicamentos esenciales y brindamos apoyo psicológico a niños y familias que han atravesado crisis severas.
              </p>
            </div>
          </div>
        </div>

        {/* TRANSPARENCIA Y EL NOMBRE */}
        <div className="bg-neutral-50 rounded-[2rem] p-12 border border-border mb-32 grid grid-cols-1 md:grid-cols-2 gap-16 items-center">
          <div>
            <h2 className="text-display-3 text-foreground mb-6">¿Por qué "Yamal"?</h2>
            <p className="text-body-lg text-muted-foreground mb-6">
              "Yamal" (جمال) significa belleza y perfección. Creemos que la verdadera belleza de un país no está en sus edificios, sino en cómo protege a sus ciudadanos más vulnerables. No enseñamos a programar a alguien que tiene hambre; primero le damos de comer.
            </p>
            <p className="text-body-lg text-muted-foreground">
              Uvico financia los costos operativos de Yamal Network, lo que significa que <strong>el 100% de tu donación va directo al campo</strong>.
            </p>
          </div>
          
          {/* Gráfico de Transparencia */}
          <div className="bg-white p-8 rounded-2xl border border-border shadow-sm">
            <div className="flex items-center gap-3 mb-6">
              <ShieldCheck className="w-6 h-6 text-accent" />
              <h3 className="text-heading font-display text-foreground">Transparencia total</h3>
            </div>
            <div className="space-y-6">
              <div>
                <div className="flex justify-between text-sm font-medium mb-2">
                  <span className="text-foreground">Impacto directo (Donaciones)</span>
                  <span className="text-foreground">100%</span>
                </div>
                <div className="w-full h-3 bg-neutral-100 rounded-full overflow-hidden">
                  <div className="h-full bg-accent w-full"></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm font-medium mb-2">
                  <span className="text-muted-foreground">Costos operativos</span>
                  <span className="text-muted-foreground">0% (Cubierto por Uvico)</span>
                </div>
                <div className="w-full h-3 bg-neutral-100 rounded-full overflow-hidden">
                  <div className="h-full bg-foreground w-0"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* HISTORIA DE IMPACTO */}
        <div className="max-w-4xl mx-auto mb-32 border-l-4 border-accent pl-8">
          <Quote className="w-10 h-10 text-neutral-300 mb-6" />
          <blockquote className="text-display-3 font-display text-foreground italic mb-6 leading-tight">
            "Cuando perdimos el techo por las lluvias, pensamos que lo habíamos perdido todo. Yamal no solo nos ayudó a reconstruir las paredes, nos devolvió la esperanza de que no estábamos solos."
          </blockquote>
          <p className="text-body font-medium text-foreground">— María G.</p>
          <p className="text-sm text-muted-foreground">Beneficiaria del programa de Vivienda (Petare, Caracas)</p>
        </div>

        {/* CTA FINAL / DONACIONES */}
        <div id="donar" className="border-t border-border pt-20">
          <div className="bg-[#0a0a0a] text-white rounded-[2rem] p-12 md:p-20 text-center relative overflow-hidden">
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-32 bg-white/5 blur-3xl rounded-full" />
            
            <div className="relative z-10 max-w-2xl mx-auto">
              {/* CORREGIDO: Agregado text-white */}
              <h2 className="text-display-2 font-display mb-6 text-white">Tu ayuda cambia vidas hoy.</h2>
              <p className="text-body-lg text-white/60 mb-10">
                Cada aporte, por pequeño que sea, se convierte en un plato de comida, un techo seguro o una medicina. Únete a la red de personas que están reconstruyendo Venezuela desde la base.
              </p>
              
              <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                <a href="mailto:yamal@uvico.shop?subject=Quiero hacer una donación">
                  <Button variant="default" size="lg" className="w-full sm:w-auto px-10 text-lg h-14">
                    Donar ahora
                  </Button>
                </a>
                <p className="text-sm text-white/40 mt-4 sm:mt-0">
                  Aceptamos Pago Móvil, Zelle, Binance y transferencias.
                </p>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  );
}
