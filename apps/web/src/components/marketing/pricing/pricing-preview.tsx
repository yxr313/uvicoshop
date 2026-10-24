import Link from "next/link";
import { PLANS, getPlanPrice, FOUNDER_SLOTS_TOTAL, FOUNDER_TRIAL_DAYS } from "@uvico/pricing";
import { Button } from "@/components/ui/button";
import { Check } from "lucide-react";

export function PricingPreview() {
  const plansArray = [PLANS.fundador, PLANS.esencial, PLANS.profesional, PLANS.empresa];

  return (
    <section className="py-24 px-4 bg-background">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-display-2 text-foreground mb-8">
            Un precio para cada etapa.
          </h2>
          <Link href="/planes">
            <Button variant="outline" size="lg">
              Ver detalles de los planes
            </Button>
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {plansArray.map((plan) => {
            const price = getPlanPrice(plan.id, "monthly");
            const isFounder = plan.isFounder;

            return (
              <div 
                key={plan.id} 
                className={`relative flex flex-col p-8 rounded-[2rem] border ${isFounder ? 'border-foreground shadow-lg bg-white' : 'border-border bg-neutral-50'} overflow-hidden`}
              >
                {isFounder && (
                  <div className="absolute top-0 inset-x-0 h-1.5 bg-accent" />
                )}
                
                <div className="mb-8">
                  <h3 className="text-heading-lg font-display text-foreground mb-2">{plan.name}</h3>
                  <div className="flex items-baseline gap-1">
                    <span className="text-[3rem] leading-none font-display font-medium tracking-tighter text-foreground">
                      ${price.amountUsd}
                    </span>
                    <span className="text-body text-muted-foreground">/mes</span>
                  </div>
                </div>

                <div className="flex-1 space-y-4 mb-8">
                  <div className="flex items-center gap-3">
                    <Check size={18} className="text-foreground" />
                    <span className="text-body-sm text-foreground">Tienda online ilimitada</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <Check size={18} className="text-foreground" />
                    <span className="text-body-sm text-foreground">Pasarelas venezolanas nativas</span>
                  </div>
                  {isFounder && (
                    <>
                      <div className="flex items-center gap-3">
                        <Check size={18} className="text-foreground" />
                        <span className="text-body-sm text-foreground">{FOUNDER_TRIAL_DAYS} días de prueba gratis</span>
                      </div>
                      <div className="flex items-center gap-3">
                        <Check size={18} className="text-foreground" />
                        <span className="text-body-sm text-foreground">Acceso a red de proveedores</span>
                      </div>
                    </>
                  )}
                  {!isFounder && (
                    <div className="flex items-center gap-3">
                      <Check size={18} className="text-foreground" />
                      <span className="text-body-sm text-foreground">Todo lo del plan anterior</span>
                    </div>
                  )}
                </div>

                <Link href={`/registro?plan=${plan.id}`} className="w-full mt-auto">
                  <Button variant={isFounder ? "default" : "dark"} className="w-full">
                    Empezar ahora
                  </Button>
                </Link>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
