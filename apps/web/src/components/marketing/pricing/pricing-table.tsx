"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Check, Minus } from "lucide-react";
// CORREGIDO: Se añadió formatUsd al import
import { formatUsd, PLANS, getPlanPrice, FOUNDER_SLOTS_TOTAL, FOUNDER_TRIAL_DAYS } from "@uvico/pricing";
import { AnimatedPrice } from "./animated-price";

interface PricingTableProps {
  cycle: "monthly" | "annual";
}

export function PricingTable({ cycle }: PricingTableProps) {
  const planIds = ["fundador", "esencial", "profesional", "empresa"] as const;
  const slotsTaken = 312;
  const slotsLeft = FOUNDER_SLOTS_TOTAL - slotsTaken;

  const categories = [
    {
      title: "Tienda y Productos",
      features: [
        { name: "Tienda online", values: ["Ilimitada", "Ilimitada", "Ilimitada", "Ilimitada"] },
        { name: "Dominio personalizado", values: [false, true, true, true] },
        { name: "Pasarelas venezolanas nativas", values: [true, true, true, true] },
        { name: "Diseño de tienda (Editor Puck)", values: ["Básico", "Avanzado", "Ilimitado", "Ilimitado"] },
        { name: "Moneda local (VES) y USD", values: [true, true, true, true] },
        { name: "Carga de productos masiva", values: [false, true, true, true] },
      ]
    },
    {
      title: "Dropshipping y Proveedores",
      features: [
        { name: "Red de proveedores verificados", values: ["Básica", "Completa", "Prioritaria", "Dedicada"] },
        { name: "Comisión por venta", values: ["5%", "5%", "5%", "5%"] },
        { name: "Sincronización de inventario", values: ["Manual", "Automática", "Tiempo real", "Tiempo real"] },
        { name: "Gestión de pedidos", values: ["Ilimitada", "Ilimitada", "Ilimitada", "Ilimitada"] },
        { name: "Conexión con múltiples proveedores", values: [false, true, true, true] },
      ]
    },
    {
      title: "Analíticas y Reportes",
      features: [
        { name: "Visitas y conversión", values: ["Básico", "Avanzado", "Avanzado", "Avanzado"] },
        { name: "Reportes exportables (CSV)", values: [false, true, true, true] },
        { name: "Dashboard de ventas en vivo", values: [false, false, true, true] },
        { name: "Análisis de cohortes y retención", values: [false, false, false, true] },
        { name: "Integración con Google Analytics", values: [true, true, true, true] },
      ]
    },
    {
      title: "Marketing y Automatización",
      features: [
        { name: "Recuperación de carritos", values: [false, "Básica", "Avanzada", "Avanzada"] },
        { name: "Cupones y descuentos", values: [true, true, true, true] },
        { name: "Email marketing (envíos/mes)", values: ["1,000", "5,000", "50,000", "Ilimitados"] },
        { name: "SEO Optimización", values: ["Básico", "Avanzado", "Avanzado", "Avanzado"] },
        { name: "Programa de afiliados", values: [false, false, true, true] },
      ]
    },
    {
      title: "Equipo y Seguridad",
      features: [
        { name: "Cuentas de staff", values: ["1", "3", "5", "Ilimitadas"] },
        { name: "Roles y permisos granulares", values: [false, true, true, true] },
        { name: "Autenticación de 2 factores (2FA)", values: [true, true, true, true] },
        { name: "API y Webhooks", values: [false, false, true, true] },
        { name: "Multi-tienda", values: [false, false, false, true] },
        { name: "SSO (Inicio de sesión único)", values: [false, false, false, true] },
        { name: "Logs de auditoría", values: [false, false, true, true] },
      ]
    },
    {
      title: "Soporte",
      features: [
        { name: "Tiempo de respuesta", values: ["48h", "24h", "12h", "4h"] },
        { name: "Chat en vivo", values: [false, false, true, true] },
        { name: "Gerente de cuenta dedicado", values: [false, false, false, true] },
        { name: "SLA garantizado", values: [false, false, false, true] },
      ]
    }
  ];

  return (
    <div className="w-full">
      <div className="grid grid-cols-5 gap-4 pb-8 border-b border-border">
        <div className="col-span-1 flex items-end pb-4">
          <span className="text-sm font-medium text-muted-foreground">Planes</span>
        </div>
        {planIds.map((id) => {
          const plan = PLANS[id];
          const price = getPlanPrice(id, cycle);
          const isPopular = id === "profesional";
          const isFounder = id === "fundador";

          return (
            <div key={id} className="col-span-1 flex flex-col items-start pt-4">
              <div className="flex items-center gap-2 mb-2 h-6">
                <h3 className="text-heading font-display text-foreground">{plan.name}</h3>
                {isPopular && (
                  <span className="bg-accent text-accent-foreground text-[10px] font-bold px-2 py-0.5 rounded">
                    Recomendado
                  </span>
                )}
              </div>
              
              <div className="flex items-baseline gap-1 mb-2">
                <AnimatedPrice value={price.effectiveMonthlyUsd} />
                <span className="text-sm text-muted-foreground">/mes</span>
              </div>
              
              <p className="text-xs text-muted-foreground mb-6 h-10">
                {isFounder 
                  ? `Quedan ${slotsLeft} de 500 cupos. ${FOUNDER_TRIAL_DAYS} días de prueba gratis.` 
                  : cycle === "annual" 
                    ? `Facturado ${formatUsd(price.amountUsd)} al año` 
                    : "Facturado mensualmente"}
              </p>
              
              <Link href={`/registro?plan=${id}`} className="block w-full mt-auto">
                <Button 
                  variant={isPopular ? "default" : "outline"} 
                  className="w-full"
                >
                  Empezar ahora
                </Button>
              </Link>
            </div>
          );
        })}
      </div>

      {categories.map((category, catIdx) => (
        <div key={catIdx} className="mb-8">
          <div className="grid grid-cols-5 gap-4 py-6 border-b border-border">
            <div className="col-span-1">
              <h4 className="text-sm font-semibold text-foreground">{category.title}</h4>
            </div>
          </div>

          <div className="divide-y divide-border/50">
            {category.features.map((feature, idx) => (
              <div key={idx} className="grid grid-cols-5 gap-4 py-3 items-center hover:bg-neutral-50/50 transition-colors">
                <div className="col-span-1 text-sm text-muted-foreground">{feature.name}</div>
                {feature.values.map((value, i) => (
                  <div key={i} className="col-span-1 text-sm text-center">
                    {typeof value === "boolean" ? (
                      value ? (
                        <Check className="inline-block w-4 h-4 text-foreground" />
                      ) : (
                        <Minus className="inline-block w-4 h-4 text-neutral-300" />
                      )
                    ) : (
                      <span className="text-foreground font-medium">{value}</span>
                    )}
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
