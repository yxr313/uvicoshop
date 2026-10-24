"use client";

import { useState } from "react";
import { BillingToggle } from "@/components/marketing/pricing/billing-toggle";
import { PricingTable } from "@/components/marketing/pricing/pricing-table";
import { FAQItem } from "@/components/marketing/pricing/faq-item";

export default function PlanesPage() {
  const [cycle, setCycle] = useState<"monthly" | "annual">("monthly");

  const faqs = [
    { q: "¿Cómo funciona la prueba gratuita de 7 días?", a: "Al registrarte en el Plan Fundador, tienes 7 días de acceso total sin necesidad de ingresar una tarjeta de crédito. Al finalizar, decides si continúas por $9/mes." },
    { q: "¿Qué pasa si supero los límites de mi plan?", a: "No bloquearemos tu tienda de inmediato. Te notificaremos cuando estés cerca del límite (ej. cuentas de staff o envíos de email) para que puedas actualizar tu plan sin interrumpir tus ventas." },
    { q: "¿Cómo se cobra la comisión de dropshipping del 5%?", a: "El 5% se descuenta automáticamente de los fondos retenidos de cada venta de dropshipping antes de liberar tu ganancia neta los viernes. No genera facturas separadas." },
    { q: "¿Puedo vender mi propia mercancía?", a: "Sí. Si tienes tu propio inventario, el 5% de comisión no aplica. Los fondos de esas ventas van directo a tu cuenta de pasarela de pago sin pasar por la custodia de Uvico." },
    { q: "¿Necesito conocimientos técnicos para usar la API?", a: "El plan Profesional incluye acceso a API y Webhooks para que puedas conectar Uvico con sistemas externos (ERPs, CRMs). El plan Empresa incluye soporte técnico dedicado para integraciones." },
    { q: "¿Qué métodos de pago aceptan para la suscripción?", a: "Aceptamos Pago Móvil, C2P, transferencias bancarias, Zelle y USDT. Los precios están fijados en dólares (USD) y puedes pagar en bolívares a la tasa oficial del BCV." },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-7xl mx-auto">
        <div className="max-w-3xl mb-12">
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">Un precio para cada etapa.</h1>
          <p className="text-body-lg text-muted-foreground">Empieza gratis hoy. Escala cuando estés listo. Sin costos ocultos.</p>
        </div>
        <BillingToggle cycle={cycle} onChange={setCycle} />
        <PricingTable cycle={cycle} />
        <div className="mt-32 max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-12 gap-12 md:gap-24">
          <div className="md:col-span-4">
            <div className="sticky top-32">
              <h2 className="text-display-3 text-foreground">Preguntas frecuentes.</h2>
              <p className="text-body text-muted-foreground mt-4 max-w-xs">¿Aún tienes dudas? <a href="mailto:contacto@uvico.shop" className="text-foreground underline underline-offset-4 hover:text-muted-foreground transition-colors">Escríbenos a contacto@uvico.shop</a></p>
            </div>
          </div>
          <div className="md:col-span-8">
            <div className="border-t border-border">
              {faqs.map((faq, index) => (
                <FAQItem key={index} question={faq.q} answer={faq.a} />
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
